package helpers

// ensures that the given url starts with "http://"
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

