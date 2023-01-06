package health

type Health struct {
	Status   string   `json:"status"`
	Database DBHealth `json:"database"`
}

type DBHealth struct {
	Status string `json:"status"`
}
