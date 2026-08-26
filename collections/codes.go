package collections

import "github.com/verax-validation/validation/internal/codes"

const (
	// CodeLen is the error code for a collection size out of range.
	CodeLen = codes.CodeCollectionLen
	// CodeUnique is the error code for duplicate elements in a collection.
	CodeUnique = codes.CodeCollectionUnique
)
