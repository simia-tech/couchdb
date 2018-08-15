package couchdb

// DatabaseRef holds the reference to a database.
type DatabaseRef struct {
	client *Client
	name   string
}

// Create returns the database create request.
func (dr *DatabaseRef) Create() *DatabaseCreateRequest {
	return &DatabaseCreateRequest{database: dr}
}

// Delete returns the database delete request.
func (dr *DatabaseRef) Delete() *DatabaseDeleteRequest {
	return &DatabaseDeleteRequest{database: dr}
}

// Document returns a reference to a document in the database.
func (dr *DatabaseRef) Document() *DocumentRef {
	return &DocumentRef{database: dr}
}
