package utils

import (
	"context"
	"net/url"
	"time"

	"github.com/otiai10/opengraph/v2"
)

type OGPInfo struct {
	URL         string
	Title       string
	Description string
	Image       string
	Domain      string
}

// FetchOGPInfo は指定されたURLからOGP情報を取得する
func FetchOGPInfo(targetUrl string, timeout time.Duration) (*OGPInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	og, err := opengraph.Fetch(targetUrl, opengraph.Intent{Context: ctx})
	if err != nil {
		// OGP取得に失敗しても、URLなどの基本情報は返せるようにする
		parsedUrl, _ := url.Parse(targetUrl)
		return &OGPInfo{
			URL:    targetUrl,
			Domain: parsedUrl.Host,
		}, err
	}

	parsedUrl, _ := url.Parse(targetUrl)
	info := &OGPInfo{
		URL:         targetUrl,
		Title:       og.Title,
		Description: og.Description,
		Domain:      parsedUrl.Host,
	}

	if len(og.Image) > 0 {
		info.Image = og.Image[0].URL
	}

	return info, nil
}
