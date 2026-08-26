package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestUUID(t *testing.T) {
	valid := []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6BA7B810-9DAD-11D1-80B4-00C04FD430C8",
	}
	for _, v := range valid {
		if err := is.UUID(v); err != nil {
			t.Errorf("UUID(%q) = %v, want nil", v, err)
		}
	}

	invalid := []string{"", "6ba7b8109dad11d180b400c04fd430c8", "6ba7b81x-9dad-11d1-80b4-00c04fd430c8"}
	for _, v := range invalid {
		if err := is.UUID(v); err == nil {
			t.Errorf("UUID(%q) should fail", v)
		}
	}
}

func TestBase64(t *testing.T) {
	if err := is.Base64("aGVsbG8="); err != nil {
		t.Errorf("Base64(hello encoded) = %v, want nil", err)
	}
	for _, v := range []string{"", "abc", "not base64!!"} {
		if err := is.Base64(v); err == nil {
			t.Errorf("Base64(%q) should fail", v)
		}
	}
}

func TestJSON(t *testing.T) {
	valid := []string{`{"k": [1, 2]}`, `[]`, `"text"`, `null`}
	for _, v := range valid {
		if err := is.JSON(v); err != nil {
			t.Errorf("JSON(%s) = %v, want nil", v, err)
		}
	}

	invalid := []string{"", `{k: 1}`, `{"unclosed": `}
	for _, v := range invalid {
		if err := is.JSON(v); err == nil {
			t.Errorf("JSON(%s) should fail", v)
		}
	}
}
