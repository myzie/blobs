package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// MaxPropertiesSize specifies the max size in bytes for Blob properties
const MaxPropertiesSize = 4 * 1024

// Blob is a stored object
type Blob struct {
	ID         string    `gorm:"size:64;primary_key;unique_index"`
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time `gorm:"index"`
	CreatedBy  string    `gorm:"size:64;index"`
	UpdatedBy  string    `gorm:"size:64;index"`
	Context    string    `gorm:"size:64;index"`
	Name       string    `gorm:"size:128"`
	Path       string    `gorm:"size:256;unique_index"`
	Size       int64
	Properties postgres.Jsonb
}

// Key used when storing the blob
func (b *Blob) Key() string {
	return fmt.Sprintf("%s/%s", b.Context, b.Path)
}

// Validate the blob
func (b *Blob) Validate() []error {

	var errs []error

	fail := func(msg string) {
		errs = append(errs, errors.New(msg))
	}

	if b.ID == "" {
		fail("Invalid id: empty")
	}
	if b.Context == "" {
		fail("Invalid context: empty")
	}
	if b.CreatedBy == "" {
		fail("Invalid created_at: empty")
	}
	if b.UpdatedBy == "" {
		fail("Invalid updated_at: empty")
	}
	if b.Size < 0 {
		fail("Invalid size: negative")
	}

	if len(b.Name) > 100 {
		fail("Invalid name: too long")
	}

	// TODO: path regex?
	if len(b.Path) == 0 {
		fail("Invalid path: empty")
	} else if len(b.Path) > 200 {
		fail("Invalid path: too long")
	} else if b.Path[0] != '/' {
		fail("Invalid path: does not start with /")
	}

	if len(b.Properties.RawMessage) > MaxPropertiesSize {
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
