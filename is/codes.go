package is

import "github.com/verax-validation/validation/internal/codes"

// Error code constants, whose values are the external contract, referenced by translation tables and programmatic checks.

const (
	// CodeEmail's value is in internal/codes, representing email-related assertion failures.
	CodeEmail = codes.CodeEmail
	// CodeAlpha's value is in internal/codes, representing alpha-related assertion failures.
	CodeAlpha = codes.CodeAlpha
	// CodeAlphanumeric's value is in internal/codes, representing alphanumeric-related assertion failures.
	CodeAlphanumeric = codes.CodeAlphanumeric
	// CodeDigit's value is in internal/codes, representing digit-related assertion failures.
	CodeDigit = codes.CodeDigit
	// CodeUTFLetter's value is in internal/codes, representing utf_letter-related assertion failures.
	CodeUTFLetter = codes.CodeUTFLetter
	// CodeUTFDigit's value is in internal/codes, representing utf_digit-related assertion failures.
	CodeUTFDigit = codes.CodeUTFDigit
	// CodeUTFNumeric's value is in internal/codes, representing utf_numeric-related assertion failures.
	CodeUTFNumeric = codes.CodeUTFNumeric
	// CodeUTFLetterNumeric's value is in internal/codes, representing utf_letter_numeric-related assertion failures.
	CodeUTFLetterNumeric = codes.CodeUTFLetterNumeric
	// CodeLowerCase's value is in internal/codes, representing lower_case-related assertion failures.
	CodeLowerCase = codes.CodeLowerCase
	// CodeUpperCase's value is in internal/codes, representing upper_case-related assertion failures.
	CodeUpperCase = codes.CodeUpperCase
	// CodeASCII's value is in internal/codes, representing ascii-related assertion failures.
	CodeASCII = codes.CodeASCII
	// CodePrintableASCII's value is in internal/codes, representing printable_ascii-related assertion failures.
	CodePrintableASCII = codes.CodePrintableASCII
	// CodeHexadecimal's value is in internal/codes, representing hexadecimal-related assertion failures.
	CodeHexadecimal = codes.CodeHexadecimal
	// CodeHexColor's value is in internal/codes, representing hex_color-related assertion failures.
	CodeHexColor = codes.CodeHexColor
	// CodeRGBColor's value is in internal/codes, representing rgb_color-related assertion failures.
	CodeRGBColor = codes.CodeRGBColor
	// CodeInt's value is in internal/codes, representing int-related assertion failures.
	CodeInt = codes.CodeInt
	// CodeFloat's value is in internal/codes, representing float-related assertion failures.
	CodeFloat = codes.CodeFloat
	// CodeUUID's value is in internal/codes, representing uuid-related assertion failures.
	CodeUUID = codes.CodeUUID
	// CodeUUIDv3's value is in internal/codes, representing uuid_v3-related assertion failures.
	CodeUUIDv3 = codes.CodeUUIDv3
	// CodeUUIDv4's value is in internal/codes, representing uuid_v4-related assertion failures.
	CodeUUIDv4 = codes.CodeUUIDv4
	// CodeUUIDv5's value is in internal/codes, representing uuid_v5-related assertion failures.
	CodeUUIDv5 = codes.CodeUUIDv5
	// CodeUUIDv7's value is in internal/codes, representing uuid_v7-related assertion failures.
	CodeUUIDv7 = codes.CodeUUIDv7
	// CodeULID's value is in internal/codes, representing ulid-related assertion failures.
	CodeULID = codes.CodeULID
	// CodeBase64's value is in internal/codes, representing base64-related assertion failures.
	CodeBase64 = codes.CodeBase64
	// CodeJSON's value is in internal/codes, representing json-related assertion failures.
	CodeJSON = codes.CodeJSON
	// CodeURL's value is in internal/codes, representing url-related assertion failures.
	CodeURL = codes.CodeURL
	// CodeRequestURL's value is in internal/codes, representing request_url-related assertion failures.
	CodeRequestURL = codes.CodeRequestURL
	// CodeRequestURI's value is in internal/codes, representing request_uri-related assertion failures.
	CodeRequestURI = codes.CodeRequestURI
	// CodeIP's value is in internal/codes, representing ip-related assertion failures.
	CodeIP = codes.CodeIP
	// CodeIPv4's value is in internal/codes, representing ipv4-related assertion failures.
	CodeIPv4 = codes.CodeIPv4
	// CodeIPv6's value is in internal/codes, representing ipv6-related assertion failures.
	CodeIPv6 = codes.CodeIPv6
	// CodeMAC's value is in internal/codes, representing mac-related assertion failures.
	CodeMAC = codes.CodeMAC
	// CodeCreditCard's value is in internal/codes, representing credit_card-related assertion failures.
	CodeCreditCard = codes.CodeCreditCard
	// CodeISBN10's value is in internal/codes, representing isbn_10-related assertion failures.
	CodeISBN10 = codes.CodeISBN10
	// CodeISBN13's value is in internal/codes, representing isbn_13-related assertion failures.
	CodeISBN13 = codes.CodeISBN13
	// CodeISBN's value is in internal/codes, representing isbn-related assertion failures.
	CodeISBN = codes.CodeISBN
	// CodeSubdomain's value is in internal/codes, representing sub_domain-related assertion failures.
	CodeSubdomain = codes.CodeSubdomain
	// CodeDomain's value is in internal/codes, representing domain-related assertion failures.
	CodeDomain = codes.CodeDomain
	// CodeDNSName's value is in internal/codes, representing dns_name-related assertion failures.
	CodeDNSName = codes.CodeDNSName
	// CodeHost's value is in internal/codes, representing host-related assertion failures.
	CodeHost = codes.CodeHost
	// CodePort's value is in internal/codes, representing port-related assertion failures.
	CodePort = codes.CodePort
	// CodeE164's value is in internal/codes, representing e164_number-related assertion failures.
	CodeE164 = codes.CodeE164
	// CodeDialString's value is in internal/codes, representing dial_string-related assertion failures.
	CodeDialString = codes.CodeDialString
	// CodeMongoID's value is in internal/codes, representing mongo_id-related assertion failures.
	CodeMongoID = codes.CodeMongoID
	// CodeDataURI's value is in internal/codes, representing data_uri-related assertion failures.
	CodeDataURI = codes.CodeDataURI
	// CodeLatitude's value is in internal/codes, representing latitude-related assertion failures.
	CodeLatitude = codes.CodeLatitude
	// CodeLongitude's value is in internal/codes, representing longitude-related assertion failures.
	CodeLongitude = codes.CodeLongitude
	// CodeSSN's value is in internal/codes, representing ssn-related assertion failures.
	CodeSSN = codes.CodeSSN
	// CodeSemver's value is in internal/codes, representing semver-related assertion failures.
	CodeSemver = codes.CodeSemver
	// CodeOrigin's value is in internal/codes, representing origin-related assertion failures.
	CodeOrigin = codes.CodeOrigin
	// CodeMultibyte's value is in internal/codes, representing multibyte-related assertion failures.
	CodeMultibyte = codes.CodeMultibyte
	// CodeFullWidth's value is in internal/codes, representing full_width-related assertion failures.
	CodeFullWidth = codes.CodeFullWidth
	// CodeHalfWidth's value is in internal/codes, representing half_width-related assertion failures.
	CodeHalfWidth = codes.CodeHalfWidth
	// CodeVariableWidth's value is in internal/codes, representing variable_width-related assertion failures.
	CodeVariableWidth = codes.CodeVariableWidth
	// CodeBase64URL's value is in internal/codes, representing base64url-related assertion failures.
	CodeBase64URL = codes.CodeBase64URL
	// CodeJWT's value is in internal/codes, representing jwt-related assertion failures.
	CodeJWT = codes.CodeJWT
	// CodeBoolean's value is in internal/codes, representing boolean-related assertion failures.
	CodeBoolean = codes.CodeBoolean
	// CodeTimeZone's value is in internal/codes, representing timezone-related assertion failures.
	CodeTimeZone = codes.CodeTimeZone
	// CodeCountryCode's value is in internal/codes, representing country_code-related assertion failures.
	CodeCountryCode = codes.CodeCountryCode
	// CodeCurrencyCode's value is in internal/codes, representing currency_code-related assertion failures.
	CodeCurrencyCode = codes.CodeCurrencyCode
	// CodeLanguageCode's value is in internal/codes, representing language_code-related assertion failures.
	CodeLanguageCode = codes.CodeLanguageCode
	// CodeCIDR's value is in internal/codes, representing cidr-related assertion failures.
	CodeCIDR = codes.CodeCIDR
	// CodeCIDRv4's value is in internal/codes, representing cidr_v4-related assertion failures.
	CodeCIDRv4 = codes.CodeCIDRv4
	// CodeCIDRv6's value is in internal/codes, representing cidr_v6-related assertion failures.
	CodeCIDRv6 = codes.CodeCIDRv6
	// CodeRGBA's value is in internal/codes, representing rgba-related assertion failures.
	CodeRGBA = codes.CodeRGBA
	// CodeHSL's value is in internal/codes, representing hsl-related assertion failures.
	CodeHSL = codes.CodeHSL
	// CodeHSLA's value is in internal/codes, representing hsla-related assertion failures.
	CodeHSLA = codes.CodeHSLA
	// CodeMobileCN's value is in internal/codes, representing mobile_cn-related assertion failures.
	CodeMobileCN = codes.CodeMobileCN
	// CodeIDCardCN's value is in internal/codes, representing id_card_cn-related assertion failures.
	CodeIDCardCN = codes.CodeIDCardCN
	// CodePostalCodeCN's value is in internal/codes, representing postal_code_cn-related assertion failures.
	CodePostalCodeCN = codes.CodePostalCodeCN
)
