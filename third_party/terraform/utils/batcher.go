package google

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const defaultBatchSendIntervalSec = 10

// RequestBatcher a global batcher object.
type RequestBatcher struct {
	sync.Mutex

	*batchingConfig
	batches map[string]*startedBatch
}

// BatchRequest represents a single request to a global batcher.
type BatchRequest struct {
	// ResourceName determine the resource name to be passed to SendF.
	ResourceName string

	// Body is this request's data to be passed to SendF, and may be combined
	// with other bodies using CombineF.
	Body interface{}

	// CombineF function determines how to combine bodies from two batches.
	CombineF batcherCombineFunc

	// CombineF function determines how to actually send a batched request to a
	// third party service.
	SendF batcherSendFunc

	// ID for debugging request. This should be specific to a single request
	// (i.e. per Terraform resource)
	DebugId string
}

// These types are meant to be the public interface to batchers. They define
// logic to manage batch data type and behavior, and require service-specific
// implementations per type of request per service.
// Function type for combine existing batches and additional batch data
type batcherCombineFunc func(body interface{}, toAdd interface{}) (interface{}, error)

// Function type for sending a batch request
type batcherSendFunc func(resourceName string, body interface{}) (interface{}, error)

// batchResponse bundles an API response (data, error) tuple.
type batchResponse struct {
	body interface{}
	err  error
}

// startedBatch refers to a processed batch whose timer to send the request has
// already been started. The responses for the request is sent to each listener
// channel, representing parallel callers that are waiting on requests
// combined into this batch.
type startedBatch struct {
	*BatchRequest

	listeners []chan batchResponse
	timer     *time.Timer
}

// batchingConfig contains user configuration for controlling batch requests.
type batchingConfig struct {
	sendAfter       time.Duration
	disableBatching bool
}

// Initializes a new batcher
func NewRequestBatcher(config *batchingConfig) *RequestBatcher {
	return &RequestBatcher{
		batchingConfig: config,
		batches:        make(map[string]*startedBatch),
	}
}

// SendRequestWithTimeout is expected to be called per parallel call.
// It manages waiting on the result of a batch request.
func (b *RequestBatcher) SendRequestWithTimeout(batchType string, request *BatchRequest, timeout time.Duration) (interface{}, error) {
	if request == nil {
		return nil, fmt.Errorf("error, cannot request batching for nil BatchRequest")
	}
	if request.CombineF == nil {
		return nil, fmt.Errorf("error, cannot request batching for BatchRequest with nil CombineF")
	}
	if request.SendF == nil {
		return nil, fmt.Errorf("error, cannot request batching for BatchRequest with nil SendF")
	}
	if b.disableBatching {
		log.Printf("[DEBUG] Batching is disabled, sending single request for %q", request.DebugId)
		return request.SendF(request.ResourceName, request.Body)
	}

	respCh, err := b.startBatchRequest(batchType, request)
	if err != nil {
		return nil, fmt.Errorf("error adding request to batch: %s", err)
	}

	timer := time.NewTimer(timeout)

	select {
	case resp := <-respCh:
		if resp.err != nil {
			return nil, fmt.Errorf("Batch %q Request %q returned error: %v", batchType, request.DebugId, resp.err)
		}
		return resp.body, nil
	case <-timer.C:
		break
	}
	return nil, fmt.Errorf("Request %s timed out after %v", batchType, timeout)
}

// startBatchRequest manages batching logic that access shared information
// (i.e. existing batches)
func (b *RequestBatcher) startBatchRequest(batchType string, newRequest *BatchRequest) (<-chan batchResponse, error) {
	b.Lock()
	defer b.Unlock()

	// The calling goroutine will need a channel to wait on for a response.
	respCh := make(chan batchResponse, 1)

	batchId := fmt.Sprintf("%s:%s", batchType, newRequest.ResourceName)
	// If batch already exists, combine this request into existing request.
	if batch, ok := b.batches[batchId]; ok {
		log.Printf("[DEBUG] Adding batch request %q to existing batch %q", newRequest.DebugId, batchId)
		if batch.CombineF == nil {
			return nil, fmt.Errorf("Provider Error: unable to add request %q to batch %q with no CombineF", newRequest.DebugId, batchId)
		}

		newBody, err := batch.CombineF(batch.Body, newRequest.Body)
		if err != nil {
			return nil, fmt.Errorf("Unable to combine request %q data into existing batch %q: %v", newRequest.DebugId, batchId, err)
		}

		batch.Body = newBody
		log.Printf("[DEBUG] Added batch request %q to batch. New batch body: %v", newRequest.DebugId, batch.Body)
		batch.listeners = append(batch.listeners, respCh)
		return respCh, nil
	}

	log.Printf("[DEBUG] Creating new batch %q from request %q", newRequest.DebugId, batchId)
	// Create a new batch.
	b.batches[batchId] = &startedBatch{
		BatchRequest: newRequest,
		listeners:    []chan batchResponse{respCh},
	}

	// Start a timer to send the request
	b.batches[batchId].timer = time.AfterFunc(b.sendAfter, func() {
		batch := b.popBatch(batchId)

		var resp batchResponse
		if batch == nil {
			log.Printf("[DEBUG] Batch not found in saved batches, running single request batch %q", batchId)
			resp = newRequest.send()
		} else {
			log.Printf("[DEBUG] Sending batch %q combining %d requests)", batchId, len(batch.listeners))
			resp = batch.send()
		}

		// Send message to all goroutines waiting on result.
		for _, ch := range batch.listeners {
			ch <- resp
			close(ch)
		}
	})

	return respCh, nil
}

func (b *RequestBatcher) popBatch(batchId string) *startedBatch {
	b.Lock()
	defer b.Unlock()

	batch, ok := b.batches[batchId]
	if !ok {
		log.Printf("[DEBUG] Batch with ID %q not found in batcher", batchId)
		return nil
	}

	delete(b.batches, batchId)
	return batch
}

func (req *BatchRequest) send() batchResponse {
	if req.SendF == nil {
		return batchResponse{
			err: fmt.Errorf("provider error: Batch request has no SendBatch function"),
		}
	}
	v, err := req.SendF(req.ResourceName, req.Body)
	return batchResponse{v, err}
}
