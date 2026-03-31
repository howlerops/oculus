package format

import (
	"testing"
	"time"
)

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{500, "500"},
		{1000, "1.0k"},
		{1500, "1.5k"},
		{1000000, "1.0M"},
	}
	for _, tt := range tests {
		if got := FormatTokens(tt.n); got != tt.want {
			t.Errorf("FormatTokens(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{500 * time.Millisecond, "500ms"},
		{2 * time.Second, "2.0s"},
		{90 * time.Second, "1m30s"},
	}
	for _, tt := range tests {
		if got := FormatDuration(tt.d); got != tt.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestTruncateString(t *testing.T) {
	if got := TruncateString("hello world", 8); got != "hello..." {
		t.Errorf("TruncateString = %q, want 'hello...'", got)
	}
	if got := TruncateString("hi", 10); got != "hi" {
		t.Errorf("TruncateString short = %q, want 'hi'", got)
	}
}
