package models

type TookDetail struct {
	Total            int `json:"total"`
	IdxTook          int `json:"idx_took"`
	WaitQueue        int `json:"wait_queue"`
	ClusterTotal     int `json:"cluster_total"`
	ClusterWaitQueue int `json:"cluster_wait_queue"`
}

type SearchResponse struct {
	Took             int        `json:"took"`
	TookDetail       TookDetail `json:"took_detail"`
	Hits             []Email    `json:"hits"`
	Total            int        `json:"total"`
	From             int        `json:"from"`
	Size             int        `json:"size"`
	CachedRatio      int        `json:"cached_ratio"`
	ScanSize         int        `json:"scan_size"`
	IdxScanSize      int        `json:"idx_scan_size"`
	ScanRecords      int        `json:"scan_records"`
	TraceID          string     `json:"trace_id"`
	IsPartial        bool       `json:"is_partial"`
	ResultCacheRatio int        `json:"result_cache_ratio"`
	OrderBy          string     `json:"order_by"`
}
