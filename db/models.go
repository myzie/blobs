package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// MaxPropertiesSize specifies the max size in bytes for Blob properties
const MaxPropertiesSize = 4 * 1024

// Model base for any type saved to the database
type Model struct {
	ID         string    `gorm:"size:64;primary_key;unique_index"`
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time `gorm:"index"`
	CreatedBy  string    `gorm:"size:64;index"`
	UpdatedBy  string    `gorm:"size:64;index"`
	Name       string    `gorm:"size:128"`
	Properties postgres.Jsonb
}

// Blob is a stored object
type Blob struct {
	Model
	Context string `gorm:"size:64;index"`
	Path    string `gorm:"size:256;unique_index"`
	Size    int64
}

// User information including email
type User struct {
	Model
	Email string `gorm:"size:320"`
}

// Group of Users
type Group struct {
	Model
}

// Context is a logical container for Blobs
type Context struct {
	Model
}

// Key used when storing the blob
func (b *Blob) Key() string {
	return fmt.Sprintf("%s/%s", b.Context, b.Path)
}

// BeforeSave is called to validate the model before saving to the database
func (b *Blob) BeforeSave() error {

	if b.ID == "" {
		return fmt.Errorf("Invalid id: empty")
	}
	if b.Context == "" {
		return fmt.Errorf("Invalid context: empty")
	}
	if b.CreatedBy == "" {
		return fmt.Errorf("Invalid created_at: empty")
	}
	if b.UpdatedBy == "" {
		return fmt.Errorf("Invalid updated_at: empty")
	}
	if b.Size < 0 {
		return fmt.Errorf("Invalid size: negative")
	}

	if len(b.Name) > 100 {
		return fmt.Errorf("Invalid name: too long")
	}

	// TODO: path regex?
	if len(b.Path) == 0 {
		return fmt.Errorf("Invalid path: empty")
	} else if len(b.Path) > 200 {
		return fmt.Errorf("Invalid path: too long")
	} else if b.Path[0] != '/' {
		return fmt.Errorf("Invalid path: does not start with /")
	}

	if len(b.Properties.RawMessage) > MaxPropertiesSize {
		return fmt.Errorf("Invalid properties: too large")
	}
	return nil
}
