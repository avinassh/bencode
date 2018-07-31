package bencode

import (
	"testing"
)

func TestDecoder(t *testing.T) {

	type testCase struct {
		input    string
		expected string
	}

	tests := []testCase{
		// string tests
		{input: "4:spam", expected: `"spam"`},
		{input: "0:", expected: ""},
		{input: "3:gg\"", expected: `"gg\""`},

		// integers
		{input: "i3e", expected: `3`},
		{input: "i-3e", expected: `-3`},
		{input: "i0e", expected: `0`},

		// list
		{input: "l1:se", expected: `["s"]`},
		{input: "l1:s2:gge", expected: `["s","gg"]`},
		{input: "l4:spam4:eggse", expected: `["spam","eggs"]`},
		{input: "le", expected: `[]`},
		{input: "li3ee", expected: `[3]`},
		{input: "li3ei-3ee", expected: `[3,-3]`},
		{input: "li3ei-3e4:spam4:eggse", expected: `[3,-3,"spam","eggs"]`},

		//map
		{input: "d3:cow3:moo4:spam4:eggse", expected: `{"cow":"moo","spam":"eggs"}`},
		{input: "d4:spaml1:a1:bee", expected: `{"spam":["a","b"]}`},
		{input: "d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee",
			expected: `{"publisher":"bob","publisher-webpage":"www.example.com","publisher.location":"home"}`},
		{input: "d4:spaml1:a1:be3:numi3ee", expected: `{"num":3,"spam":["a","b"]}`},
	}

	for _, test := range tests {

		enc := NewBencoder(test.input)
		value := enc.Parse()
		got := string(value.JsonValue)

		if got != test.expected {
			t.Errorf("Expected: %s, got %s", test.expected, got)
		}
	}
}
