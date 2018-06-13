package main

import (
	"encoding/json"
	"time"

	"github.com/myzie/blobs/db"
)

type errorView struct {
	Error string `json:"error"`
}

type blobView struct {
	ID         string          `json:"id"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	CreatedBy  string          `json:"created_by"`
	UpdatedBy  string          `json:"updated_by"`
	Path       string          `json:"path"`
	Size       int64           `json:"size"`
	Properties json.RawMessage `json:"properties"`
}

func newBlobView(blob *db.Blob) *blobView {
	return &blobView{
		ID:         blob.ID,
		CreatedAt:  blob.CreatedAt,
		CreatedBy:  blob.CreatedBy,
		UpdatedAt:  blob.UpdatedAt,
		UpdatedBy:  blob.UpdatedBy,
		Path:       blob.Path,
		Size:       blob.Size,
		Properties: blob.Properties.RawMessage,
	}
}
