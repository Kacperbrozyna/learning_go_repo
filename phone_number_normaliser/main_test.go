package main

import "testing"

func TestNormalis(t *testing.T) {
	test_cases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, test_case := range test_cases {
		t.Run(test_case.input, func(t *testing.T) {
			actual := normalise(test_case.input)
			if actual != test_case.want {
				t.Errorf("got %s; want %s", actual, test_case.want)
			}
		})
	}
}
