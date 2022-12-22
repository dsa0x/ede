package evaluator

import (
	"fmt"
	"testing"
)

func TestEval_InfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1.0 == 1.00", true},
		{"1.0 != 1.0", false},
		{"1.1 == 1.2", false},
		{"1.1 != 1.2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"true || false", true},
		{"false || true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			evaluated := testEval(tt.input)
			if !testBooleanObject(t, evaluated, tt.expected) {
				t.FailNow()
			}
		})
	}
}
