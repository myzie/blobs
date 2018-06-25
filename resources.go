package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// BlobUploadAttributes contains fields sent by a client in an upload form
type BlobUploadAttributes struct {
	Context    string                 `json:"context" form:"context"`
	Name       string                 `json:"name" form:"name"`
	Path       string                 `json:"path" form:"path"`
	Hash       string                 `json:"hash" form:"hash"`
	Size       int64                  `json:"size" form:"size"`
	Properties map[string]interface{} `json:"properties" form:"properties"`
}

// Normalize attributes to standard form. Especially the path format.
func (attrs *BlobUploadAttributes) Normalize() {
	// Normalize by trimming trailing slash; adding preceding slash
	if strings.HasSuffix(attrs.Path, "/") {
		attrs.Path = attrs.Path[:len(attrs.Path)-1]
	}
	if !strings.HasPrefix(attrs.Path, "/") {
		attrs.Path = "/" + attrs.Path
	}
}

// Validate checks whether the attributes are valid
func (attrs *BlobUploadAttributes) Validate() error {
	if attrs.Size < 0 {
		return errors.New("Invalid size: negative")
	}
	if attrs.Size > MaxUploadSize {
		return fmt.Errorf("File size too large: %d", attrs.Size)
	}
	if !strings.HasPrefix(attrs.Path, "/") {
		return fmt.Errorf("Invalid path: '%s'", attrs.Path)
	}
	if len(attrs.Properties) > 1000 {
		return fmt.Errorf("Invalid properties: too many keys")
	}
	return nil
}

// MarshalProperties returns the Properties field marshaled as JSON
func (attrs *BlobUploadAttributes) MarshalProperties() ([]byte, error) {
	return json.Marshal(attrs.Properties)
}
