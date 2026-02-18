package config

import "testing"

func TestEnvironmentApp_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    EnvironmentApp
		wantErr bool
	}{
		{
			name:    "Valid Production",
			input:   "production",
			want:    EnvironmentProduction,
			wantErr: false,
		},
		{
			name:    "Valid Development",
			input:   "development",
			want:    EnvironmentDevelopment,
			wantErr: false,
		},
		{
			name:    "Invalid Typo",
			input:   "dev",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid Empty",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e EnvironmentApp

			err := e.UnmarshalText([]byte(tt.input))

			if (err != nil) != tt.wantErr {
				t.Fatalf("got %v, want %v", err, tt.wantErr)
			}

			if e != tt.want {
				t.Errorf("got %v, want %v", e, tt.want)
			}
		})
	}
}
