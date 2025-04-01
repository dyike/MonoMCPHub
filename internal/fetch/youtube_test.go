package fetch

import (
	"net/http"
	"testing"

	"github.com/kkdai/youtube/v2"
)

func TestFetchYouTubeTranscript(t *testing.T) {
	youtubeClient := &youtube.Client{
		HTTPClient: &http.Client{},
	}
	url := "https://www.youtube.com/watch?v=vNDjoNuT9kM"
	videoID, _ := extractVideoID(url)
	t.Logf("videoID: %v", videoID)
	video, err := youtubeClient.GetVideo(videoID)
	if err != nil {
		t.Fatalf("failed to get video: %v", err)
	}
	t.Logf("video: %v", video)
}
