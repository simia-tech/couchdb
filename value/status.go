package value

// Status holds all infos of a status response.
type Status struct {
	OK     bool   `json:"ok"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}
