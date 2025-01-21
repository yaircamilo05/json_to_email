package models

type TookDetail struct {
	Total            int `json:"total"`
	IdxTook          int `json:"idx_took"`
	WaitQueue        int `json:"wait_queue"`
	ClusterTotal     int `json:"cluster_total"`
	ClusterWaitQueue int `json:"cluster_wait_queue"`
}

type SearchResponse struct {
	Hits        []Email `json:"hits"`
	ScanRecords int     `json:"scan_records"`
}
