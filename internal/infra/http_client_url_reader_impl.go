package infra

import (
	"errors"
	"fmt"
	"github.com/lloydmeta/lb-src-ranger-k8s/internal/domain"
	"io"
	"net/http"
)

type HttpClientUrlReader struct {
	client *http.Client
}

func MkHttpClientUrlReader(client *http.Client) domain.UrlReader {
	return &HttpClientUrlReader{client: client}
}

func (h *HttpClientUrlReader) GetUrl(url string) (io.Reader, error) {
	if resp, err := h.client.Get(url); err != nil {
		return nil, err
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return resp.Body, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid response code [%v]", resp.StatusCode))
		}
	}
}
