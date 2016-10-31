package tworker

import (
	"net/http"
	"net/url"
)

func GetByProxy(url_addr, proxy_addr string) (*http.Response, error) {
	var err error
	request, _ := http.NewRequest("GET", url_addr, nil)
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	return client.Do(request)
}
