package value

// InstanceInfo holds information about the couchdb instance.
type InstanceInfo struct {
	CouchDB  string   `json:"couchdb"`
	Version  string   `json:"version"`
	GitSHA   string   `json:"git_sha"`
	UUID     string   `json:"uuid"`
	Features []string `json:"features"`
	Vendor   Vendor   `json:"vendor"`
}
