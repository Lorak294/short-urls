package validation

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestValidteUrl(t *testing.T) {
	testcases := []struct {
		inputUrl string
		expectedResult bool
	}{
		{"", false},
		{"https://go.dev", true},
		{"http://go.dev", true},
		{"://go.dev", false},
		{"go.dev", false},
		{"some_text", false},
		{"www.go.dev", false},
		{"349088123908123", false},
		{"ąćźż\n\t", false},
	}

	for _, tc := range testcases {
		res := ValidteUrl(tc.inputUrl)
		if res != tc.expectedResult {
			t.Errorf("ValidteUrl() inputUrl=%v, result=%v, expectedResult=%v",tc.inputUrl, res,tc.expectedResult)
		}
	}
}

func TestValidateDomainError(t *testing.T) {

	domain := "test_domain"
	os.Setenv("DOMAIN", domain)

	testcases := []struct {
		inputUrl string
		expectedResult bool
	}{
		{domain, false},
		{"https://" + domain, false},
		{"http://" + domain, false},
		{"www." + domain, false},
		{"https://" + domain + "/someaddress", false},
		{"http://" + domain + "/someaddress", false},
		{"https://go.dev", true},
		{"http://go.dev", true},
		{"www.go.dev", true},
	}

	for _, tc := range testcases {
		res := ValidateDomainError(tc.inputUrl)
		if res != tc.expectedResult {
			t.Errorf("ValidateDomainError() inputUrl=%v, result=%v, expectedResult=%v",tc.inputUrl, res,tc.expectedResult)
		}
	}

	_ = os.Unsetenv("DOMAIN")
}

func TestValidteShort(t *testing.T) {
	testcases := []struct {
		inputShort string
		expectedResult bool
	}{
		{"", true},
		{strings.Repeat("c", MIN_SHORT_LEN), true},
		{strings.Repeat("c", MIN_SHORT_LEN-1), false},
		{strings.Repeat("c", MIN_SHORT_LEN+1), true},
		{strings.Repeat("c", MAX_SHORT_LEN), true},
		{strings.Repeat("c", MAX_SHORT_LEN-1), true},
		{strings.Repeat("c", MAX_SHORT_LEN+1), false},
		{strings.Repeat("$", MIN_SHORT_LEN), false},
		{strings.Repeat("$", MIN_SHORT_LEN-1), false},
		{strings.Repeat("$", MIN_SHORT_LEN+1), false},
		{strings.Repeat("$", MAX_SHORT_LEN), false},
		{strings.Repeat("$", MAX_SHORT_LEN-1), false},
		{strings.Repeat("$", MIN_SHORT_LEN+1), false},
	}

	for _, tc := range testcases {
		res := ValidateShort(tc.inputShort)
		if res != tc.expectedResult {
			t.Errorf("ValidateShort() inputShort=%v, result=%v, expectedResult=%v",tc.inputShort, res,tc.expectedResult)
		}
	}
}

func TestValidteExpiry(t *testing.T) {
	testcases := []struct {
		inputExpiry time.Duration
		expectedResult bool
	}{
		{0, true},
		{MAX_EXPIRY, true},
		{MAX_EXPIRY+1, false},
		{MAX_EXPIRY-1, true},
		{-1, false},
	}

	for _, tc := range testcases {
		res := ValidateExpiry(tc.inputExpiry)
		if res != tc.expectedResult {
			t.Errorf("ValidateExpiry() inputExpiry=%v, result=%v, expectedResult=%v",tc.inputExpiry, res,tc.expectedResult)
		}
	}
}