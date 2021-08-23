package common

type Redirect struct {
	Identifier   string `json:"identifier"`
	Location     string `json:"location"`
	CreationDate string `json:"creation_date"`
	Clicks       int    `json:"clicks"`
}

type ErrorResponse struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}
