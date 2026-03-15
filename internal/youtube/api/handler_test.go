package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"apiservices/youtube-thumbnail-api/internal/youtube/service"
)

func TestThumbnailEndpoint(t *testing.T) {
	t.Parallel()

	h := NewHandler(service.NewService())
	req := httptest.NewRequest(http.MethodPost, "/v1/youtube/thumbnail", strings.NewReader(`{"video_url_or_id":"https://www.youtube.com/watch?v=ifHYS1Ji9Ag","quality":"maxres"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}

	var payload struct {
		Data service.ThumbnailResponse `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload.Data.VideoID != "ifHYS1Ji9Ag" {
		t.Fatalf("unexpected video id: %s", payload.Data.VideoID)
	}
	if payload.Data.Thumbnails["default"] != "https://img.youtube.com/vi/ifHYS1Ji9Ag/default.jpg" {
		t.Fatalf("unexpected default thumbnail: %#v", payload.Data.Thumbnails)
	}
	if payload.Data.Thumbnails["maxres"] != "https://img.youtube.com/vi/ifHYS1Ji9Ag/maxresdefault.jpg" {
		t.Fatalf("unexpected maxres thumbnail: %#v", payload.Data.Thumbnails)
	}
}

func TestRemovedRoutesReturnNotFound(t *testing.T) {
	t.Parallel()

	h := NewHandler(service.NewService())

	for _, path := range []string{"/v1/youtube/metadata", "/v1/youtube/process"} {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404 for %s, got %d body=%s", path, rr.Code, rr.Body.String())
		}
	}
}
