package generator

import "testing"

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "simple",
			in:   "foo_bar",
			want: "fooBar",
		},
		{
			name: "double underscore",
			in:   "foo__bar",
			want: "fooBar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCase(tt.in); got != tt.want {
				t.Errorf("camelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpperCamelCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "simple",
			in:   "foo_bar",
			want: "FooBar",
		},
		{
			name: "double underscore",
			in:   "foo__bar",
			want: "FooBar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upperCamelCase(tt.in); got != tt.want {
				t.Errorf("upperCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
