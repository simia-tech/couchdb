package couchdb

// DocumentRef holds the reference to a document.
type DocumentRef struct {
	database *DatabaseRef
}

// Create returns the document create request.
func (dr *DocumentRef) Create(document Document) *DocumentCreateRequest {
	return &DocumentCreateRequest{document: dr, doc: document}
}

// Fetch returns the document fetch request.
func (dr *DocumentRef) Fetch(id string) *DocumentFetchRequest {
	return &DocumentFetchRequest{document: dr, id: id}
}

// FetchMeta returns the document fetch meta request.
func (dr *DocumentRef) FetchMeta(id string) *DocumentFetchMetaRequest {
	return &DocumentFetchMetaRequest{document: dr, id: id}
}

// Update returns the document update request.
func (dr *DocumentRef) Update(id string, document Document) *DocumentUpdateRequest {
	return &DocumentUpdateRequest{document: dr, id: id, doc: document}
}

// Delete returns the document delete request.
func (dr *DocumentRef) Delete(id, revision string) *DocumentDeleteRequest {
	return &DocumentDeleteRequest{document: dr, id: id, revision: revision}
}
