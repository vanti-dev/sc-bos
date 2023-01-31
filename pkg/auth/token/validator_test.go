package token

import (
	"errors"
	"reflect"
	"testing"
)

func TestValidatorSet_modify(t *testing.T) {
	vs := ValidatorSet{}
	v1, v2, v3 := AlwaysValid(&Claims{}), NeverValid(errors.New("foo")), NeverValid(errors.New("bar"))
	v4 := AlwaysValid(&Claims{})

	vs.Append(v1)
	if !reflect.DeepEqual(vs, ValidatorSet{v1}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Append(v2)
	if !reflect.DeepEqual(vs, ValidatorSet{v1, v2}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Delete(v1)
	if !reflect.DeepEqual(vs, ValidatorSet{v2}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Append(v3)
	if !reflect.DeepEqual(vs, ValidatorSet{v2, v3}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Delete(v2)
	if !reflect.DeepEqual(vs, ValidatorSet{v3}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Delete(v4)
	if !reflect.DeepEqual(vs, ValidatorSet{v3}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Delete(v1)
	if !reflect.DeepEqual(vs, ValidatorSet{v3}) {
		t.Fatalf("ValidatorSet not expected")
	}
	vs.Delete(v3)
	if !reflect.DeepEqual(vs, ValidatorSet{}) {
		t.Fatalf("ValidatorSet not expected")
	}
}
