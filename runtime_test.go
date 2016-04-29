package template

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	type testElement struct {
		t string
		r map[string]string
	}

	tests := []testElement{
		testElement{`asdf`, map[string]string{"0": "asdf"}},             // 0
		testElement{`<? ?>asdf`, map[string]string{"0": "asdf"}},        // 1
		testElement{`<? 1  ?>asdf`, map[string]string{"1": "asdf"}},     // 2
		testElement{`<?  abc ?>asdf`, map[string]string{"abc": "asdf"}}, // 3
		testElement{`asdf<? 	?>`, map[string]string{"0": "asdf", "1": ""}}, // 4

		testElement{``, map[string]string{"0": ""}}, // 5

		testElement{`<? 213 ?>s1<? ?>s2`, map[string]string{"213": "s1", "214": "s2"}},                     // 6
		testElement{`<? 213 ?>s1<? a ?>s2<? ?>s3`, map[string]string{"213": "s1", "a": "s2", "214": "s3"}}, // 7
	}

	for i, test := range tests {
		if r := ParseString(test.t); !reflect.DeepEqual(r, test.r) {
			t.Errorf("%v: Expected: %v Got: %v", i, test.r, r)
		}
	}
}
