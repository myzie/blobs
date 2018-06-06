package main

import (
	"github.com/myzie/base"
	"github.com/myzie/blobs/db"
	"github.com/myzie/blobs/store"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

func main() {

	var sizeLimit string
	flag.StringVar(&sizeLimit, "blob-size-limit", "100M", "Blob size limit")

	log.Infof("Blob size limit: %s", sizeLimit)

	base := base.Must()

	objStoreSettings := base.Settings.ObjectStore

	objStore, err := store.NewMinioObjectStore(store.MinioOpts{
		Bucket: objStoreSettings.Bucket,
		Region: objStoreSettings.Region,
		URL:    objStoreSettings.URL,
		UseSSL: !objStoreSettings.DisableSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	blobDB := db.NewStandardDB(base.DB)

	serviceOpts := blobsServiceOpts{
		Base:      base,
		Store:     objStore,
		Database:  blobDB,
		SizeLimit: sizeLimit,
	}

	service := newBlobsService(serviceOpts)

	if err := service.DB.AutoMigrate(db.Blob{}).Error; err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
