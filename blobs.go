package main

import (
	"encoding/json"
	"io"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/myzie/base"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type blobsService struct {
	*base.Base
}

// newBlobsService returns an HTTP interface for blobs
func newBlobsService(base *base.Base, sizeLimit string) *blobsService {
	svc := &blobsService{Base: base}
	group := svc.Echo.Group("/blobs")
	group.Use(middleware.BodyLimit(sizeLimit))
	group.GET("", svc.List)
	group.GET("/*", svc.Get)
	group.PUT("/*", svc.Put)
	group.POST("", svc.Post)
	group.DELETE("/*", svc.Delete)
	return svc
}

func (svc *blobsService) Get(c echo.Context) error {

	path := "/" + c.ParamValues()[0]
	blob := &Blob{}

	if err := svc.DB.Where("path = ?", path).First(blob).Error; err != nil {
		if err.Error() == "record not found" {
			return c.JSON(NotFound, errorView{"Blob not found"})
		}
		log.WithError(err).Error("Get failed")
		return c.JSON(InternalServerError, errorView{"Blob not found"})
	}
	return c.JSON(OK, blob)
}

func (svc *blobsService) Put(c echo.Context) error {

	var attrs BlobAttributes
	if err := c.Bind(&attrs); err != nil {
		return c.JSON(BadRequest, errorView{"Failed to bind attributes"})
	}

	propJSON, err := json.Marshal(attrs.Properties)
	if err != nil {
		return c.JSON(BadRequest, errorView{"Bad properties"})
	}
	log.Infof("Blob attributes: %+v properties: %s", attrs, string(propJSON))

	b := &Blob{
		Name:       attrs.Name,
		Path:       attrs.Path,
		Extension:  attrs.Extension,
		Properties: postgres.Jsonb{RawMessage: json.RawMessage(propJSON)},
	}
	b.ID = uuid.Must(uuid.NewV4()).String()

	if err := svc.DB.Save(b).Error; err != nil {
		log.WithError(err).Error("Failed to save blob")
		return c.JSON(InternalServerError, errorView{"Failed to update blob"})
	}
	return c.JSON(OK, b)
}

func (svc *blobsService) Post(c echo.Context) error {
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	return c.JSON(BadRequest, errorView{"Form file missing"})
	// }
	// src, err := file.Open()
	// if err != nil {
	// 	return c.JSON(BadRequest, errorView{"Form file error"})
	// }
	// defer src.Close()
	// return svc.soundUpload(c, src)
	return nil
}

func (svc *blobsService) Delete(c echo.Context) error {
	// soundPath := "/" + c.ParamValues()[0]
	// s, err := svc.App.SoundStore().Get(app.SoundKey{Path: soundPath})
	// if err != nil {
	// 	return c.JSON(NotFound, errorView{"Sound not found"})
	// }
	// if err = svc.App.SoundStore().Delete(s); err != nil {
	// 	log.WithError(err).Error("Sound delete failed")
	// 	return c.JSON(InternalServerError, errorView{"Sound delete failed"})
	// }
	return c.NoContent(NoContent)
}

func (svc *blobsService) List(c echo.Context) error {
	// query := app.SoundQuery{
	// 	Offset:    0,
	// 	Limit:     100,
	// 	PathBegin: "",
	// 	PathEnd:   "",
	// }
	// sounds, err := svc.App.SoundStore().List(query)
	// if err != nil {
	// 	log.WithError(err).Error("List sounds failed")
	// 	return c.JSON(InternalServerError, errorView{"List sounds failed"})
	// }
	// views := make([]*SoundView, len(sounds))
	// for i, sound := range sounds {
	// 	views[i] = soundView(sound)
	// }
	// return c.JSON(OK, views)
	return nil
}

