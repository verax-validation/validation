package rules

import "github.com/verax-validation/validation/internal/codes"

// Error code constants, whose values are the external contract, referenced by translation tables and programmatic checks.

const (
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
)
