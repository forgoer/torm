package utils

import "testing"

func TestStudlyCase(t *testing.T) {
	a := "xx_yy"
	b := StudlyCase(a)
	expected := "XxYy"
	if b != expected {
		t.Errorf("Expected the `StudlyCase` of %s to be %s but instead got %s !", a, expected, b)
	}
}

func TestSnakeCase(t *testing.T) {
	a := "XxYy"
	b := SnakeCase(a)
	expected := "xx_yy"
	if b != expected {
		t.Errorf("Expected the `SnakeCase` of %s to be %s but instead got %s !", a, expected, b)
	}

	a = "XxYY"
	b = SnakeCase(a)
	expected = "xx_yy"
	if b != expected {
		t.Errorf("Expected the `SnakeCase` of %s to be %s but instead got %s !", a, expected, b)
	}
}