package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestIntFloat(t *testing.T) {
	checkRules(t, "Int", is.Int,
		[]string{"0", "42", "-7", "+123"},
		[]string{"", "1.5", "1e3", "abc"})

	checkRules(t, "Float", is.Float,
		[]string{"0.5", "-1.25", ".5", "5.", "1e3", "-2.5E-2"},
		[]string{"", "abc", "1,5"})
}

func TestColorRules(t *testing.T) {
	checkRules(t, "HexColor", is.HexColor,
		[]string{"#fff", "#FFFFFF", "#AbCdEf"},
		[]string{"", "fff", "#ffff", "#GGGGGG", "#12345"})

	checkRules(t, "RGBColor", is.RGBColor,
		[]string{"rgb(0, 0, 0)", "rgb(255,255,255)", "rgb(12 , 34 , 56)"},
		[]string{"", "rgb(256, 0, 0)", "rgba(0, 0, 0, 1)", "rgb(-1, 0, 0)"})
}
