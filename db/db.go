package db

//go:generate mockgen -source=db.go -package db -destination mock.go

// Query used to list Blobs
type Query struct {
	Offset     int
	Limit      int
	OrderBy    string
	Context    string
	PathPrefix string
}

// Key used to look up a Blob
type Key struct {
	ID      string
	Context string
	Path    string
}

// Database holding Blob metadata
type Database interface {

	// Get a Blob using the key
	Get(Key) (*Blob, error)

	// Save the Blob to the Database which updates all its fields
	Save(*Blob) error

	// Delete the Blob from the Database
	Delete(*Blob) error

	// List Blobs matching the query
	List(Query) ([]*Blob, error)
}
