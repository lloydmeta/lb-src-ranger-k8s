package domain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestCidrsFetcherImpl_Fetch_ok_empty(t *testing.T) {
	mockReader := mockUrlReader{s: "1.1.1.0/25"}
	fetcher := MkCidrsFetcher(&mockReader)
	r, err := fetcher.Fetch(&[]string{})
	assert.Equal(t, 0, mockReader.getCalled)
	assert.Nil(t, err)
	assert.Empty(t, r)
}
func TestCidrsFetcherImpl_Fetch_ok(t *testing.T) {
	mockReader := mockUrlReader{s: "1.1.1.0/25"}
	fetcher := MkCidrsFetcher(&mockReader)
	r, err := fetcher.Fetch(&[]string{"hello"})
	assert.Equal(t, 1, mockReader.getCalled)
	assert.Nil(t, err)
	assert.Equal(t, []Cidr{"1.1.1.0/25"}, r)
}

func TestCidrsFetcherImpl_Fetch_err(t *testing.T) {
	mockReader := mockUrlReader{err: true}
	fetcher := MkCidrsFetcher(&mockReader)
	r, err := fetcher.Fetch(&[]string{"hello"})
	assert.Equal(t, 1, mockReader.getCalled)
	assert.NotNil(t, err)
	assert.Empty(t, r)
}

func Test_isValidCidr(t *testing.T) {
	assert.False(t, isValidCidr("lol"))
	assert.True(t, isValidCidr("1.1.1.0/25"))
}

type mockUrlReader struct {
	err       bool
	s         string
	getCalled int
}

func (m *mockUrlReader) GetUrl(url string) (io.Reader, error) {
	m.getCalled += 1
	if m.err {
		return nil, errors.New("ugh")
	} else {
		return strings.NewReader(m.s), nil
	}
}
