package service

type ThumbnailRequest struct {
	VideoURLorID string `json:"video_url_or_id"`
	Quality      string `json:"quality,omitempty"`
}

type ThumbnailResponse struct {
	VideoID    string            `json:"video_id"`
	VideoURL   string            `json:"video_url"`
	Thumbnails map[string]string `json:"thumbnails"`
}
