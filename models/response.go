package models

type SearchResponse struct {
	Hits        []Email `json:"hits"`
	ScanRecords int     `json:"scan_records"`
}
