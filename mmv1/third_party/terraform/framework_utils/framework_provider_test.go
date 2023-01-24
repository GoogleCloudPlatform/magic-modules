package google

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var fwProviders map[string]*frameworkTestProvider

type frameworkTestProvider struct {
	ProdProvider frameworkProvider
	TestName     string
}

func NewFrameworkTestProvider(testName string) *frameworkTestProvider {
	return &frameworkTestProvider{
		ProdProvider: frameworkProvider{
			version: "test",
		},
		TestName: testName,
	}
}

func (p *frameworkTestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if isVcrEnabled() {
		configsLock.RLock()
		_, ok := fwProviders[p.TestName]
		configsLock.RUnlock()
		if ok {
			return
		}
		p.ProdProvider.Configure(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			return
		}
		var vcrMode recorder.Mode
		switch vcrEnv := os.Getenv("VCR_MODE"); vcrEnv {
		case "RECORDING":
			vcrMode = recorder.ModeRecording
		case "REPLAYING":
			vcrMode = recorder.ModeReplaying
			// When replaying, set the poll interval low to speed up tests
			p.ProdProvider.pollInterval = 10 * time.Millisecond
		default:
			log.Printf("[DEBUG] No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", vcrEnv)
			return
		}

		envPath := os.Getenv("VCR_PATH")
		if envPath == "" {
			log.Print("[DEBUG] No environment var set for VCR_PATH, skipping VCR")
			return
		}
		path := filepath.Join(envPath, vcrFileName(p.TestName))

		rec, err := recorder.NewAsMode(path, vcrMode, p.ProdProvider.client.Transport)
		if err != nil {
			resp.Diagnostics.AddError("error creating record as new mode", err.Error())
			return
		}
		// Defines how VCR will match requests to responses.
		rec.SetMatcher(func(r *http.Request, i cassette.Request) bool {
			// Default matcher compares method and URL only
			if !cassette.DefaultMatcher(r, i) {
				return false
			}
			if r.Body == nil {
				return true
			}
			contentType := r.Header.Get("Content-Type")
			// If body contains media, don't try to compare
			if strings.Contains(contentType, "multipart/related") {
				return true
			}

			var b bytes.Buffer
			if _, err := b.ReadFrom(r.Body); err != nil {
				log.Printf("[DEBUG] Failed to read request body from cassette: %v", err)
				return false
			}
			r.Body = ioutil.NopCloser(&b)
			reqBody := b.String()
			// If body matches identically, we are done
			if reqBody == i.Body {
				return true
			}

			// JSON might be the same, but reordered. Try parsing json and comparing
			if strings.Contains(contentType, "application/json") {
				var reqJson, cassetteJson interface{}
				if err := json.Unmarshal([]byte(reqBody), &reqJson); err != nil {
					log.Printf("[DEBUG] Failed to unmarshall request json: %v", err)
					return false
				}
				if err := json.Unmarshal([]byte(i.Body), &cassetteJson); err != nil {
					log.Printf("[DEBUG] Failed to unmarshall cassette json: %v", err)
					return false
				}
				return reflect.DeepEqual(reqJson, cassetteJson)
			}
			return false
		})
		p.ProdProvider.client.Transport = rec
		configsLock.Lock()
		fwProviders[p.TestName] = p
		configsLock.Unlock()
		return
	} else {
		log.Print("[DEBUG] VCR_PATH or VCR_MODE not set, skipping VCR")
	}
}

// General test utils
func MuxedProviders(testName string) (func() tfprotov5.ProviderServer, error) {
	ctx := context.Background()

	// plugin framework provider
	downgradedFrameworkProvider, err := tf6to5server.DowngradeServer(
		context.Background(),
		providerserver.NewProtocol6(&NewFrameworkTestProvider(testName).ProdProvider),
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	providers := []func() tfprotov5.ProviderServer{
		func() tfprotov5.ProviderServer {
			return downgradedFrameworkProvider // framework provider
		},
		Provider().GRPCProvider, // sdk provider
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}

func getTestAccFrameworkProviders(testName string, c resource.TestCase) map[string]func() (tfprotov5.ProviderServer, error) {
	myFunc := func() (tfprotov5.ProviderServer, error) {
		prov, err := MuxedProviders(testName)
		return prov(), err
	}

	var testProvider string
	providerMapKeys := reflect.ValueOf(c.ProtoV5ProviderFactories).MapKeys()
	if strings.Contains(providerMapKeys[0].String(), "google-beta") {
		testProvider = "google-beta"
	} else {
		testProvider = "google"
	}
	return map[string]func() (tfprotov5.ProviderServer, error){
		testProvider: myFunc,
	}
}

func getTestFwProvider(t *testing.T) *frameworkTestProvider {
	configsLock.RLock()
	fwProvider, ok := fwProviders[t.Name()]
	configsLock.RUnlock()
	if ok {
		return fwProvider
	}

	return NewFrameworkTestProvider(t.Name())
}
