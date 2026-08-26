// Package codes centrally defines all error code constants.
// Error code values are the external contract and must not change after release;
// public packages re-export them by domain, and users should reference the re-exported constants rather than literals.
package codes

const (
	// CodeStructNil is the error code returned when the object being validated is nil
	CodeStructNil = "validation.struct_nil"

	// CodeRequired is the error code for a failed required-field validation
	CodeRequired = "validation.required"
	// CodeLength is the error code for a length out of range
	CodeLength = "validation.length"
	// CodeMin is the error code for being below the lower bound
	CodeMin = "validation.min"
	// CodeMax is the error code for being above the upper bound
	CodeMax = "validation.max"
	// CodeBetween is the error code for being outside a range
	CodeBetween = "validation.between"
	// CodeIn is the error code for not being in the allowed value list
	CodeIn = "validation.in"
	// CodeNotIn is the error code for hitting a forbidden value in the list
	CodeNotIn = "validation.not_in"
	// CodeMatch is the error code for a failed regex match
	CodeMatch = "validation.match"
	// CodeDate is the error code for an invalid date format
	CodeDate = "validation.date"
	// CodeEq is the error code for a value not equal to the target
	CodeEq = "validation.eq"
	// CodeNe is the error code for a value hitting a forbidden target
	CodeNe = "validation.ne"
	// CodeGt is the error code for not being greater than the lower bound
	CodeGt = "validation.gt"
	// CodeLt is the error code for not being less than the upper bound
	CodeLt = "validation.lt"
	// CodeContains is the error code for not containing the given substring
	CodeContains = "validation.contains"
	// CodeStartWith is the error code for not starting with the given prefix
	CodeStartWith = "validation.start_with"
	// CodeEndWith is the error code for not ending with the given suffix
	CodeEndWith = "validation.end_with"
	// CodeExcludes is the error code for containing a forbidden substring
	CodeExcludes = "validation.excludes"
	// CodeContainsAny is the error code for not containing any of the given characters
	CodeContainsAny = "validation.contains_any"
	// CodeNotNil is the error code for a nil value
	CodeNotNil = "validation.not_nil"
	// CodeExactLen is the error code for a length not equal to the given value
	CodeExactLen = "validation.exact_len"
	// CodeMultipleOf is the error code for not being a multiple of the given base
	CodeMultipleOf = "validation.multiple_of"
	// CodeFieldEq is the error code for a failed cross-field equality comparison
	CodeFieldEq = "validation.field_eq"
	// CodeFieldNe is the error code for a failed cross-field inequality comparison
	CodeFieldNe = "validation.field_ne"
	// CodeFieldGt is the error code for a failed cross-field greater-than comparison
	CodeFieldGt = "validation.field_gt"
	// CodeFieldGte is the error code for a failed cross-field greater-or-equal comparison
	CodeFieldGte = "validation.field_gte"
	// CodeFieldLt is the error code for a failed cross-field less-than comparison
	CodeFieldLt = "validation.field_lt"
	// CodeFieldLte is the error code for a failed cross-field less-or-equal comparison
	CodeFieldLte = "validation.field_lte"
	// CodeCollectionLen is the error code for a collection size out of range
	CodeCollectionLen = "validation.collections.len"
	// CodeCollectionUnique is the error code for duplicate elements in a collection
	CodeCollectionUnique = "validation.collections.unique"

	// CodeEmail is the error code for a failed email format assertion
	CodeEmail = "validation.is.email"
	// CodeAlpha is the error code for a failed English-letters-only assertion
	CodeAlpha = "validation.is.alpha"
	// CodeAlphanumeric is the error code for a failed English-letters-and-digits-only assertion
	CodeAlphanumeric = "validation.is.alphanumeric"
	// CodeDigit is the error code for a failed digits-only assertion
	CodeDigit = "validation.is.digit"
	// CodeUTFLetter is the error code for a failed unicode-letters-only assertion
	CodeUTFLetter = "validation.is.utf_letter"
	// CodeUTFDigit is the error code for a failed unicode-decimal-digits-only assertion
	CodeUTFDigit = "validation.is.utf_digit"
	// CodeUTFNumeric is the error code for a failed unicode-number-only assertion
	CodeUTFNumeric = "validation.is.utf_numeric"
	// CodeUTFLetterNumeric is the error code for a failed unicode-letters-or-numbers-only assertion
	CodeUTFLetterNumeric = "validation.is.utf_letter_numeric"
	// CodeLowerCase is the error code for a failed all-lowercase assertion
	CodeLowerCase = "validation.is.lower_case"
	// CodeUpperCase is the error code for a failed all-uppercase assertion
	CodeUpperCase = "validation.is.upper_case"
	// CodeASCII is the error code for a failed ASCII-only assertion
	CodeASCII = "validation.is.ascii"
	// CodePrintableASCII is the error code for a failed printable-ASCII-only assertion
	CodePrintableASCII = "validation.is.printable_ascii"
	// CodeHexadecimal is the error code for a failed hexadecimal-number assertion
	CodeHexadecimal = "validation.is.hexadecimal"
	// CodeHexColor is the error code for a failed hex-color-code assertion
	CodeHexColor = "validation.is.hex_color"
	// CodeRGBColor is the error code for a failed RGB-color assertion
	CodeRGBColor = "validation.is.rgb_color"
	// CodeInt is the error code for a failed integer assertion
	CodeInt = "validation.is.int"
	// CodeFloat is the error code for a failed float assertion
	CodeFloat = "validation.is.float"
	// CodeUUID is the error code for a failed generic-UUID assertion
	CodeUUID = "validation.is.uuid"
	// CodeUUIDv3 is the error code for a failed UUID v3 assertion
	CodeUUIDv3 = "validation.is.uuid_v3"
	// CodeUUIDv4 is the error code for a failed UUID v4 assertion
	CodeUUIDv4 = "validation.is.uuid_v4"
	// CodeUUIDv5 is the error code for a failed UUID v5 assertion
	CodeUUIDv5 = "validation.is.uuid_v5"
	// CodeUUIDv7 is the error code for a failed UUID v7 assertion
	CodeUUIDv7 = "validation.is.uuid_v7"
	// CodeULID is the error code for a failed ULID assertion
	CodeULID = "validation.is.ulid"
	// CodeBase64 is the error code for a failed base64 assertion
	CodeBase64 = "validation.is.base64"
	// CodeJSON is the error code for a failed JSON-format assertion
	CodeJSON = "validation.is.json"
	// CodeURL is the error code for a failed URL assertion
	CodeURL = "validation.is.url"
	// CodeRequestURL is the error code for a failed request-URL assertion
	CodeRequestURL = "validation.is.request_url"
	// CodeRequestURI is the error code for a failed request-URI assertion
	CodeRequestURI = "validation.is.request_uri"
	// CodeIP is the error code for a failed IP assertion
	CodeIP = "validation.is.ip"
	// CodeIPv4 is the error code for a failed IPv4 assertion
	CodeIPv4 = "validation.is.ipv4"
	// CodeIPv6 is the error code for a failed IPv6 assertion
	CodeIPv6 = "validation.is.ipv6"
	// CodeMAC is the error code for a failed MAC assertion
	CodeMAC = "validation.is.mac"
	// CodeCreditCard is the error code for a failed credit card number assertion
	CodeCreditCard = "validation.is.credit_card"
	// CodeISBN10 is the error code for a failed ISBN-10 assertion
	CodeISBN10 = "validation.is.isbn_10"
	// CodeISBN13 is the error code for a failed ISBN-13 assertion
	CodeISBN13 = "validation.is.isbn_13"
	// CodeISBN is the error code for a failed ISBN assertion
	CodeISBN = "validation.is.isbn"
	// CodeSubdomain is the error code for a failed subdomain assertion
	CodeSubdomain = "validation.is.sub_domain"
	// CodeDomain is the error code for a failed domain assertion
	CodeDomain = "validation.is.domain"
	// CodeDNSName is the error code for a failed DNS-name assertion
	CodeDNSName = "validation.is.dns_name"
	// CodeHost is the error code for a failed host assertion
	CodeHost = "validation.is.host"
	// CodePort is the error code for a failed port assertion
	CodePort = "validation.is.port"
	// CodeE164 is the error code for a failed E.164 number assertion
	CodeE164 = "validation.is.e164_number"
	// CodeDialString is the error code for a failed dial-address assertion
	CodeDialString = "validation.is.dial_string"
	// CodeMongoID is the error code for a failed MongoDB ObjectID assertion
	CodeMongoID = "validation.is.mongo_id"
	// CodeDataURI is the error code for a failed Data URI assertion
	CodeDataURI = "validation.is.data_uri"
	// CodeBase64URL is the error code for a failed base64url assertion
	CodeBase64URL = "validation.is.base64_url"
	// CodeJWT is the error code for a failed JWT assertion
	CodeJWT = "validation.is.jwt"
	// CodeBoolean is the error code for a failed boolean assertion
	CodeBoolean = "validation.is.boolean"
	// CodeTimeZone is the error code for a failed timezone assertion
	CodeTimeZone = "validation.is.timezone"
	// CodeCountryCode is the error code for a failed country-code assertion
	CodeCountryCode = "validation.is.country_code"
	// CodeCurrencyCode is the error code for a failed currency-code assertion
	CodeCurrencyCode = "validation.is.currency_code"
	// CodeLanguageCode is the error code for a failed language-code assertion
	CodeLanguageCode = "validation.is.language_code"
	// CodeCIDR is the error code for a failed CIDR assertion
	CodeCIDR = "validation.is.cidr"
	// CodeCIDRv4 is the error code for a failed IPv4 CIDR assertion
	CodeCIDRv4 = "validation.is.cidr_v4"
	// CodeCIDRv6 is the error code for a failed IPv6 CIDR assertion
	CodeCIDRv6 = "validation.is.cidr_v6"
	// CodeRGBA is the error code for a failed RGBA-color assertion
	CodeRGBA = "validation.is.rgba"
	// CodeHSL is the error code for a failed HSL-color assertion
	CodeHSL = "validation.is.hsl"
	// CodeHSLA is the error code for a failed HSLA-color assertion
	CodeHSLA = "validation.is.hsla"
	// CodeMobileCN is the error code for a failed mainland China mobile number assertion
	CodeMobileCN = "validation.is.mobile_cn"
	// CodeIDCardCN is the error code for a failed mainland China ID card number assertion
	CodeIDCardCN = "validation.is.id_card_cn"
	// CodePostalCodeCN is the error code for a failed mainland China postal code assertion
	CodePostalCodeCN = "validation.is.postal_code_cn"
	// CodeLatitude is the error code for a failed latitude assertion
	CodeLatitude = "validation.is.latitude"
	// CodeLongitude is the error code for a failed longitude assertion
	CodeLongitude = "validation.is.longitude"
	// CodeSSN is the error code for a failed social security number assertion
	CodeSSN = "validation.is.ssn"
	// CodeSemver is the error code for a failed semantic version assertion
	CodeSemver = "validation.is.semver"
	// CodeOrigin is the error code for a failed origin assertion
	CodeOrigin = "validation.is.origin"
	// CodeMultibyte is the error code for a failed multibyte-character assertion
	CodeMultibyte = "validation.is.multibyte"
	// CodeFullWidth is the error code for a failed full-width-character assertion
	CodeFullWidth = "validation.is.full_width"
	// CodeHalfWidth is the error code for a failed half-width-character assertion
	CodeHalfWidth = "validation.is.half_width"
	// CodeVariableWidth is the error code for a failed both-full-and-half-width assertion
	CodeVariableWidth = "validation.is.variable_width"
)
