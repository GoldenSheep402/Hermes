package utils

import (
	"io"
	"net/http"
	"net/url"
)

func UrlAppend(rawUrl string, queries map[string][]string) string {
	uu, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	q1 := uu.Query()
	for k, v := range queries {
		q1[k] = v
	}
	uu.RawQuery = q1.Encode()
	return uu.String()
}

func GetMPFDContentType(src io.Reader) (string, error) {
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}
