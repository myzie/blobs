package main

import (
	"os"

	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"

	"github.com/myzie/blobs/gen/server/restapi"
	"github.com/myzie/blobs/gen/server/restapi/operations"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	opts := struct {
		Host string `long:"host" default:"127.0.0.1" description:"Host address to use for listening"`
		Port int    `long:"port" default:"8080" description:"Port to use for listening"`
	}{}

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	spec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewBlobsAPI(spec)
	api.Logger = log.Infof
	server := restapi.NewServer(api)
	server.Host = opts.Host
	server.Port = opts.Port
	defer server.Shutdown()

	service := &service{}
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.ListBlobsHandler = operations.ListBlobsHandlerFunc(service.ListBlobs)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
