package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestGeoRules(t *testing.T) {
	checkRules(t, "Latitude", is.Latitude,
		[]string{"0", "45.5", "-90", "90"},
		[]string{"", "90.1", "-91", "abc", "NaN", "nan"})

	checkRules(t, "Longitude", is.Longitude,
		[]string{"0", "-122.4", "180", "-180"},
		[]string{"", "180.1", "-181", "abc", "NaN", "nan"})
}

func TestSSN(t *testing.T) {
	checkRules(t, "SSN", is.SSN,
		[]string{"123-45-6789", "782-01-2345"},
		[]string{"", "000-45-6789", "666-45-6789", "900-45-6789",
			"123-00-6789", "123-45-0000", "12a-45-6789"})
}

func TestSemver(t *testing.T) {
	checkRules(t, "Semver", is.Semver,
		[]string{"1.0.0", "0.2.3", "1.0.0-alpha", "1.0.0-alpha.1+build.5", "10.20.30"},
		[]string{"", "1.0", "v1.0.0", "01.0.0", "1.0.0-"})
}

func TestOrigin(t *testing.T) {
	checkRules(t, "Origin", is.Origin,
		[]string{"https://example.com", "http://localhost:3000", "https://api.example.co.uk:8443"},
		[]string{"", "example.com", "ftp://example.com", "https://example.com/path"})
}

func TestDataURI(t *testing.T) {
	checkRules(t, "DataURI", is.DataURI,
		[]string{
			"data:image/png;base64,iVBORw0KGgo=",
			"data:text/plain;charset=utf-8,hello",
			"data:,plain",
		},
		[]string{"", "http://example.com/a.png", "data:image/png;base64,!!not-base64"})
}
