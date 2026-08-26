// Package is provides high-frequency format assertion rules, all of type verax.Rule[string].
//
// All rules are strict mode: an empty string is always considered invalid;
// for optional fields wrap with rules.Optional, e.g. rules.Optional(is.Email).
//
// Usage notes:
//   - CountryCode/CurrencyCode/LanguageCode include the complete ISO standard data tables;
//   - Latitude/Longitude only accept string form; for numeric ranges use rules.Between directly.
package is

import (
	"encoding/base64"
	"regexp"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/internal/codes"
)

var (
	// ErrEmail is the error returned when Email validation fails.
	ErrEmail = verax.NewError(codes.CodeEmail, "must be a valid email address")
	// ErrAlpha is the error returned when Alpha validation fails.
	ErrAlpha = verax.NewError(codes.CodeAlpha, "must contain English letters only")
	// ErrAlphanumeric is the error returned when Alphanumeric validation fails.
	ErrAlphanumeric = verax.NewError(codes.CodeAlphanumeric, "must contain English letters and digits only")
	// ErrDigit is the error returned when Digit validation fails.
	ErrDigit = verax.NewError(codes.CodeDigit, "must contain digits only")
	// ErrUTFLetter is the error returned when UTFLetter validation fails.
	ErrUTFLetter = verax.NewError(codes.CodeUTFLetter, "must contain unicode letter characters only")
	// ErrUTFDigit is the error returned when UTFDigit validation fails.
	ErrUTFDigit = verax.NewError(codes.CodeUTFDigit, "must contain unicode digit characters only")
	// ErrUTFNumeric is the error returned when UTFNumeric validation fails.
	ErrUTFNumeric = verax.NewError(codes.CodeUTFNumeric, "must contain unicode number characters only")
	// ErrUTFLetterNumeric is the error returned when UTFLetterNumeric validation fails.
	ErrUTFLetterNumeric = verax.NewError(codes.CodeUTFLetterNumeric, "must contain unicode letters and numbers only")
	// ErrLowerCase is the error returned when LowerCase validation fails.
	ErrLowerCase = verax.NewError(codes.CodeLowerCase, "must be in lower case")
	// ErrUpperCase is the error returned when UpperCase validation fails.
	ErrUpperCase = verax.NewError(codes.CodeUpperCase, "must be in upper case")
	// ErrASCII is the error returned when ASCII validation fails.
	ErrASCII = verax.NewError(codes.CodeASCII, "must contain ASCII characters only")
	// ErrPrintableASCII is the error returned when PrintableASCII validation fails.
	ErrPrintableASCII = verax.NewError(codes.CodePrintableASCII, "must contain printable ASCII characters only")
	// ErrHexadecimal is the error returned when Hexadecimal validation fails.
	ErrHexadecimal = verax.NewError(codes.CodeHexadecimal, "must be a valid hexadecimal number")
	// ErrHexColor is the error returned when HexColor validation fails.
	ErrHexColor = verax.NewError(codes.CodeHexColor, "must be a valid hexadecimal color code")
	// ErrRGBColor is the error returned when RGBColor validation fails.
	ErrRGBColor = verax.NewError(codes.CodeRGBColor, "must be a valid RGB color code")
	// ErrInt is the error returned when Int validation fails.
	ErrInt = verax.NewError(codes.CodeInt, "must be an integer")
	// ErrFloat is the error returned when Float validation fails.
	ErrFloat = verax.NewError(codes.CodeFloat, "must be a floating point number")
	// ErrUUID is the error returned when UUID validation fails.
	ErrUUID = verax.NewError(codes.CodeUUID, "must be a valid UUID")
	// ErrUUIDv3 is the error returned when UUIDv3 validation fails.
	ErrUUIDv3 = verax.NewError(codes.CodeUUIDv3, "must be a valid UUID v3")
	// ErrUUIDv4 is the error returned when UUIDv4 validation fails.
	ErrUUIDv4 = verax.NewError(codes.CodeUUIDv4, "must be a valid UUID v4")
	// ErrUUIDv5 is the error returned when UUIDv5 validation fails.
	ErrUUIDv5 = verax.NewError(codes.CodeUUIDv5, "must be a valid UUID v5")
	// ErrUUIDv7 is the error returned when UUIDv7 validation fails.
	ErrUUIDv7 = verax.NewError(codes.CodeUUIDv7, "must be a valid UUID v7")
	// ErrULID is the error returned when ULID validation fails.
	ErrULID = verax.NewError(codes.CodeULID, "must be a valid ULID")
	// ErrBase64 is the error returned when Base64 validation fails.
	ErrBase64 = verax.NewError(codes.CodeBase64, "must be encoded in base64")
	// ErrJSON is the error returned when JSON validation fails.
	ErrJSON = verax.NewError(codes.CodeJSON, "must be in valid JSON format")
	// ErrURL is the error returned when URL validation fails.
	ErrURL = verax.NewError(codes.CodeURL, "must be a valid URL")
	// ErrIP is the error returned when IP validation fails.
	ErrIP = verax.NewError(codes.CodeIP, "must be a valid IP address")
	// ErrIPv4 is the error returned when IPv4 validation fails.
	ErrIPv4 = verax.NewError(codes.CodeIPv4, "must be a valid IPv4 address")
	// ErrIPv6 is the error returned when IPv6 validation fails.
	ErrIPv6 = verax.NewError(codes.CodeIPv6, "must be a valid IPv6 address")
	// ErrMAC is the error returned when MAC validation fails.
	ErrMAC = verax.NewError(codes.CodeMAC, "must be a valid MAC address")
	// ErrCreditCard is the error returned when CreditCard validation fails.
	ErrCreditCard = verax.NewError(codes.CodeCreditCard, "must be a valid credit card number")
	// ErrISBN10 is the error returned when ISBN10 validation fails.
	ErrISBN10 = verax.NewError(codes.CodeISBN10, "must be a valid ISBN-10")
	// ErrISBN13 is the error returned when ISBN13 validation fails.
	ErrISBN13 = verax.NewError(codes.CodeISBN13, "must be a valid ISBN-13")
	// ErrISBN is the error returned when ISBN validation fails.
	ErrISBN = verax.NewError(codes.CodeISBN, "must be a valid ISBN")
	// ErrSubdomain is the error returned when Subdomain validation fails.
	ErrSubdomain = verax.NewError(codes.CodeSubdomain, "must be a valid subdomain")
	// ErrDomain is the error returned when Domain validation fails.
	ErrDomain = verax.NewError(codes.CodeDomain, "must be a valid domain")
	// ErrDNSName is the error returned when DNSName validation fails.
	ErrDNSName = verax.NewError(codes.CodeDNSName, "must be a valid DNS name")
	// ErrHost is the error returned when Host validation fails.
	ErrHost = verax.NewError(codes.CodeHost, "must be a valid IP address or DNS name")
	// ErrPort is the error returned when Port validation fails.
	ErrPort = verax.NewError(codes.CodePort, "must be a valid TCP port")
	// ErrE164 is the error returned when E164 validation fails.
	ErrE164 = verax.NewError(codes.CodeE164, "must be a valid E164 number")
	// ErrDialString is the error returned when DialString validation fails.
	ErrDialString = verax.NewError(codes.CodeDialString, "must be a valid dial string")
	// ErrMongoID is the error returned when MongoID validation fails.
	ErrMongoID = verax.NewError(codes.CodeMongoID, "must be a valid MongoDB ObjectID")
	// ErrDataURI is the error returned when DataURI validation fails.
	ErrDataURI = verax.NewError(codes.CodeDataURI, "must be a valid data URI")
	// ErrLatitude is the error returned when Latitude validation fails.
	ErrLatitude = verax.NewError(codes.CodeLatitude, "must be a valid latitude")
	// ErrLongitude is the error returned when Longitude validation fails.
	ErrLongitude = verax.NewError(codes.CodeLongitude, "must be a valid longitude")
	// ErrSSN is the error returned when SSN validation fails.
	ErrSSN = verax.NewError(codes.CodeSSN, "must be a valid social security number")
	// ErrSemver is the error returned when Semver validation fails.
	ErrSemver = verax.NewError(codes.CodeSemver, "must be a valid semantic version")
	// ErrOrigin is the error returned when Origin validation fails.
	ErrOrigin = verax.NewError(codes.CodeOrigin, "must be a valid origin")
	// ErrMultibyte is the error returned when Multibyte validation fails.
	ErrMultibyte = verax.NewError(codes.CodeMultibyte, "must contain multibyte characters")
	// ErrFullWidth is the error returned when FullWidth validation fails.
	ErrFullWidth = verax.NewError(codes.CodeFullWidth, "must contain full-width characters only")
	// ErrHalfWidth is the error returned when HalfWidth validation fails.
	ErrHalfWidth = verax.NewError(codes.CodeHalfWidth, "must contain half-width characters only")
	// ErrVariableWidth is the error returned when VariableWidth validation fails.
	ErrVariableWidth = verax.NewError(codes.CodeVariableWidth, "must contain both full-width and half-width characters")
	// ErrBase64URL is the error returned when Base64URL validation fails.
	ErrBase64URL = verax.NewError(codes.CodeBase64URL, "must be a valid base64url string")
	// ErrJWT is the error returned when JWT validation fails.
	ErrJWT = verax.NewError(codes.CodeJWT, "must be a valid JWT")
	// ErrBoolean is the error returned when Boolean validation fails.
	ErrBoolean = verax.NewError(codes.CodeBoolean, "must be a valid boolean value")
	// ErrTimeZone is the error returned when TimeZone validation fails.
	ErrTimeZone = verax.NewError(codes.CodeTimeZone, "must be a valid timezone")
	// ErrCountryCode is the error returned when CountryCode validation fails.
	ErrCountryCode = verax.NewError(codes.CodeCountryCode, "must be a valid ISO 3166-1 alpha-2 country code")
	// ErrCurrencyCode is the error returned when CurrencyCode validation fails.
	ErrCurrencyCode = verax.NewError(codes.CodeCurrencyCode, "must be a valid ISO 4217 currency code")
	// ErrLanguageCode is the error returned when LanguageCode validation fails.
	ErrLanguageCode = verax.NewError(codes.CodeLanguageCode, "must be a valid ISO 639-1 language code")
	// ErrCIDR is the error returned when CIDR validation fails.
	ErrCIDR = verax.NewError(codes.CodeCIDR, "must be a valid CIDR notation")
	// ErrCIDRv4 is the error returned when CIDRv4 validation fails.
	ErrCIDRv4 = verax.NewError(codes.CodeCIDRv4, "must be a valid IPv4 CIDR notation")
	// ErrCIDRv6 is the error returned when CIDRv6 validation fails.
	ErrCIDRv6 = verax.NewError(codes.CodeCIDRv6, "must be a valid IPv6 CIDR notation")
	// ErrRGBA is the error returned when RGBA validation fails.
	ErrRGBA = verax.NewError(codes.CodeRGBA, "must be a valid RGBA color")
	// ErrHSL is the error returned when HSL validation fails.
	ErrHSL = verax.NewError(codes.CodeHSL, "must be a valid HSL color")
	// ErrHSLA is the error returned when HSLA validation fails.
	ErrHSLA = verax.NewError(codes.CodeHSLA, "must be a valid HSLA color")
	// ErrMobileCN is the error returned when MobileCN validation fails.
	ErrMobileCN = verax.NewError(codes.CodeMobileCN, "must be a valid mainland China mobile number")
	// ErrIDCardCN is the error returned when IDCardCN validation fails.
	ErrIDCardCN = verax.NewError(codes.CodeIDCardCN, "must be a valid mainland China ID card number")
	// ErrPostalCodeCN is the error returned when PostalCodeCN validation fails.
	ErrPostalCodeCN = verax.NewError(codes.CodePostalCodeCN, "must be a valid mainland China postal code")
)

