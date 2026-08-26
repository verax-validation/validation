package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestBase64URL(t *testing.T) {
	checkRules(t, "Base64URL", is.Base64URL,
		[]string{"aGVsbG8", "U3Vubnk", "aGVsbG8_aGVsbG8"},
		[]string{"", "aGVsbG8=", "aGVs bG8", "aGVsbG8+bA", "a"})
}

func TestJWT(t *testing.T) {
	checkRules(t, "JWT", is.JWT,
		[]string{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
		},
		[]string{"", "abc", "a.b.c.d", "..", "a.bc.def!", "a.b.c"})
}

func TestBoolean(t *testing.T) {
	checkRules(t, "Boolean", is.Boolean,
		[]string{"true", "TRUE", "True", "false", "1", "0", "t", "f"},
		[]string{"", "yes", "2", "on"})
}

func TestTimeZone(t *testing.T) {
	checkRules(t, "TimeZone", is.TimeZone,
		[]string{"UTC", "Asia/Shanghai"},
		[]string{"", "Mars/Olympus"})
}

func TestRGBA(t *testing.T) {
	checkRules(t, "RGBA", is.RGBA,
		[]string{"rgba(255, 0, 0, 1)", "rgba(0,0,0,0.5)", "rgba(100, 200, 255, 50%)", "rgba(0, 0, 0, 0)"},
		[]string{"", "rgb(255,0,0)", "rgba(256, 0, 0, 1)", "rgba(0, 0, 0, 2)", "rgba(0, 0, 0, 1.5)"})
}

func TestHSL(t *testing.T) {
	checkRules(t, "HSL", is.HSL,
		[]string{"hsl(0, 100%, 50%)", "hsl(120, 50%, 0%)", "hsl(360, 0%, 100%)", "hsl(0,0%,0%)"},
		[]string{"", "hsl(361, 100%, 50%)", "hsl(0, 100, 50%)", "hsl(0, 100%, 50)"})
}

func TestHSLA(t *testing.T) {
	checkRules(t, "HSLA", is.HSLA,
		[]string{"hsla(120, 50%, 50%, 1)", "hsla(0, 0%, 0%, 0.5)", "hsla(0, 0%, 0%, 100%)"},
		[]string{"", "hsla(120, 50%, 50%)", "hsla(0, 0%, 0%, 2)", "hsla(0, 0%, 0%, 101%)"})
}
