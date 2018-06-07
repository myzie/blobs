package db

//go:generate mockgen -source=db.go -package db -destination mock.go

// Query used to list Blobs
type Query struct {
	Offset  int
	Limit   int
	OrderBy string
}

// Database holding Blob metadata
type Database interface {

	// Get a Blob with the given path
	Get(path string) (*Blob, error)

	// Save the Blob to the Database which updates all its fields
	Save(*Blob) error

	// Update the specified Blob fields
	Update(*Blob, []string) error

	// List Blobs matching the query
	List(Query) ([]*Blob, error)

	// Delete the Blob from the Database
	Delete(*Blob) error
}
