package db

import (
	"github.com/jinzhu/gorm"
)

type standardDB struct {
	gormDB *gorm.DB
}

// NewStandardDB returns an interface to a Blob Database
func NewStandardDB(gormDB *gorm.DB) Database {
	return &standardDB{gormDB: gormDB}
}

// Get a Blob with the given path
func (db *standardDB) Get(path string) (*Blob, error) {
	blob := &Blob{}
	err := db.gormDB.Where("path = ?", path).First(blob).Error
	if err != nil {
		return nil, err
	}
	return blob, nil
}

// Save the Blob to the Database which updates all its fields
func (db *standardDB) Save(blob *Blob) error {
	return db.gormDB.Save(blob).Error
}

// Update the specified Blob fields
func (db *standardDB) Update(blob *Blob, fields []string) error {
	return db.gormDB.Model(blob).Select(fields).Updates(blob).Error
}

// List Blobs matching the query
func (db *standardDB) List(q Query) ([]*Blob, error) {

	context := &Blob{Context: q.Context}

	var blobs []*Blob

	err := db.gormDB.Where(context).
		Order(q.OrderBy).
		Offset(q.Offset).
		Limit(q.Limit).
		Find(&blobs).Error

	if err != nil {
		return nil, err
	}
	return blobs, nil
}

// Delete the Blob from the Database
func (db *standardDB) Delete(blob *Blob) error {
	return db.gormDB.Delete(blob).Error
}
