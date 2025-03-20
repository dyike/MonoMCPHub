package unsplash

import (
	"net/url"
	"os"
	"testing"
	"time"
)

func TestSearchPhotos(t *testing.T) {
	cfg := &UnsplashConfig{
		AccessKey: os.Getenv("UNSPLASH_API_KEY"),
		Timeout:   10 * time.Second,
	}
	client := NewUnsplashClient(cfg)
	query := url.Values{}
	query.Add("query", "nature")
	query.Add("page", "1")
	query.Add("per_page", "10")
	query.Add("order_by", "relevant")
	photos, err := client.SearchPhotos(query)
	if err != nil {
		t.Fatalf("Failed to search photos: %v", err)
	}
	t.Logf("Found photos: %v", photos)
}
