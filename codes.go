package verax

import "github.com/verax-validation/validation/internal/codes"

// Error code constants, whose values are the external contract, referenced by translation tables and programmatic checks.
// All constants are defined centrally in internal/codes and re-exported here by domain;
// users should reference the exported names from this package or the rule packages, not literals.
const (
	// CodeStructNil is the error code returned when the object being validated is nil
	CodeStructNil = codes.CodeStructNil

	// CodeRequired is the error code for a failed required-field validation
	CodeRequired = codes.CodeRequired
	// CodeLength is the error code for a length out of range
	CodeLength = codes.CodeLength
	// CodeMin is the error code for being below the lower bound
	CodeMin = codes.CodeMin
	// CodeMax is the error code for being above the upper bound
	CodeMax = codes.CodeMax
	// CodeBetween is the error code for being outside a range
	CodeBetween = codes.CodeBetween
	// CodeIn is the error code for not being in the allowed value list
	CodeIn = codes.CodeIn
	// CodeNotIn is the error code for hitting a forbidden value in the list
	CodeNotIn = codes.CodeNotIn
	// CodeMatch is the error code for a failed regex match
	CodeMatch = codes.CodeMatch
	// CodeDate is the error code for an invalid date format
	CodeDate = codes.CodeDate
	// CodeEq is the error code for a value not equal to the target
	CodeEq = codes.CodeEq
	// CodeNe is the error code for a value hitting a forbidden target
	CodeNe = codes.CodeNe
	// CodeGt is the error code for not being greater than the lower bound
	CodeGt = codes.CodeGt
	// CodeLt is the error code for not being less than the upper bound
	CodeLt = codes.CodeLt
	// CodeContains is the error code for not containing the given substring
	CodeContains = codes.CodeContains
	// CodeStartWith is the error code for not starting with the given prefix
	CodeStartWith = codes.CodeStartWith
	// CodeEndWith is the error code for not ending with the given suffix
	CodeEndWith = codes.CodeEndWith
	// CodeExcludes is the error code for containing a forbidden substring
	CodeExcludes = codes.CodeExcludes
	// CodeContainsAny is the error code for not containing any of the given characters
	CodeContainsAny = codes.CodeContainsAny
	// CodeNotNil is the error code for a nil value
	CodeNotNil = codes.CodeNotNil
	// CodeExactLen is the error code for a length not equal to the given value
	CodeExactLen = codes.CodeExactLen
	// CodeMultipleOf is the error code for not being a multiple of the given base
	CodeMultipleOf = codes.CodeMultipleOf
	// CodeFieldEq is the error code for a failed cross-field equality comparison
	CodeFieldEq = codes.CodeFieldEq
	// CodeFieldNe is the error code for a failed cross-field inequality comparison
	CodeFieldNe = codes.CodeFieldNe
	// CodeFieldGt is the error code for a failed cross-field greater-than comparison
	CodeFieldGt = codes.CodeFieldGt
	// CodeFieldGte is the error code for a failed cross-field greater-or-equal comparison
	CodeFieldGte = codes.CodeFieldGte
	// CodeFieldLt is the error code for a failed cross-field less-than comparison
	CodeFieldLt = codes.CodeFieldLt
	// CodeFieldLte is the error code for a failed cross-field less-or-equal comparison
	CodeFieldLte = codes.CodeFieldLte
	// CodeCollectionLen is the error code for a collection size out of range
	CodeCollectionLen = codes.CodeCollectionLen
	// CodeCollectionUnique is the error code for duplicate elements in a collection
	CodeCollectionUnique = codes.CodeCollectionUnique

	// CodeEmail is the error code for a failed email format assertion
	CodeEmail = codes.CodeEmail
	// CodeAlpha is the error code for a failed English-letters-only assertion
	CodeAlpha = codes.CodeAlpha
	// CodeAlphanumeric is the error code for a failed English-letters-and-digits-only assertion
	CodeAlphanumeric = codes.CodeAlphanumeric
	// CodeDigit is the error code for a failed digits-only assertion
	CodeDigit = codes.CodeDigit
	// CodeUTFLetter is the error code for a failed unicode-letters-only assertion
	CodeUTFLetter = codes.CodeUTFLetter
	// CodeUTFDigit is the error code for a failed unicode-decimal-digits-only assertion
	CodeUTFDigit = codes.CodeUTFDigit
	// CodeUTFNumeric is the error code for a failed unicode-number-only assertion
	CodeUTFNumeric = codes.CodeUTFNumeric
	// CodeUTFLetterNumeric is the error code for a failed unicode-letters-or-numbers-only assertion
	CodeUTFLetterNumeric = codes.CodeUTFLetterNumeric
	// CodeLowerCase is the error code for a failed all-lowercase assertion
	CodeLowerCase = codes.CodeLowerCase
	// CodeUpperCase is the error code for a failed all-uppercase assertion
	CodeUpperCase = codes.CodeUpperCase
	// CodeASCII is the error code for a failed ASCII-only assertion
	CodeASCII = codes.CodeASCII
	// CodePrintableASCII is the error code for a failed printable-ASCII-only assertion
	CodePrintableASCII = codes.CodePrintableASCII
	// CodeHexadecimal is the error code for a failed hexadecimal-number assertion
	CodeHexadecimal = codes.CodeHexadecimal
	// CodeHexColor is the error code for a failed hex-color-code assertion
	CodeHexColor = codes.CodeHexColor
	// CodeRGBColor is the error code for a failed RGB-color assertion
	CodeRGBColor = codes.CodeRGBColor
	// CodeInt is the error code for a failed integer assertion
	CodeInt = codes.CodeInt
	// CodeFloat is the error code for a failed float assertion
	CodeFloat = codes.CodeFloat
	// CodeUUID is the error code for a failed UUID assertion
	CodeUUID = codes.CodeUUID
	// CodeUUIDv3 is the error code for a failed UUID v3 assertion
	CodeUUIDv3 = codes.CodeUUIDv3
	// CodeUUIDv4 is the error code for a failed UUID v4 assertion
	CodeUUIDv4 = codes.CodeUUIDv4
	// CodeUUIDv5 is the error code for a failed UUID v5 assertion
	CodeUUIDv5 = codes.CodeUUIDv5
	// CodeUUIDv7 is the error code for a failed UUID v7 assertion
	CodeUUIDv7 = codes.CodeUUIDv7
	// CodeULID is the error code for a failed ULID assertion
	CodeULID = codes.CodeULID
	// CodeBase64 is the error code for a failed base64 assertion
	CodeBase64 = codes.CodeBase64
	// CodeJSON is the error code for a failed JSON-format assertion
	CodeJSON = codes.CodeJSON
	// CodeURL is the error code for a failed URL assertion
	CodeURL = codes.CodeURL
	// CodeRequestURL is the error code for a failed request-URL assertion
	CodeRequestURL = codes.CodeRequestURL
	// CodeRequestURI is the error code for a failed request-URI assertion
	CodeRequestURI = codes.CodeRequestURI
	// CodeIP is the error code for a failed IP assertion
	CodeIP = codes.CodeIP
	// CodeIPv4 is the error code for a failed IPv4 assertion
	CodeIPv4 = codes.CodeIPv4
	// CodeIPv6 is the error code for a failed IPv6 assertion
	CodeIPv6 = codes.CodeIPv6
	// CodeMAC is the error code for a failed MAC assertion
	CodeMAC = codes.CodeMAC
	// CodeCreditCard is the error code for a failed credit card number assertion
	CodeCreditCard = codes.CodeCreditCard
	// CodeISBN10 is the error code for a failed ISBN-10 assertion
	CodeISBN10 = codes.CodeISBN10
	// CodeISBN13 is the error code for a failed ISBN-13 assertion
	CodeISBN13 = codes.CodeISBN13
	// CodeISBN is the error code for a failed ISBN assertion
	CodeISBN = codes.CodeISBN
	// CodeSubdomain is the error code for a failed subdomain assertion
	CodeSubdomain = codes.CodeSubdomain
	// CodeDomain is the error code for a failed domain assertion
	CodeDomain = codes.CodeDomain
	// CodeDNSName is the error code for a failed DNS-name assertion
	CodeDNSName = codes.CodeDNSName
	// CodeHost is the error code for a failed host assertion
	CodeHost = codes.CodeHost
	// CodePort is the error code for a failed port assertion
	CodePort = codes.CodePort
	// CodeE164 is the error code for a failed E.164 number assertion
	CodeE164 = codes.CodeE164
	// CodeDialString is the error code for a failed dial-address assertion
	CodeDialString = codes.CodeDialString
	// CodeMongoID is the error code for a failed MongoDB ObjectID assertion
	CodeMongoID = codes.CodeMongoID
	// CodeDataURI is the error code for a failed Data URI assertion
	CodeDataURI = codes.CodeDataURI
	// CodeBase64URL is the error code for a failed base64url assertion
	CodeBase64URL = codes.CodeBase64URL
	// CodeJWT is the error code for a failed JWT assertion
	CodeJWT = codes.CodeJWT
	// CodeBoolean is the error code for a failed boolean assertion
	CodeBoolean = codes.CodeBoolean
	// CodeTimeZone is the error code for a failed timezone assertion
	CodeTimeZone = codes.CodeTimeZone
	// CodeCountryCode is the error code for a failed country-code assertion
	CodeCountryCode = codes.CodeCountryCode
	// CodeCurrencyCode is the error code for a failed currency-code assertion
	CodeCurrencyCode = codes.CodeCurrencyCode
	// CodeLanguageCode is the error code for a failed language-code assertion
	CodeLanguageCode = codes.CodeLanguageCode
	// CodeCIDR is the error code for a failed CIDR assertion
	CodeCIDR = codes.CodeCIDR
	// CodeCIDRv4 is the error code for a failed IPv4 CIDR assertion
	CodeCIDRv4 = codes.CodeCIDRv4
	// CodeCIDRv6 is the error code for a failed IPv6 CIDR assertion
	CodeCIDRv6 = codes.CodeCIDRv6
	// CodeRGBA is the error code for a failed RGBA-color assertion
	CodeRGBA = codes.CodeRGBA
	// CodeHSL is the error code for a failed HSL-color assertion
	CodeHSL = codes.CodeHSL
	// CodeHSLA is the error code for a failed HSLA-color assertion
	CodeHSLA = codes.CodeHSLA
	// CodeMobileCN is the error code for a failed mainland China mobile number assertion
	CodeMobileCN = codes.CodeMobileCN
	// CodeIDCardCN is the error code for a failed mainland China ID card number assertion
	CodeIDCardCN = codes.CodeIDCardCN
	// CodePostalCodeCN is the error code for a failed mainland China postal code assertion
	CodePostalCodeCN = codes.CodePostalCodeCN
	// CodeLatitude is the error code for a failed latitude assertion
	CodeLatitude = codes.CodeLatitude
	// CodeLongitude is the error code for a failed longitude assertion
	CodeLongitude = codes.CodeLongitude
	// CodeSSN is the error code for a failed social security number assertion
	CodeSSN = codes.CodeSSN
	// CodeSemver is the error code for a failed semantic version assertion
	CodeSemver = codes.CodeSemver
	// CodeOrigin is the error code for a failed origin assertion
	CodeOrigin = codes.CodeOrigin
	// CodeMultibyte is the error code for a failed multibyte-character assertion
	CodeMultibyte = codes.CodeMultibyte
	// CodeFullWidth is the error code for a failed full-width-character assertion
	CodeFullWidth = codes.CodeFullWidth
	// CodeHalfWidth is the error code for a failed half-width-character assertion
	CodeHalfWidth = codes.CodeHalfWidth
	// CodeVariableWidth is the error code for a failed both-full-and-half-width assertion
	CodeVariableWidth = codes.CodeVariableWidth
)
