package domain

import (
	"bufio"
	"net"
	"strings"
)

type Cidr string

func MkCidrsFetcher(reader UrlReader) CidrsFetcher {
	return &cidrsFetcherImpl{
		reader: reader,
	}
}

type CidrsFetcher interface {
	Fetch(srcUrls *[]string) ([]Cidr, error)
}

type cidrsFetcherImpl struct {
	reader UrlReader
}

func (l *cidrsFetcherImpl) Fetch(srcUrls *[]string) ([]Cidr, error) {
	return fetchCidrsFromSrcUrls(srcUrls, l.reader)
}

func fetchCidrsFromSrcUrls(urls *[]string, reader UrlReader) ([]Cidr, error) {
	var allIps []Cidr
	for _, srcIPUrl := range *urls {
		if cidrStrs, err := fetchCidrStrs(reader, &srcIPUrl); err != nil {
			return nil, err
		} else {
			for _, cidrStr := range cidrStrs {
				trimmed := strings.TrimSpace(cidrStr)
				if isValidCidr(trimmed) {
					allIps = append(allIps, Cidr(trimmed))
				}
			}
		}
	}
	return allIps, nil
}

func fetchCidrStrs(getter UrlReader, address *string) ([]string, error) {
	if response, err := getter.GetUrl(*address); err != nil {
		return nil, err
	} else {
		scanner := bufio.NewScanner(response)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, nil
	}
}

func isValidCidr(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}
