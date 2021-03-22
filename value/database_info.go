package value

// DatabaseInfo holds infos about the database.
type DatabaseInfo struct {
	Name                  string  `json:"db_name"`
	PurgeSequence         string  `json:"purge_seq"`
	UpdateSequence        string  `json:"update_seq"`
	DocumentDeletionCount uint    `json:"doc_del_count"`
	DocumentCount         uint    `json:"doc_count"`
	DiskFormatVersion     uint    `json:"disk_format_version"`
	CompactRunning        bool    `json:"compact_running"`
	Cluster               Cluster `json:"cluster"`
}
