package validation

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"short-urls/contracts"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

func ValidateShortenRequest(request contracts.ShortenRequest) error {

	var validation_errors error = nil
	if !ValidteUrl(request.Url) {
		validation_errors = errors.Join(validation_errors, errors.New("provided Url is invalid"))
	}
	if !ValidateDomainError(request.Url) {
		validation_errors = errors.Join(validation_errors, errors.New("provided Url cannot be shortened"))
	}
	if !ValidateShort(request.CustomShort) {
		validation_errors = errors.Join(validation_errors,fmt.Errorf("short must be an alphanumeric string of (%d,%d) length", MIN_SHORT_LEN, MAX_SHORT_LEN))
	}
	if !ValidateExpiry(request.Expiry) {
		validation_errors = errors.Join(validation_errors, fmt.Errorf("maximum allowed expiry is (%d)",MAX_EXPIRY))
	}
	return validation_errors
}

func ValidteUrl(url string) bool {
	return govalidator.IsURL(url)
}

// checks if the given url contains current domain
func ValidateDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "https://", "", 1)
	newURL = strings.Replace(newURL, "http://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	return newURL != os.Getenv("DOMAIN")
}

func ValidateShort(short string) bool {
	return len(short) == 0 || (regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(short) && len(short) >= MIN_SHORT_LEN && len(short) <= MAX_SHORT_LEN)
}

func ValidateExpiry(expiry time.Duration) bool {
	return expiry >= 0 && expiry < MAX_EXPIRY
}