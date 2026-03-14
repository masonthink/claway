package service

import "testing"

func TestSanitizeAvatarURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"empty string", "", ""},
		{"valid https", "https://pbs.twimg.com/profile_images/abc/photo.jpg", "https://pbs.twimg.com/profile_images/abc/photo.jpg"},
		{"http rejected", "http://example.com/avatar.jpg", ""},
		{"javascript rejected", "javascript:alert(1)", ""},
		{"data URI rejected", "data:image/png;base64,abc", ""},
		{"no host rejected", "https://", ""},
		{"ftp rejected", "ftp://example.com/file", ""},
		{"relative path rejected", "/images/avatar.jpg", ""},
		{"malformed URL", "://not-a-url", ""},
		{"https with query params", "https://example.com/avatar.jpg?size=400", "https://example.com/avatar.jpg?size=400"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeAvatarURL(tt.raw)
			if got != tt.want {
				t.Errorf("sanitizeAvatarURL(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
