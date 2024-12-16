package helpers

import (
	"os"
	"strings"
)

// ensures that the given url starts with "http://"
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// checks if the given url contains current domain 
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url,"https://","",1)
	newURL = strings.Replace(newURL,"http://","",1)
	newURL = strings.Replace(newURL,"www.","",1)
	newURL = strings.Split(newURL,"/")[0]
	
	return newURL != os.Getenv("DOMAIN")
}