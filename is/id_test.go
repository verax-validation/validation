package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestUUIDVersionRules(t *testing.T) {
	checkRules(t, "UUIDv3", is.UUIDv3,
		[]string{"b4c6a8e4-c2b1-3f5a-9c0d-1e2f3a4b5c6d"},
		[]string{"", "4f8c2a22-9e57-4b1a-9c0d-1e2f3a4b5c6d", "zzz"})

	checkRules(t, "UUIDv4", is.UUIDv4,
		[]string{"4f8c2a22-9e57-4b1a-9c0d-1e2f3a4b5c6d"},
		[]string{"", "b4c6a8e4-c2b1-3f5a-9c0d-1e2f3a4b5c6d"})

	checkRules(t, "UUIDv5", is.UUIDv5,
		[]string{"74738ff5-5367-5958-9aee-98fffdcd1876"},
		[]string{""})

	checkRules(t, "UUIDv7", is.UUIDv7,
		[]string{"01890a5d-ac96-774b-bcce-b302099a8057"},
		[]string{"", "4f8c2a22-9e57-4b1a-cccc-cccccccccccc"})
}

func TestULIDAndMongoID(t *testing.T) {
	checkRules(t, "ULID", is.ULID,
		[]string{"01ARZ3NDEKTSV4RRFFQ69G5FAV", "01arz3ndektsv4rrffq69g5fav"},
		[]string{"", "01ARZ3NDEKTSV4RRFFQ69G5FA", "ILOU" + "ARZ3NDEKTSV4RRFFQ69G5F"})

	checkRules(t, "MongoID", is.MongoID,
		[]string{"507f1f77bcf86cd799439011"},
		[]string{"", "507f1f77bcf86cd79943901", "507g1f77bcf86cd799439011"})
}
