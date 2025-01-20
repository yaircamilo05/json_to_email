package models

type SchemaStats struct {
	DocTimeMin     int64   `json:"doc_time_min"`
	DocTimeMax     int64   `json:"doc_time_max"`
	CreatedAt      int64   `json:"created_at"`
	DocNum         int64   `json:"doc_num"`
	FileNum        int64   `json:"file_num"`
	StorageSize    float64 `json:"storage_size"`
	CompressedSize float64 `json:"compressed_size"`
	IndexSize      float64 `json:"index_size"`
}

type Schema struct {
	Name        string      `json:"name"`
	StorageType string      `json:"storage_type"`
	StreamType  string      `json:"stream_type"`
	Stats       SchemaStats `json:"stats"`
}

type GetSchemasResponse struct {
	List []Schema `json:"list"`
}
