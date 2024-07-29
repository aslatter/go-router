package router

import "testing"

func Test_applyPrefixToPattern(t *testing.T) {
	type args struct {
		prefix  string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{},
			want: "",
		},
		{
			name: "nil prefix",
			args: args{
				pattern: "/some/path",
			},
			want: "/some/path",
		},
		{
			name: "nil prefix with method",
			args: args{
				pattern: "GET /some/path",
			},
			want: "GET /some/path",
		},
		{
			name: "basic prefix",
			args: args{
				prefix:  "/pre",
				pattern: "/some/path",
			},
			want: "/pre/some/path",
		},
		{
			name: "compound prefix",
			args: args{
				prefix:  "/pre/fix",
				pattern: "/some/path",
			},
			want: "/pre/fix/some/path",
		},
		{
			name: "slash prefix",
			args: args{
				prefix:  "/",
				pattern: "/some/path",
			},
			want: "/some/path",
		},
		{
			name: "slash compound prefix",
			args: args{
				prefix:  "/pre/fix",
				pattern: "/some/path",
			},
			want: "/pre/fix/some/path",
		},
		{
			name: "trailing slash compound prefix",
			args: args{
				prefix:  "/pre/fix/",
				pattern: "/some/path",
			},
			want: "/pre/fix/some/path",
		},
		{
			name: "basic prefix with method",
			args: args{
				prefix:  "/pre",
				pattern: "POST /some/path",
			},
			want: "POST /pre/some/path",
		},
		{
			name: "compound prefix with method",
			args: args{
				prefix:  "/pre/fix",
				pattern: "PUT /some/path",
			},
			want: "PUT /pre/fix/some/path",
		},
		{
			name: "trailing slash compound prefix with method",
			args: args{
				prefix:  "/pre/fix/",
				pattern: "DELETE /some/path",
			},
			want: "DELETE /pre/fix/some/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyPrefixToPattern(tt.args.prefix, tt.args.pattern); got != tt.want {
				t.Errorf("applyPrefixToPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
