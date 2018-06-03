package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	minio "github.com/minio/minio-go"
	"github.com/myzie/base"
	"github.com/myzie/blobs/store"
	log "github.com/sirupsen/logrus"
)

// MaxUploadSize defines the max file size in bytes for uploads
const MaxUploadSize = 100 * 1024 * 1024

type blobsService struct {
	*base.Base
	Store store.ObjectStore
}

// newBlobsService returns an HTTP interface for blobs
func newBlobsService(base *base.Base, store store.ObjectStore, sizeLimit string) *blobsService {
	svc := &blobsService{Base: base, Store: store}
	group := svc.Echo.Group("/blobs")
	group.Use(middleware.BodyLimit(sizeLimit))
	group.GET("", svc.List)
	group.GET("/*", svc.Get)
	group.PUT("/*", svc.Put)
	group.POST("", svc.Post)
	group.DELETE("/*", svc.Delete)
	return svc
}

func (svc *blobsService) getBlobWithPath(path string) (*Blob, error) {
	blob := &Blob{}
	if err := svc.DB.Where("path = ?", path).First(blob).Error; err != nil {
		return nil, err
	}
	return blob, nil
}

func (svc *blobsService) updateBlob(blob *Blob) error {
	fields := []string{"name", "extension", "properties"}
	return svc.DB.Model(blob).Select(fields).Updates(blob).Error
}

func (svc *blobsService) saveBlob(blob *Blob) error {
	return svc.DB.Save(blob).Error
}

func (svc *blobsService) Get(c echo.Context) error {

	// Look up blob at the specified path
	path := "/" + c.ParamValues()[0]
	blob, err := svc.getBlobWithPath(path)
	if err != nil {
		log.Infof("404 error: %+v", reflect.TypeOf(err))
		if err.Error() == "record not found" {
			return c.JSON(NotFound, errorView{"Blob not found"})
		}
		log.WithError(err).Error("Get failed")
		return c.JSON(InternalServerError, errorView{"Failed to look up Blob"})
	}

	// Return the blob metadata if JSON content was requested
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		return c.JSON(OK, blob)
	}

	// Otherwise return the object itself
	getOpts := minio.GetObjectOptions{}
	obj, err := svc.Store.Get(blob.Key(), getOpts)
	if err != nil {
		return c.JSON(InternalServerError, errorView{"Failed to get object"})
	}
	return c.Stream(OK, "application/octet-stream", obj)
}

func (svc *blobsService) Put(c echo.Context) error {

	path := "/" + c.ParamValues()[0]

	var attrs BlobUpdateAttributes
	if err := c.Bind(&attrs); err != nil {
		return c.JSON(BadRequest, errorView{"Bad attributes"})
	}
	if err := attrs.Validate(); err != nil {
		return c.JSON(BadRequest, errorView{err.Error()})
	}

	blob, err := svc.getBlobWithPath(path)
	if err != nil {
		return c.JSON(NotFound, errorView{"Not found"})
	}
	propJSON, err := attrs.MarshalProperties()
	if err != nil {
		return c.JSON(BadRequest, errorView{"Bad properties"})
	}
	log.Infof("Blob attributes: %+v properties: %s", attrs, string(propJSON))

	blob.Name = attrs.Name
	blob.Properties = postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}

	fields := []string{"name", "properties"}
	if err := svc.DB.Model(blob).Select(fields).Updates(blob).Error; err != nil {
		log.WithError(err).Error("Failed to update blob")
		return c.JSON(InternalServerError, errorView{"Failed to update blob"})
	}
	return c.JSON(OK, blob)
}

