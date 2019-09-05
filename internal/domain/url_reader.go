package domain

import (
	"io"
)

type UrlReader interface {
	GetUrl(url string) (io.Reader, error)
}
