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

type Text struct {
	Identifier   string `json:"identifier"`
	Content      string `json:"content"`
	Preview      string `json:"preview"`
	DeletionKey  string `json:"deletion_key"`
	CreationDate string `json:"creation_date"`
	Views        int    `json:"views"`
}

type TextCreated struct {
	LocationURL string `json:"location_url"`
	Content     string `json:"content"`
	DeletionURL string `json:"deletion_url"`
}

type Image struct {
	Identifier   string `json:"identifier"`
	FileLocation string `json:"file_location"`
	DeletionKey  string `json:"deletion_key"`
	CreationDate string `json:"creation_date"`
	FileSize     int64  `json:"file_size"`
	Views        int    `json:"views"`
}

type ImageCreated struct {
	LocationURL string `json:"location_url"`
	FileSize    int64  `json:"file_size"`
	DeletionURL string `json:"deletion_url"`
}

type ErrorResponse struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}
