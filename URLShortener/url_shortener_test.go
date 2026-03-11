package URLShortener

import (
	"testing"
)

func TestGetShortenedURL(t *testing.T) {
	rawUrl := "http://www.google.com"

	urlshortener := NewURLShortener()
	shortenedUrl, err := urlshortener.GetShortenedURL(rawUrl)
	if err != nil {
		t.Error(err)
	}

	if shortenedUrl != "a" {
		t.Error("Shortened URL does not match")
	}

	rawUrl = "https://www.yt.com"

	shortenedUrl, err = urlshortener.GetShortenedURL(rawUrl)
	if err != nil {
		t.Error(err)
	}

	if shortenedUrl != "b" {
		t.Error("Shortened URL does not match")
	}
}
