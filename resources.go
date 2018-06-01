package main

// BlobAttributes sent by a client
type BlobAttributes struct {
	Name       string                 `json:"name" form:"name"`
	Extension  string                 `json:"extension" form:"extension"`
	Path       string                 `json:"path" form:"path"`
	Hash       string                 `json:"hash" form:"hash"`
	Properties map[string]interface{} `json:"properites" form:"properties"`
}