func (svc *blobsService) soundUpload(c echo.Context, fileReader io.Reader) error {

	// attrs, err := svc.bindSoundAttrs(c)
	// if err != nil {
	// 	log.WithError(err).Error("Attributes error")
	// 	return c.JSON(BadRequest, errorView{"Attributes error"})
	// }

	// snd, err := svc.App.SoundStore().Get(app.SoundKey{Path: attrs.Path})
	// if err != nil && err.Error() != "record not found" {
	// 	log.WithError(err).Error("Failed to lookup sound")
	// 	return c.JSON(InternalServerError, errorView{"Failed to lookup sound"})
	// }

	// props, err := attrs.MarshalProperties()
	// if err != nil {
	// 	log.WithError(err).Error("Bad properties JSON")
	// 	return c.JSON(BadRequest, errorView{"Bad properties JSON"})
	// }

	// successCode := OK

	// if snd == nil {
	// 	// Create sound
	// 	snd = &app.Sound{
	// 		Name:        attrs.Name,
	// 		Path:        attrs.Path,
	// 		Description: attrs.Description,
	// 		Transcript:  attrs.Transcript,
	// 		Properties:  postgres.Jsonb{RawMessage: json.RawMessage(props)},
	// 	}
	// 	if err := svc.App.SoundStore().Create(snd); err != nil {
	// 		log.WithError(err).Error("Failed to create sound")
	// 		return c.JSON(InternalServerError, errorView{"Failed to create sound"})
	// 	}
	// 	successCode = Created
	// } else {
	// 	// Update sound
	// 	snd.Description = attrs.Description
	// 	snd.Transcript = attrs.Transcript
	// 	snd.Properties = postgres.Jsonb{RawMessage: json.RawMessage(props)}
	// 	if err := svc.App.SoundStore().Save(snd); err != nil {
	// 		log.WithError(err).Error("Failed to update sound")
	// 		return c.JSON(InternalServerError, errorView{"Failed to update sound"})
	// 	}
	// }

	// // Upload sound file to S3 bucket
	// if !attrs.JSONOnly && fileReader != nil {
	// 	if attrs.Size <= 0 {
	// 		return c.JSON(BadRequest, errorView{"Error detecting sound size"})
	// 	}
	// 	reader := io.LimitReader(fileReader, attrs.Size)
	// 	if err := svc.App.UploadSound(reader, attrs); err != nil {
	// 		log.WithError(err).Error("Error saving file")
	// 		return c.JSON(InternalServerError, errorView{"Error saving file"})
	// 	}
	// 	// TODO: retrieve ETag from S3?
	// }
	// return c.JSON(successCode, soundView(snd))
	return nil
}

// func (svc *blobsService) bindSoundAttrs(c echo.Context) (*app.UploadAttributes, error) {

// 	var err error
// 	attrs := &app.UploadAttributes{}
// 	method := c.Request().Method

// 	// Bind handles forms via POST and JSON via PUT
// 	if err := c.Bind(attrs); err != nil {
// 		if method == "POST" {
// 			return nil, fmt.Errorf("Invalid form: %s", err.Error())
// 		}
// 	} else if method == "PUT" {
// 		// If Bind worked on a PUT, it must've been Content-Type
// 		// application/json. So there isn't a sound file as well.
// 		attrs.JSONOnly = true
// 	}

// 	// PUT requests provide sound path and name in the URL
// 	if method == "PUT" {
// 		attrs.Path = c.ParamValues()[0]
// 		parts := strings.Split(attrs.Path, "/")
// 		attrs.Name = parts[len(parts)-1]
// 	}

// 	// Name and size headers have highest precedence
// 	name := c.Request().Header.Get("Sound-Name")
// 	if name != "" {
// 		attrs.Name = name
// 	}
// 	sizeStr := c.Request().Header.Get("Sound-Size")
// 	if sizeStr != "" {
// 		attrs.Size, err = strconv.ParseInt(sizeStr, 10, 64)
// 		if err != nil {
// 			return nil, errors.New("Invalid Sound-Size header")
// 		}
// 	}
// 	attrs.Normalize()

// 	// log.WithFields(log.Fields{
// 	// 	"name":        attrs.Name,
// 	// 	"path":        attrs.Path,
// 	// 	"size":        attrs.Size,
// 	// 	"transcript":  attrs.Transcript,
// 	// 	"description": attrs.Description,
// 	// 	"tags":        attrs.Tags,
// 	// 	"properties":  attrs.Properties,
// 	// }).Info("Sound parameters")

// 	return attrs, nil
// }
