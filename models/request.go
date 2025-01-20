package models

type ProcessRequest struct {
	DirPath    string `json:"dir_path"`
	StreamName string `json:"stream_name"`
}

type GetAllEmailsRequest struct {
	StreamName string `json:"stream_name"`
}
