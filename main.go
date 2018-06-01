package main

import (
	"github.com/myzie/base"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

func main() {

	var sizeLimit string
	flag.StringVar(&sizeLimit, "blob-size-limit", "100M", "Blob size limit")

	log.Infof("Blob size limit: %s", sizeLimit)

	service := newBlobsService(base.Must(), sizeLimit)

	if err := service.DB.AutoMigrate(Blob{}).Error; err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