var (
	// emailPattern is the email format regex, using the same rule as input[type=email] in the WHATWG HTML spec
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	// uuidPattern is the standard 36-character hyphenated UUID form, segment lengths 8-4-4-4-12, hex case-insensitive
	uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// base64Std is the standard base64 encoder (RFC 4648, with padding), used for Base64 and DataURI decode checks
	base64Std = base64.StdEncoding

	// English letters only
	alphaPattern = regexp.MustCompile(`^[a-zA-Z]+$`)
	// English letters and digits
	alphanumericPattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	// ASCII digits only
	digitPattern = regexp.MustCompile(`^[0-9]+$`)
	// hex sequence, without 0x prefix
	hexadecimalPattern = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	// integer, allowing a sign
	intPattern = regexp.MustCompile(`^[-+]?[0-9]+$`)
	// decimal float, allowing omitted digits around the dot and scientific notation
	floatPattern = regexp.MustCompile(`^[-+]?(?:[0-9]*\.[0-9]+|[0-9]+\.?[0-9]*)(?:[eE][-+]?[0-9]+)?$`)

	// #RGB or #RRGGBB
	hexColorPattern = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	// rgb(r, g, b) with components in 0-255
	rgbColorPattern = regexp.MustCompile(`^rgb\(\s*(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\s*,\s*(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\s*,\s*(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\s*\)$`)

	// single domain label, without dots, length 1-63
	subdomainPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?$`)
	// full domain: multiple labels + alphabetic or internationalized TLD
	domainPattern = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+(?:[a-zA-Z]{1,63}|xn--[a-z0-9]{1,59})$`)
	// DNS name: dot-separated labels, labels not starting/ending with a hyphen, trailing root dot allowed
	dnsNamePattern = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.?$`)
	// E.164 international phone number: starts with +, first digit non-zero, at most 15 digits
	e164Pattern = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	// CORS origin: http(s) + host + optional port, no path
	originPattern = regexp.MustCompile(`^https?://([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(:\d{1,5})?$`)
	// RFC 2397 data URI: optional media type/charset/base64 marker + payload
	dataURIPattern = regexp.MustCompile(`^data:([\w.+-]+/[\w.+-]+)?(;charset=[\w-]+)?(;base64)?,(.*)$`)
	// skeletal format of US social security numbers, area/group details checked by predicates
	ssnPattern = regexp.MustCompile(`^\d{3}-\d{2}-\d{4}$`)
	// official recommended regex of the semantic versioning 2.0.0 spec
	semverPattern = regexp.MustCompile(`^((0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)$`)
	// ULID: 26-character Crockford Base32, excluding the confusable I/L/O/U
	ulidPattern = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Za-hjkmnp-tv-z]{26}$`)
	// MongoDB ObjectID: 24 hex digits
	mongoIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)

	// Full-width/half-width determination uses inclusive semantics:
	// FullWidth = a full-width character exists, HalfWidth = a half-width character exists, VariableWidth = both exist.
	// The half-width set consists of three parts: the printable ASCII range (U+0020 to U+007E),
	// the half-width katakana and vertical-forms range (U+FF61 to U+FFDC etc.), and the hex literal characters 012789abcdef.
	fullWidthSearch = regexp.MustCompile(`[^\x{0020}-\x{007E}\x{FF61}-\x{FF9F}\x{FFA0}-\x{FFDC}\x{FFE8}-\x{FFEE}012789abcdef]`)
	halfWidthSearch = regexp.MustCompile(`[\x{0020}-\x{007E}\x{FF61}-\x{FF9F}\x{FFA0}-\x{FFDC}\x{FFE8}-\x{FFEE}012789abcdef]`)

	// any non-ASCII byte counts as a multibyte character
	multibyteRe = regexp.MustCompile(`[^\x00-\x7F]`)

	// base64url charset (RFC 4648 section 5, without padding)
	base64URLPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	// mainland China mobile number: starts with 1, second digit 3-9, 11 digits in total
	mobileCNPattern = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 18-digit ID card skeleton: 6-digit region + 8-digit birth date + 3-digit sequence + 1 check digit
	idCardCNPattern = regexp.MustCompile(`^[1-9]\d{5}(?:18|19|20)\d{2}(?:0[1-9]|1[0-2])(?:0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]$`)
	// mainland China postal code: 6 digits
	postalCodeCNPattern = regexp.MustCompile(`^\d{6}$`)
	// rgba(r, g, b, a), r/g/b 0-255, a is a 0-1 decimal or 0-100%
	rgbaPattern = regexp.MustCompile(`^rgba\(\s*(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\s*,\s*(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\s*,\s*(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\s*,\s*(0(?:\.\d+)?|1(?:\.0+)?|[0-9]{1,2}%|100%)\s*\)$`)
	// hsl(h, s%, l%), h 0-360, s/l 0-100
	hslPattern = regexp.MustCompile(`^hsl\(\s*(36[0]|3[0-5][0-9]|[12][0-9]{2}|[1-9]?[0-9])\s*,\s*([0-9]{1,2}|100)%\s*,\s*([0-9]{1,2}|100)%\s*\)$`)
	// hsla(h, s%, l%, a), h 0-360, s/l 0-100, a is a 0-1 decimal or 0-100%
	hslaPattern = regexp.MustCompile(`^hsla\(\s*(36[0]|3[0-5][0-9]|[12][0-9]{2}|[1-9]?[0-9])\s*,\s*([0-9]{1,2}|100)%\s*,\s*([0-9]{1,2}|100)%\s*,\s*(0(?:\.\d+)?|1(?:\.0+)?|[0-9]{1,2}%|100%)\s*\)$`)
)
