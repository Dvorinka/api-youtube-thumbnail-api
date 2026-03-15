package service

import (
	"context"
	"testing"
)

func TestGenerateThumbnailsLowQualityAlias(t *testing.T) {
	t.Parallel()

	svc := NewService()
	result, err := svc.GenerateThumbnails(context.Background(), ThumbnailRequest{
		VideoURLorID: "https://youtu.be/ifHYS1Ji9Ag",
		Quality:      "low",
	})
	if err != nil {
		t.Fatalf("generate thumbnails: %v", err)
	}

	if result.VideoID != "ifHYS1Ji9Ag" {
		t.Fatalf("unexpected video id: %s", result.VideoID)
	}
	if result.Thumbnails["low"] != "https://img.youtube.com/vi/ifHYS1Ji9Ag/default.jpg" {
		t.Fatalf("unexpected thumbnails map: %#v", result.Thumbnails)
	}
}

func TestGenerateThumbnailsInvalidQuality(t *testing.T) {
	t.Parallel()

	svc := NewService()
	_, err := svc.GenerateThumbnails(context.Background(), ThumbnailRequest{
		VideoURLorID: "ifHYS1Ji9Ag",
		Quality:      "ultra",
	})
	if err == nil {
		t.Fatal("expected invalid quality error")
	}
}
