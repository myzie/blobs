package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// BlobAttributes sent by a client
type BlobAttributes struct {
	Name       string                 `json:"name" form:"name"`
	Path       string                 `json:"path" form:"path"`
	Hash       string                 `json:"hash" form:"hash"`
	Size       int64                  `json:"size" form:"size"`
	Properties map[string]interface{} `json:"properites" form:"properties"`
}

// Normalize attributes to standard form. Especially the path format.
func (attrs *BlobAttributes) Normalize() {
	if attrs.Path == "" && attrs.Name != "" {
		attrs.Path = "/" + attrs.Name
	}
	// Normalize by trimming trailing slash; adding preceding slash
	if strings.HasSuffix(attrs.Path, "/") {
		attrs.Path = attrs.Path[:len(attrs.Path)-1]
	}
	if !strings.HasPrefix(attrs.Path, "/") {
		attrs.Path = "/" + attrs.Path
	}
}

// Validate checks whether the attributes are valid
func (attrs *BlobAttributes) Validate() error {
	if attrs.Name == "" {
		return errors.New("Name was not specified")
	}
	if strings.Contains(attrs.Name, "/") {
		return errors.New("Name must not contain '/'")
	}
	if attrs.Size < 1 {
		return errors.New("File is empty")
	}
	if attrs.Size > MaxUploadSize {
		return fmt.Errorf("File size too large: %d", attrs.Size)
	}
	if !strings.HasPrefix(attrs.Path, "/") {
		return fmt.Errorf("Invalid path: '%s'", attrs.Path)
	}
	return nil
}

// Key returns the path within the destination bucket for this upload
func (attrs *BlobAttributes) Key() string {
	return attrs.Path[1:]
}

// MarshalProperties returns the Properties field marshaled as JSON
func (attrs *BlobAttributes) MarshalProperties() ([]byte, error) {
	return json.Marshal(attrs.Properties)
}
