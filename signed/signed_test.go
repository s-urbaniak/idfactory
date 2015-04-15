package signed

import "testing"

var secret = []byte("something very secret")
var wrongSecret = []byte("wrong secret")

func TestValidation(t *testing.T) {
	s := New(secret)

	if !s.Validate(secret) {
		t.Fatal("signed uuid must validate itself successfully")
	}

	if s.Validate(wrongSecret) {
		t.Fatal("signed uuid must not validate with a wrong secret")
	}
}

func assertFailedParser(src string, t *testing.T) {
	_, err := Parse("")

	if err == nil {
		t.Fatal("an empty string must fail")
	}
}

func assertParser(src string, t *testing.T) *Signed {
	s, err := Parse(src)

	if err != nil {
		t.Fatal(err)
	}

	return s
}

func TestParserFailures(t *testing.T) {
	assertFailedParser("", t)
	assertFailedParser("abc", t)
	assertFailedParser("abc~abc", t)
	assertFailedParser("~", t)
	assertFailedParser("abc~", t)
	assertFailedParser("~abc", t)
}

func TestParser(t *testing.T) {
	s := assertParser("6ac34e8f-8780-4509-83f8-e15a345632e7~AjLn6tiDcBNgx8zMR0xH0Q3rjZxdBuv7Gpdz4MJN7rc", t)

	if !s.Validate(secret) {
		t.Fatal("signed uuid must validate itself successfully")
	}
}
