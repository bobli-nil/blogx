package xss

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func XssFilter(content string) (newContent string, err error) {
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(content)))
	if err != nil {
		return
	}
	contentDoc.Find("script").Remove()
	contentDoc.Find("image").Remove()
	contentDoc.Find("iframe").Remove()
	contentDoc.Find("video").Remove()
	contentDoc.Find("audio").Remove()
	return contentDoc.Text(), nil
}
