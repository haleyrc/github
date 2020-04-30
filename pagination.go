package github

import "strings"

// Link: <https://api.github.com/repositories/172137510/issues?page=2>; rel="next", <https://api.github.com/repositories/172137510/issues?page=40>; rel="last"

func (c *Client) parseLinkHeader(hdr string) map[string]string {
	links := make(map[string]string)
	resources := strings.Split(hdr, ", ")
	for _, resource := range resources {
		parts := strings.Split(resource, "; ")
		url, rel := parts[0], parts[1]
		links[c.parseRel(rel)] = c.parseURL(url)
	}
	return links
}

// TODO (RCH): This is unlikely to work all the time, but it works for the
// sample so /shrug
func (c *Client) parseRel(rel string) string {
	rel = strings.TrimPrefix(rel, "rel=\"")
	rel = strings.TrimSuffix(rel, "\"")
	return rel
}

// TODO (RCH): This is unlikely to work all the time, but it works for the
// sample so /shrug
func (c *Client) parseURL(url string) string {
	url = strings.TrimPrefix(url, "<")
	url = strings.TrimPrefix(url, c.URL)
	url = strings.TrimSuffix(url, ">")
	return url
}
