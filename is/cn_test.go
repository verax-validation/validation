package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestMobileCN(t *testing.T) {
	checkRules(t, "MobileCN", is.MobileCN,
		[]string{"13812345678", "19912345678", "15012345678"},
		[]string{"", "12345678901", "12812345678", "1381234567", "138123456789"})
}

func TestIDCardCN(t *testing.T) {
	checkRules(t, "IDCardCN", is.IDCardCN,
		[]string{"11010519491231002X", "110105199003078887"},
		[]string{"", "11010519491231002", "110105194912310021", "11010519490230002X", "123"})
}

func TestPostalCodeCN(t *testing.T) {
	checkRules(t, "PostalCodeCN", is.PostalCodeCN,
		[]string{"310000", "100000", "518000"},
		[]string{"", "12345", "1234567", "12345a", "31000 "})
}
