package main

import (
	"os"

	"github.com/myzie/env"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	var sizeLimit string
	flag.StringVar(&sizeLimit, "blob-size-limit", "100M", "Blob size limit")

	e := env.Must()
	e.Settings.Log()

	log.Infof("Blob size limit: %s", sizeLimit)

	newBlobsService(e, sizeLimit)
	err := e.Run()
	log.Fatal(err)
}
