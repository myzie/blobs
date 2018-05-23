package main

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/myzie/blobs/gen/server/restapi/operations"
	log "github.com/sirupsen/logrus"
)

type service struct {
}

func (s *service) ListBlobs(params operations.ListBlobsParams) middleware.Responder {
	log.Info("ListBlobs")
	return operations.NewListBlobsOK()
}
