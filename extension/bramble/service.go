package bramble

type Service struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Schema  string `json:"schema"`
}
