package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type standardDB struct {
	gormDB *gorm.DB
}

// NewStandardDB returns an interface to a Blob Database
func NewStandardDB(gormDB *gorm.DB) Database {
	return &standardDB{gormDB: gormDB}
}

// Get a Blob using the provided key
func (db *standardDB) Get(key Key) (*Blob, error) {
	if key.ID == "" && (key.Path == "" || key.Context == "") {
		return nil, errors.New("Invalid key")
	}
	where := &Blob{
		Model:   Model{ID: key.ID},
		Context: key.Context,
		Path:    key.Path,
	}
	blob := &Blob{}
	if err := db.gormDB.Where(where).First(blob).Error; err != nil {
		return nil, err
	}
	return blob, nil
}

// Save the Blob to the Database which updates all its fields
func (db *standardDB) Save(blob *Blob) error {
	if blob.ID == "" {
		return errors.New("Invalid empty ID")
	}
	return db.gormDB.Save(blob).Error
}

// Delete the Blob from the Database
func (db *standardDB) Delete(blob *Blob) error {
	if blob.ID == "" {
		return errors.New("Invalid empty ID")
	}
	return db.gormDB.Delete(blob).Error
}

// List Blobs matching the query
func (db *standardDB) List(q Query) ([]*Blob, error) {

	var blobs []*Blob

	err := db.gormDB.
		Order(q.OrderBy).
		Offset(q.Offset).
		Limit(q.Limit).
		Find(&blobs).Error

	if err != nil {
		return nil, err
	}
	return blobs, nil
}
