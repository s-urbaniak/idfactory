package signed

import (
	"testing"
)

var secret = []byte("something very secret")
var wrongSecret = []byte("wrong secret")

func TestValidation(t *testing.T) {
	s := New(secret)

	if !s.Validate(secret) {
		t.Error("signed uuid must validate itself successfully")
	}

	if s.Validate(wrongSecret) {
		t.Error("signed uuid must not validate with a wrong secret")
	}
}

func assertFailedParser(src string, t *testing.T) {
	_, err := Parse("")

	if err == nil {
		t.Error("an empty string must fail")
	}
}

func assertParser(src string, t *testing.T) *Signed {
	s, err := Parse(src)

	if err != nil {
		t.Error(err)
	}

	return s
}

func TestParserFailures(t *testing.T) {
	assertFailedParser("", t)
	assertFailedParser("abc", t)
	assertFailedParser("abc:abc", t)
	assertFailedParser(":", t)
	assertFailedParser("abc:", t)
	assertFailedParser(":abc", t)
}

func TestParser(t *testing.T) {
	s := assertParser("3be0d941-0996-49bd-8de1-82a7bdf6282a:E7uy78y1pia2CYmsza/w6OIRVsd5yYgjuykDDcx4RzQ=", t)

	if !s.Validate(secret) {
		t.Error("signed uuid must validate itself successfully")
	}
}
