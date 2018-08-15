package couchdb

// Document defines a document.
type Document map[string]interface{}

// ID returns the document's id if defined.
func (d Document) ID() string {
	if v, ok := d["_id"].(string); ok {
		return v
	}
	return ""
}

// Revision returns the document's revision if defined.
func (d Document) Revision() string {
	if v, ok := d["_rev"].(string); ok {
		return v
	}
	return ""
}
