package urlvalidation

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			url:     "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com/file",
			wantErr: true,
		},
		{
			name:    "localhost IP",
			url:     "http://127.0.0.1:8080",
			wantErr: true,
		},
		{
			name:    "private IPv4",
			url:     "http://10.0.0.1",
			wantErr: true,
		},
		{
			name:    "private IPv4",
			url:     "http://192.168.1.10",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.url)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"Validate(%q) error = %v, wantErr = %v",
					tt.url,
					err,
					tt.wantErr,
				)
			}
		})
	}
}
