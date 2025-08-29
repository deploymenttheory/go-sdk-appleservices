package itunes_search

import "net/url"

func (c *Client) buildQueryString(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return values.Encode()
}
