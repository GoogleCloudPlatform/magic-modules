/*
* Copyright 2026 Google LLC. All Rights Reserved.
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
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/api/option"
)

func TestWriteToGCSBucket(t *testing.T) {
	const (
		wantBucket  = "nightly-test-data"
		wantObject  = "test-metadata/ga/2026-11-05-ga.json"
		wantPayload = `[{"name":"TestAccExample","status":"SUCCESS"}]`
	)

	var (
		gotMethod     string
		gotPath       string
		gotMetaName   string
		gotMetaType   string
		gotPartType   string
		gotMediaBytes []byte
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path

		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Errorf("ParseMediaType(%q): %v", r.Header.Get("Content-Type"), err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if mediaType != "multipart/related" {
			t.Errorf("Content-Type mediaType = %q, want multipart/related", mediaType)
		}

		mr := multipart.NewReader(r.Body, params["boundary"])
		metaPart, err := mr.NextPart()
		if err != nil {
			t.Errorf("reading metadata part: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var meta struct {
			Name        string `json:"name"`
			ContentType string `json:"contentType"`
		}
		if err := json.NewDecoder(metaPart).Decode(&meta); err != nil {
			t.Errorf("decoding metadata part: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gotMetaName = meta.Name
		gotMetaType = meta.ContentType

		mediaPart, err := mr.NextPart()
		if err != nil {
			t.Errorf("reading media part: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gotPartType = mediaPart.Header.Get("Content-Type")
		gotMediaBytes, err = io.ReadAll(mediaPart)
		if err != nil {
			t.Errorf("reading media bytes: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"name":"`+wantObject+`"}`)
	}))
	defer srv.Close()

	srcPath := filepath.Join(t.TempDir(), "payload.json")
	if err := os.WriteFile(srcPath, []byte(wantPayload), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	client := NewClient(option.WithEndpoint(srv.URL+"/storage/v1/"), option.WithoutAuthentication())
	if err := client.WriteToGCSBucket(wantBucket, wantObject, srcPath); err != nil {
		t.Fatalf("WriteToGCSBucket returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want %q", gotMethod, http.MethodPost)
	}
	wantPath := "/upload/storage/v1/b/" + wantBucket + "/o"
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
	if gotMetaName != wantObject {
		t.Errorf("metadata name = %q, want %q", gotMetaName, wantObject)
	}
	if gotMetaType != "application/json" {
		t.Errorf("metadata contentType = %q, want application/json", gotMetaType)
	}
	if gotPartType != "application/json" {
		t.Errorf("media part Content-Type = %q, want application/json", gotPartType)
	}
	if string(gotMediaBytes) != wantPayload {
		t.Errorf("uploaded payload = %q, want %q", string(gotMediaBytes), wantPayload)
	}
}

func TestWriteToGCSBucket_MissingLocalFile(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(option.WithEndpoint(srv.URL+"/storage/v1/"), option.WithoutAuthentication())
	err := client.WriteToGCSBucket("bucket", "obj", filepath.Join(t.TempDir(), "does-not-exist.json"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got: %v", err)
	}
	if called {
		t.Errorf("expected HTTP server not to be called when local file is missing")
	}
}

func TestDownloadFile_Success(t *testing.T) {
	const (
		wantBucket  = "nightly-test-data"
		wantObject  = "test-metadata/beta/2026-11-05-beta.json"
		wantContent = `[{"name":"TestAccSample","status":"FAILURE"}]`
	)

	var (
		gotMethod string
		gotPath   string
		gotAlt    string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAlt = r.URL.Query().Get("alt")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, wantContent)
	}))
	defer srv.Close()

	destPath := filepath.Join(t.TempDir(), "downloaded.json")
	client := NewClient(option.WithEndpoint(srv.URL+"/storage/v1/"), option.WithoutAuthentication())
	if err := client.DownloadFile(wantBucket, wantObject, destPath); err != nil {
		t.Fatalf("DownloadFile returned error: %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want %q", gotMethod, http.MethodGet)
	}
	wantPath := "/storage/v1/b/" + wantBucket + "/o/" + wantObject
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
	if gotAlt != "media" {
		t.Errorf("query alt = %q, want media", gotAlt)
	}

	gotBytes, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", destPath, err)
	}
	if string(gotBytes) != wantContent {
		t.Errorf("downloaded content = %q, want %q", string(gotBytes), wantContent)
	}
}

func TestDownloadFile_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":{"code":404,"message":"No such object"}}`)
	}))
	defer srv.Close()

	destPath := filepath.Join(t.TempDir(), "should-not-be-created.json")
	client := NewClient(option.WithEndpoint(srv.URL+"/storage/v1/"), option.WithoutAuthentication())
	err := client.DownloadFile("nightly-test-data", "missing/object.json", destPath)
	if !errors.Is(err, ErrObjectNotExist) {
		t.Fatalf("DownloadFile error = %v, want errors.Is(..., ErrObjectNotExist)", err)
	}

	if _, statErr := os.Stat(destPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("destination file %q should not be created on 404, stat err = %v", destPath, statErr)
	}
}

func TestDownloadFile_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	destPath := filepath.Join(t.TempDir(), "should-not-be-created.json")
	client := NewClient(option.WithEndpoint(srv.URL+"/storage/v1/"), option.WithoutAuthentication())
	err := client.DownloadFile("nightly-test-data", "broken/object.json", destPath)
	if err == nil {
		t.Fatalf("expected non-nil error on 500 response")
	}
	if errors.Is(err, ErrObjectNotExist) {
		t.Errorf("500 error should not match ErrObjectNotExist: %v", err)
	}
	if _, statErr := os.Stat(destPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("destination file %q should not be created on 500, stat err = %v", destPath, statErr)
	}
}
