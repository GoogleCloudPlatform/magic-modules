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
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func (gcs *Client) WriteToGCSBucket(bucket, object, filePath string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	writer := client.Bucket(bucket).Object(object).NewWriter(ctx)
	writer.ContentType = "application/json"

	if _, err = io.Copy(writer, file); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Printf("File uploaded to bucket %s as %s\n", bucket, object)
	return nil
}

func (gcs *Client) DownloadFile(bucket, object, filePath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	reader, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer reader.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("file.Close: %w", err)
	}

	fmt.Printf("Obcject %s downloaded from bucket %s as %s\n", object, bucket, filePath)

	return nil

}
