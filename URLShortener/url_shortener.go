package URLShortener

import (
	"errors"
	"net/url"
	"sync"
)

type URLShortener struct {
	mu         sync.Mutex
	listOfUrls map[string]string
	counter    int64
}

var (
	ErrInvalidURL = errors.New("invalid URL")
)

const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewURLShortener() *URLShortener {
	return &URLShortener{
		listOfUrls: make(map[string]string),
		counter:    0,
	}
}

func isURLValid(rawUrl string) bool {
	if len(rawUrl) == 0 {
		return false
	}

	parsedUrl, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return false
	}

	return (parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https") && parsedUrl.Hostname() != ""
}

func (sh *URLShortener) GetShortenedURL(rawUrl string) (string, error) {
	//validate the url
	if !isURLValid(rawUrl) {
		return "", ErrInvalidURL
	}

	sh.mu.Lock()
	defer sh.mu.Unlock()

	//check if the url already exist
	if shortenedUrl, isExist := sh.listOfUrls[rawUrl]; isExist {
		return shortenedUrl, nil
	}

	//shorten the url
	//save it in map
	//increment the counter
	shortenedURL := shortenURL(sh.counter)
	sh.listOfUrls[rawUrl] = shortenedURL
	sh.counter++

	return shortenedURL, nil
}

func shortenURL(counter int64) string {
	if counter == 0 {
		return string(alphabets[0])
	}

	var base62EncodedReversed []byte

	for counter > 0 {
		remainder := counter % 62
		base62EncodedReversed = append(base62EncodedReversed, alphabets[remainder])
		counter = counter / 62
	}

	var base62Encoded []byte
	for i := len(base62EncodedReversed) - 1; i >= 0; i-- {
		base62Encoded = append(base62Encoded, base62EncodedReversed[i])
	}

	return string(base62Encoded)
}
