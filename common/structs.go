package common

type Redirect struct {
	Identifier   string `json:"identifier"`
	Location     string `json:"location"`
	DeletionKey  string `json:"deletion_key"`
	CreationDate string `json:"creation_date"`
	Clicks       int    `json:"clicks"`
}

type RedirectCreated struct {
	ShortenedURL string `json:"shortened_url"`
	TargetURL    string `json:"target_url"`
	DeletionURL  string `json:"deletion_url"`
}

type ErrorResponse struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}
