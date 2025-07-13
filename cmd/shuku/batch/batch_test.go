package batch

import (
	"testing"
)

func TestBoolToString(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{
			name:     "true の場合",
			input:    true,
			expected: "有効",
		},
		{
			name:     "false の場合",
			input:    false,
			expected: "無効",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := boolToString(tt.input)
			if result != tt.expected {
				t.Errorf("boolToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "バイト単位",
			size:     512,
			expected: "512 B",
		},
		{
			name:     "キロバイト単位",
			size:     1536, // 1.5 KB
			expected: "1.5 KB",
		},
		{
			name:     "メガバイト単位",
			size:     2097152, // 2.0 MB
			expected: "2.0 MB",
		},
		{
			name:     "ギガバイト単位",
			size:     3221225472, // 3.0 GB
			expected: "3.0 GB",
		},
		{
			name:     "ゼロバイト",
			size:     0,
			expected: "0 B",
		},
		{
			name:     "1バイト",
			size:     1,
			expected: "1 B",
		},
		{
			name:     "1キロバイト",
			size:     1024,
			expected: "1.0 KB",
		},
		{
			name:     "1メガバイト",
			size:     1048576,
			expected: "1.0 MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileSize(tt.size)
			if result != tt.expected {
				t.Errorf("formatFileSize() = %v, want %v", result, tt.expected)
			}
		})
	}
}