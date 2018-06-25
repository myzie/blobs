package main

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path"

	jwt "github.com/dgrijalva/jwt-go"
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
	group.POST("", svc.Post)
	group.DELETE("/*", svc.Delete)

	return svc
}

func (svc *blobsService) Get(c echo.Context) error {

	// Access context from either:
	// A) the JWT
	// B) an API key mapping
	// C) the query string

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
		return c.JSON(OK, newBlobView(blob))
	}

	// Otherwise return the object itself
	getOpts := minio.GetObjectOptions{}
	obj, err := svc.Store.Get(blob.Key(), getOpts)
	if err != nil {
		return c.JSON(InternalServerError, errorView{"Failed to get object"})
	}
	return c.Stream(OK, "application/octet-stream", obj)
}

func (svc *blobsService) Post(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*base.JWTClaims)
	userID := claims.Subject

	var attrs BlobUploadAttributes
	if err := c.Bind(&attrs); err != nil {
		return c.JSON(BadRequest, errorView{"Failed to bind attributes"})
	}
	attrs.Normalize()
	if err := attrs.Validate(); err != nil {
		return c.JSON(BadRequest, errorView{err.Error()})
	}
	propJSON, err := attrs.MarshalProperties()
	if err != nil {
		return c.JSON(BadRequest, errorView{"JSON error"})
	}
	pgrsJSON := postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}

	// Determine if a blob already exists at that path
	blob, err := svc.Database.Get(attrs.Path)
	if err != nil {
		if err.Error() != "record not found" {
			log.WithError(err).Error("Get failed")
			return c.JSON(InternalServerError, errorView{"Blob lookup failed"})
		}

		blob = &db.Blob{
			ID:         uid(),
			CreatedBy:  userID,
			UpdatedBy:  userID,
			Path:       attrs.Path,
			Size:       attrs.Size,
			Properties: pgrsJSON,
		}

		log.WithFields(log.Fields{
			"id":         blob.ID,
			"created_by": blob.CreatedBy,
			"updated_by": blob.UpdatedBy,
			"path":       blob.Path,
			"size":       attrs.Size,
		}).Info("Creating blob")

		if err := svc.Database.Save(blob); err != nil {
			log.WithError(err).Error("Save failed")
			return c.JSON(InternalServerError, errorView{"Save failed"})
		}
	} else {
		blob.Size = attrs.Size
		blob.Properties = pgrsJSON
		blob.UpdatedBy = userID

		log.WithFields(log.Fields{
			"id":         blob.ID,
			"created_by": blob.CreatedBy,
			"updated_by": blob.UpdatedBy,
			"path":       blob.Path,
			"size":       attrs.Size,
		}).Info("Updating blob")

		fields := []string{"properties", "size", "updated_by"}
		if err := svc.Database.Update(blob, fields); err != nil {
			log.WithError(err).Error("Update failed")
			return c.JSON(InternalServerError, errorView{"Update failed"})
		}
	}

	// TODO: allow no file to be attached

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(BadRequest, errorView{"Form file missing"})
	}

	return c.JSON(OK, newBlobView(blob))
}

func (svc *blobsService) saveFile(c echo.Context, blob *models.Blob, hdr *multipart.FileHeader) error {

	src, err := file.Open()
	if err != nil {
		return c.JSON(BadRequest, errorView{"Form file error"})
	}
	defer src.Close()

	log.WithFields(log.Fields{
		"id":   blob.ID,
		"key":  blob.Key(),
		"size": attrs.Size,
	}).Info("Upload starting")

	fileName := path.Base(attrs.Path)

	opts := minio.PutObjectOptions{
		ContentType:        "application/octet-stream",
		ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, fileName),
		UserMetadata: map[string]string{
			"id":   blob.ID,
			"path": blob.Path,
		},
	}
	n, err := svc.Store.Put(blob.Key(), src, attrs.Size, opts)
	if err != nil {
		log.WithError(err).Error("Error saving file to bucket")
		return c.JSON(InternalServerError, errorView{"Error saving file"})
	}
	if n != attrs.Size {
		log.Error("Uploaded file size incorrect")
		return c.JSON(InternalServerError, errorView{"Error saving file"})
	}

	log.WithFields(log.Fields{
		"id":   blob.ID,
		"key":  blob.Key(),
		"size": attrs.Size,
	}).Info("Upload complete")
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
