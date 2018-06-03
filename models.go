package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// MaxPropertiesSize specifies the max size in bytes for Blob properties
const MaxPropertiesSize = 10 * 1024

var nameRegex = regexp.MustCompile(`^[0-9A-Za-z_][A-Za-z0-9-_ ]*(\.[a-zA-Z0-9]+)?$`)

// Blob is a stored object
type Blob struct {
	ID         string    `gorm:"size:50;primary_key;unique_index"`
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time `gorm:"index"`
	CreatedBy  string    `gorm:"size:50;index"`
	UpdatedBy  string    `gorm:"size:50;index"`
	Name       string    `gorm:"size:100"`
	Extension  string    `gorm:"size:100"`
	Path       string    `gorm:"size:250;unique_index"`
	Hash       string    `gorm:"size:250"`
	Properties postgres.Jsonb
}

// Key used when storing the blob
func (b *Blob) Key() string {
	return fmt.Sprintf("%s/%s%s", b.ID, b.ID, b.Extension)
}

// Validate the blob
func (b *Blob) Validate() []error {

	var errs []error

	fail := func(msg string) {
		errs = append(errs, errors.New(msg))
	}

	// Name validation
	if len(b.Name) == 0 {
		fail("Invalid name: empty")
	} else if len(b.Name) > 100 {
		fail("Invalid name: too long")
	} else if !nameRegex.MatchString(b.Name) {
		fail("Invalid name: unsupported characters")
	}

	// Path validation
	if len(b.Path) == 0 {
		fail("Invalid path: empty")
	} else if len(b.Path) > 200 {
		fail("Invalid path: too long")
	} else if b.Path[0] != '/' {
		fail("Invalid path: does not start with /")
	}

	propBytes := []byte(b.Properties.RawMessage)
	if len(propBytes) > MaxPropertiesSize {
		fail("Invalid properties: too large")
	}

	return errs
}

// BeforeSave is called as the Blob is being saved to the database
func (b *Blob) BeforeSave() error {
	errs := b.Validate()
	if errs != nil {
		return errs[0]
	}
	return nil
}
