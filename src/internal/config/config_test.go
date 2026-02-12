package config

import "testing"

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVal   EnvironmentApp
		wantProd bool
		wantDev  bool
	}{
		{
			name:     "Production Environment",
			envVal:   EnvironmentProduction,
			wantProd: true,
			wantDev:  false,
		},
		{
			name:     "Development Environment",
			envVal:   EnvironmentDevelopment,
			wantProd: false,
			wantDev:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envInstance = appEnv{AppEnv: tt.envVal}

			LoadConfig()

			got := Config()

			if got.IsProduction != tt.wantProd {
				t.Errorf("IsProduction = %v, want %v", got.IsProduction, tt.wantProd)
			}

			if got.IsDevelopment != tt.wantDev {
				t.Errorf("IsDevelopment = %v, want %v", got.IsDevelopment, tt.wantDev)
			}
		})
	}
}