func (svc *blobsService) Post(c echo.Context) error {

	var attrs BlobUploadAttributes
	if err := c.Bind(&attrs); err != nil {
		return c.JSON(BadRequest, errorView{"Failed to bind attributes"})
	}

	attrs.Normalize()
	if err := attrs.Validate(); err != nil {
		return c.JSON(BadRequest, errorView{err.Error()})
	}

	blobExt := attrs.Extension()
	propJSON, err := attrs.MarshalProperties()
	if err != nil {
		return c.JSON(BadRequest, errorView{"JSON error"})
	}

	// Determine if a blob already exists at that path
	blob, err := svc.getBlobWithPath(attrs.Path)
	if err != nil {
		if err.Error() != "record not found" {
			log.WithError(err).Error("Get failed")
			return c.JSON(InternalServerError, errorView{"Blob lookup failed"})
		}
	}

	// Create or update the blob
	if blob == nil {

		blob = &Blob{
			ID:         uid(),
			Name:       attrs.Name,
			Extension:  blobExt,
			Path:       attrs.Path,
			Hash:       "", // TODO
			Properties: postgres.Jsonb{RawMessage: json.RawMessage(propJSON)},
		}

		log.WithFields(log.Fields{
			"id":   blob.ID,
			"name": blob.Name,
			"path": blob.Path,
		}).Info("Creating blob")

		if err := svc.saveBlob(blob); err != nil {
			log.WithError(err).Error("Save failed")
			return c.JSON(InternalServerError, errorView{"Save failed"})
		}
	} else {

		blob.Name = attrs.Name
		blob.Properties = postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}

		log.WithFields(log.Fields{
			"id":   blob.ID,
			"name": blob.Name,
			"path": blob.Path,
		}).Info("Updating blob")

		if err := svc.updateBlob(blob); err != nil {
			log.WithError(err).Error("Update failed")
			return c.JSON(InternalServerError, errorView{"Update failed"})
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(BadRequest, errorView{"Form file missing"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(BadRequest, errorView{"Form file error"})
	}
	defer src.Close()

	reader := io.LimitReader(src, attrs.Size)

	log.WithFields(log.Fields{
		"name":       attrs.Name,
		"extension":  blobExt,
		"size":       attrs.Size,
		"properties": attrs.Properties,
		"key":        blob.Key(),
	}).Info("Upload starting")

	opts := minio.PutObjectOptions{
		ContentType:        "application/octet-stream",
		ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, attrs.Name),
	}
	n, err := svc.Store.Put(blob.Key(), reader, attrs.Size, opts)
	if err != nil {
		log.WithError(err).Error("Error saving file to bucket")
		return c.JSON(InternalServerError, errorView{"Error saving file"})
	}
	if n != attrs.Size {
		log.Error("Uploaded file size incorrect")
		return c.JSON(InternalServerError, errorView{"Error saving file"})
	}

	log.WithFields(log.Fields{
		"key":      attrs.Key(),
		"size":     attrs.Size,
		"filename": attrs.Name,
	}).Info("Upload complete")

	return c.JSON(OK, errorView{})
}

func (svc *blobsService) Delete(c echo.Context) error {

	// Look up blob at the specified path
	path := "/" + c.ParamValues()[0]
	blob := &Blob{}
	if err := svc.DB.Where("path = ?", path).First(blob).Error; err != nil {
		return c.JSON(NotFound, errorView{"Blob not found"})
	}

	// Remove object from S3
	if err := svc.Store.Remove(blob.Key()); err != nil {
		log.WithError(err).Error("Failed to delete object")
		return c.JSON(InternalServerError, errorView{"Failed to delete object"})
	}

	// Remove database entry
	if err := svc.DB.Delete(blob).Error; err != nil {
		log.WithError(err).Error("Failed to delete blob")
		return c.JSON(InternalServerError, errorView{"Failed to delete blob"})
	}

	log.WithFields(log.Fields{"id": blob.ID, "path": path, "name": blob.Name}).
		Info("Blob deleted")

	return c.NoContent(NoContent)
}

type listQueryParameters struct {
	Offset  int    `query:"offset"`
	Limit   int    `query:"limit"`
	OrderBy string `query:"order_by"`
}

func (svc *blobsService) List(c echo.Context) error {

	var params listQueryParameters
	if err := c.Bind(&params); err != nil {
		return c.JSON(BadRequest, errorView{"Bad parameters"})
	}

	if params.Offset < 0 {
		return c.JSON(BadRequest, errorView{"Invalid offset"})
	}
	if params.Limit < 0 {
		return c.JSON(BadRequest, errorView{"Invalid limit"})
	}
	if params.Limit == 0 {
		params.Limit = 1000
	}
	if params.OrderBy == "" {
		params.OrderBy = "path"
	}

	var blobs []*Blob

	err := svc.DB.Order(params.OrderBy).
		Offset(params.Offset).
		Limit(params.Limit).
		Find(&blobs).Error

	if err != nil {
		log.WithError(err).Error("Blob query failed")
		return c.JSON(InternalServerError, errorView{"Blob query failed"})
	}
	return c.JSON(OK, blobs)
}
