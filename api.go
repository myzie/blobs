package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"path"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	minio "github.com/minio/minio-go"
	"github.com/myzie/base"
	"github.com/myzie/blobs/db"
	"github.com/myzie/blobs/store"
	log "github.com/sirupsen/logrus"
)

// MaxUploadSize defines the max file size in bytes for uploads
const MaxUploadSize = 100 * 1024 * 1024

type blobsServiceOpts struct {
	Base      *base.Base
	Store     store.ObjectStore
	Database  db.Database
	SizeLimit string
}

type blobsService struct {
	*base.Base
	Store    store.ObjectStore
	Database db.Database
}

// newBlobsService returns an HTTP interface for blobs
func newBlobsService(opts blobsServiceOpts) *blobsService {

	svc := &blobsService{
		Base:     opts.Base,
		Store:    opts.Store,
		Database: opts.Database,
	}

	group := svc.Echo.Group("/blobs")
	group.Use(svc.JWTMiddleware())
	group.Use(middleware.BodyLimit(opts.SizeLimit))
	group.GET("", svc.List)
	group.GET("/*", svc.Get)
	group.PUT("/*", svc.Put)
	group.POST("", svc.Post)
	group.DELETE("/*", svc.Delete)

	return svc
}

func (svc *blobsService) Get(c echo.Context) error {

	// Look up blob at the specified path
	path := "/" + c.ParamValues()[0]
	blob, err := svc.Database.Get(path)
	if err != nil {
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

	// Require application/json
	contentType := c.Request().Header.Get("Content-Type")
	if contentType != "application/json" {
		return c.JSON(BadRequest, errorView{"Only JSON is accepted"})
	}

	// Reject request if item does not exist
	path := "/" + c.ParamValues()[0]
	blob, err := svc.Database.Get(path)
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(NotFound, errorView{"Blob not found"})
		}
		log.WithError(err).Error("Get failed")
		return c.JSON(InternalServerError, errorView{"Failed to look up Blob"})
	}

	var props BlobProperties
	if err := c.Bind(&props); err != nil {
		return c.JSON(BadRequest, errorView{"Bad properties"})
	}
	propJSON, err := json.Marshal(props.Properties)
	if err != nil {
		return c.JSON(BadRequest, errorView{"Bad properties"})
	}
	if len(propJSON) > db.MaxPropertiesSize {
		return c.JSON(BadRequest, errorView{"Properties too large"})
	}

	blob.Properties = postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}

	fields := []string{"properties"}
	if err := svc.Database.Update(blob, fields); err != nil {
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
	blob, err := svc.Database.Get(attrs.Path)
	if err != nil {
		if err.Error() != "record not found" {
			log.WithError(err).Error("Get failed")
			return c.JSON(InternalServerError, errorView{"Blob lookup failed"})
		}
	}

	// Create or update the blob
	if blob == nil {

		blob = &db.Blob{
			ID:         uid(),
			Path:       attrs.Path,
			Hash:       "", // TODO
			Properties: postgres.Jsonb{RawMessage: json.RawMessage(propJSON)},
		}

		log.WithFields(log.Fields{
			"id":   blob.ID,
			"path": blob.Path,
		}).Info("Creating blob")

		if err := svc.Database.Save(blob); err != nil {
			log.WithError(err).Error("Save failed")
			return c.JSON(InternalServerError, errorView{"Save failed"})
		}
	} else {

		blob.Properties = postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}

		log.WithFields(log.Fields{
			"id":   blob.ID,
			"path": blob.Path,
		}).Info("Updating blob")

		fields := []string{"name", "extension", "properties"}

		if err := svc.Database.Update(blob, fields); err != nil {
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

	sha256 := sha256.New()
	reader := io.TeeReader(io.LimitReader(src, attrs.Size), sha256)

	log.WithFields(log.Fields{
		"extension":  blobExt,
		"size":       attrs.Size,
		"properties": attrs.Properties,
		"key":        blob.Key(),
	}).Info("Upload starting")

	fileName := path.Base(attrs.Path)

	opts := minio.PutObjectOptions{
		ContentType:        "application/octet-stream",
		ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, fileName),
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

	sha256Hash := fmt.Sprintf("%x", sha256.Sum(nil))

	log.WithFields(log.Fields{
		"key":    attrs.Key(),
		"size":   attrs.Size,
		"sha256": sha256Hash,
	}).Info("Upload complete")

	return c.JSON(OK, errorView{})
}

func (svc *blobsService) Delete(c echo.Context) error {

	// Look up blob at the specified path
	path := "/" + c.ParamValues()[0]
	blob, err := svc.Database.Get(path)
	if err != nil {
		return c.JSON(NotFound, errorView{"Blob not found"})
	}

	// Remove object from S3
	if err := svc.Store.Remove(blob.Key()); err != nil {
		log.WithError(err).Error("Failed to delete object")
		return c.JSON(InternalServerError, errorView{"Failed to delete object"})
	}

	// Remove database entry
	if err := svc.Database.Delete(blob); err != nil {
		log.WithError(err).Error("Failed to delete blob")
		return c.JSON(InternalServerError, errorView{"Failed to delete blob"})
	}

	log.WithFields(log.Fields{"id": blob.ID, "path": path}).
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

	query := db.Query{
		Offset:  params.Offset,
		Limit:   params.Limit,
		OrderBy: params.OrderBy,
	}

	blobs, err := svc.Database.List(query)
	if err != nil {
		log.WithError(err).Error("Blob query failed")
		return c.JSON(InternalServerError, errorView{"Blob query failed"})
	}
	return c.JSON(OK, blobs)
}
