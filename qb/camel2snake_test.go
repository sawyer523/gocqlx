package qb

import (
	"testing"
)

func TestCamelToSnakeASCII(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "camel_case"},
		{"camelCase", "camel_case"},
		{"Camel2Snake", "camel2_snake"},
		{"Camel_Snake", "camel_snake"},
		{"Camel__Snake", "camel__snake"},
		{"Camel2Snake3", "camel2_snake3"},
		{"Camel2_Snake3", "camel2_snake3"},
		{"Camel2__Snake3", "camel2__snake3"},
		{"Camel2Snake3Test", "camel2_snake3_test"},
		{"Camel2_Snake3_Test", "camel2_snake3_test"},
		{"Camel2__Snake3__Test", "camel2__snake3__test"},
		{"Camel2Snake3Test4", "camel2_snake3_test4"},
		{"Camel2_Snake3_Test4", "camel2_snake3_test4"},
		{"Camel2__Snake3__Test4", "camel2__snake3__test4"},
		{"Camel2Snake3Test4_", "camel2_snake3_test4_"},
		{"Camel2_Snake3_Test4_", "camel2_snake3_test4_"},
		{"Camel2__Snake3__Test4_", "camel2__snake3__test4_"},
		{"Camel2Snake3Test4__", "camel2_snake3_test4__"},
		{"Camel2_Snake3_Test4__", "camel2_snake3_test4__"},
		{"Camel2__Snake3__Test4__", "camel2__snake3__test4__"},
		{"Camel2Snake3Test4__5", "camel2_snake3_test4__5"},
		{"Camel2_Snake3_Test4__5", "camel2_snake3_test4__5"},
		{"Camel2__Snake3__Test4__5", "camel2__snake3__test4__5"},
		{"Camel2Snake3Test4__5_", "camel2_snake3_test4__5_"},
		{"Camel2_Snake3_Test4__5_", "camel2_snake3_test4__5_"},
		{"Camel2__Snake3__Test4__5_", "camel2__snake3__test4__5_"},
		{"\"Camel2Snake3Test4\"", "camel2_snake3_test4"},
	}

	for _, c := range cases {
		got := camelToSnakeASCII(c.input)
		if got != c.expected {
			t.Errorf("camelToSnakeASCII(%v) == %v, want %v", c.input, got, c.expected)
		}
	}
}
