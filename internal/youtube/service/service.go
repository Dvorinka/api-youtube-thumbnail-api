package service

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var videoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateThumbnails(_ context.Context, req ThumbnailRequest) (ThumbnailResponse, error) {
	videoID, err := s.extractVideoID(req.VideoURLorID)
	if err != nil {
		return ThumbnailResponse{}, err
	}

	requestedQuality, selectedPath, err := normalizeQuality(req.Quality)
	if err != nil {
		return ThumbnailResponse{}, err
	}

	thumbnails := map[string]string{
		"default": buildThumbnailURL(videoID, "default.jpg"),
	}
	if requestedQuality != "default" {
		thumbnails[requestedQuality] = buildThumbnailURL(videoID, selectedPath)
	}

	return ThumbnailResponse{
		VideoID:    videoID,
		VideoURL:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Thumbnails: thumbnails,
	}, nil
}

func normalizeQuality(raw string) (string, string, error) {
	quality := strings.ToLower(strings.TrimSpace(raw))
	if quality == "" {
		quality = "maxres"
	}

	switch quality {
	case "default", "low":
		return quality, "default.jpg", nil
	case "medium":
		return quality, "mqdefault.jpg", nil
	case "high":
		return quality, "hqdefault.jpg", nil
	case "maxres":
		return quality, "maxresdefault.jpg", nil
	default:
		return "", "", fmt.Errorf("quality must be one of: maxres, high, medium, low, default")
	}
}

func (s *Service) extractVideoID(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty video URL or ID")
	}

	if videoIDPattern.MatchString(input) {
		return input, nil
	}

	parsed, err := url.Parse(input)
	if err != nil || parsed.Host == "" {
		return "", fmt.Errorf("could not extract valid video ID from URL or input")
	}

	host := strings.ToLower(parsed.Hostname())
	var videoID string

	switch {
	case strings.Contains(host, "youtu.be"):
		videoID = strings.Trim(parsed.Path, "/")
	case strings.Contains(host, "youtube.com"), strings.Contains(host, "youtube-nocookie.com"):
		switch {
		case parsed.Path == "/watch":
			videoID = parsed.Query().Get("v")
		case strings.HasPrefix(parsed.Path, "/embed/"):
			videoID = strings.TrimPrefix(parsed.Path, "/embed/")
		case strings.HasPrefix(parsed.Path, "/shorts/"):
			videoID = strings.TrimPrefix(parsed.Path, "/shorts/")
		}
	}

	if idx := strings.IndexByte(videoID, '/'); idx >= 0 {
		videoID = videoID[:idx]
	}
	videoID = strings.TrimSpace(videoID)
	if !videoIDPattern.MatchString(videoID) {
		return "", fmt.Errorf("could not extract valid video ID from URL or input")
	}

	return videoID, nil
}

func buildThumbnailURL(videoID, path string) string {
	return fmt.Sprintf("https://img.youtube.com/vi/%s/%s", videoID, path)
}
