package infra

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpClientUrlReader_GetUrl_ok(t *testing.T) {
	called := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called += 1
		_, _ = fmt.Fprintln(w, "123")
	}))
	defer ts.Close()
	mockServerURL := ts.URL

	urlReader := MkHttpClientUrlReader(http.DefaultClient)
	resp, err := urlReader.GetUrl(mockServerURL)
	data, _ := ioutil.ReadAll(resp)
	bodyAsString := string(data)
	assert.Equal(t, 1, called)
	assert.Nil(t, err)
	assert.Equal(t, "123\n", bodyAsString)
}

func TestHttpClientUrlReader_GetUrl_err(t *testing.T) {
	called := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called += 1
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	mockServerURL := ts.URL

	urlReader := MkHttpClientUrlReader(http.DefaultClient)
	_, err := urlReader.GetUrl(mockServerURL)
	assert.Equal(t, 1, called)
	assert.Contains(t, err.Error(), "404")
}
