/*
* Copyright 2025 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cloudstorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

// ErrObjectNotExist is returned when a requested GCS object does not exist.
var ErrObjectNotExist = errors.New("storage: object doesn't exist")

func (gcs *Client) WriteToGCSBucket(bucket, object, filePath string) error {
	ctx := context.Background()

	client, err := storage.NewService(ctx, gcs.opts...)
	if err != nil {
		return fmt.Errorf("storage.NewService: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	obj := &storage.Object{
		Name:        object,
		ContentType: "application/json",
	}
	if _, err := client.Objects.Insert(bucket, obj).Media(file, googleapi.ContentType("application/json")).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Objects.Insert(%q): %w", object, err)
	}

	fmt.Printf("File uploaded to bucket %s as %s\n", bucket, object)
	return nil
}

func (gcs *Client) DownloadFile(bucket, object, filePath string) error {
	ctx := context.Background()
	client, err := storage.NewService(ctx, gcs.opts...)
	if err != nil {
		return fmt.Errorf("storage.NewService: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	resp, err := client.Objects.Get(bucket, object).Context(ctx).Download()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == http.StatusNotFound {
			return fmt.Errorf("Objects.Get(%q): %w", object, ErrObjectNotExist)
		}
		return fmt.Errorf("Objects.Get(%q): %w", object, err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("file.Close: %w", err)
	}

	fmt.Printf("Obcject %s downloaded from bucket %s as %s\n", object, bucket, filePath)

	return nil

}
