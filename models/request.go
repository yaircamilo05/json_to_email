package models

type ProcessRequest struct {
	DirPath    string `json:"dir_path"`
	StreamName string `json:"stream_name"`
}

type Query struct {
	SQL       string `json:"sql"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	From      int    `json:"from"`
	Size      int    `json:"size"`
}

type GetAllEmailsRequest struct {
	Schema string `json:"schema"`
	SQL    string `json:"sql"`
	From   int    `json:"from"`
	Size   int    `json:"size"`
}

type searchEmailsResponse struct {
	SQL        Query  `json:"sql"`
	Searchtype string `json:"searchtype"`
	Timeout    string `json:"timeout"`
}
