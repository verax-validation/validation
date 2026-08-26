package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestWidthRules(t *testing.T) {
	// inclusive semantics: passing when a target character exists, not requiring the whole string
	checkRules(t, "Multibyte", is.Multibyte,
		[]string{"中文", "café", "あいう"},
		[]string{"", "abc123"})

	checkRules(t, "FullWidth", is.FullWidth,
		[]string{"ＨＥＬＬＯ", "１２３", "ａｂｃ"},
		[]string{"", "hello123"})

	checkRules(t, "HalfWidth", is.HalfWidth,
		[]string{"hello123", "!@#", "ｱｲｳ"},
		[]string{"", "１２３４", "中文"})

	checkRules(t, "VariableWidth", is.VariableWidth,
		[]string{"１２３abc", "ＡＢＣ123"},
		[]string{"", "abcdef", "１２３４"})
}
