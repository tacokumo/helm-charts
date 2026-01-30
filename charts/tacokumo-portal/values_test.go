package tacokumo_portal

import (
	"os"
	"path/filepath"
	"testing"

	admin "github.com/tacokumo/helm-charts/charts/tacokumo-admin"
	"gopkg.in/yaml.v3"
)

func TestLoadAndValidateValuesYAML(t *testing.T) {
	// Get the path to the values.yaml file
	valuesPath := filepath.Join("values.yaml")

	// Read the values.yaml file
	data, err := os.ReadFile(valuesPath)
	if err != nil {
		t.Fatalf("Failed to read values.yaml: %v", err)
	}

	// Parse the YAML
	var values Values
	err = yaml.Unmarshal(data, &values)
	if err != nil {
		t.Fatalf("Failed to unmarshal values.yaml: %v", err)
	}

	// Validate the configuration
	err = values.Validate()
	if err != nil {
		t.Errorf("values.yaml validation failed: %v", err)
	}

	// Additional specific checks for portal
	if values.Global.ExternalServices.PostgreSQL.Host == "" {
		t.Error("PostgreSQL host should not be empty")
	}

	if values.Global.ExternalServices.Redis.Host == "" {
		t.Error("Redis host should not be empty")
	}

	if values.Portal.Config.BaseDomain == "" {
		t.Error("Base domain should not be empty")
	}

	if values.Portal.ReplicaCount < 1 {
		t.Error("Replica count should be at least 1")
	}

	// Check portal-specific database naming
	if values.Portal.Config.PortalDB.InitialConnRetry < 1 {
		t.Error("Portal DB initial connection retry should be at least 1")
	}
}

func TestPortalConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PortalConfig
		wantErr bool
	}{
		{
			name: "valid portal config",
			config: PortalConfig{
				ReplicaCount: 1,
				Config: PortalAppConfig{
					Addr:       "0.0.0.0",
					Port:       "8080",
					BaseDomain: "example.com",
					LogLevel:   "info",
					PortalDB: PortalDBConfig{
						Port:             5432,
						InitialConnRetry: 10,
					},
					Auth: admin.AuthAppConfig{
						CallbackURL:  "https://portal.example.com/auth/callback",
						FrontendURL:  "https://portal.example.com",
						GitHubOrg:    "test-org",
						SessionTTL:   "24h",
						CookieSecure: true,
					},
					Redis: admin.RedisAppConfig{
						Port:             6379,
						DB:               0,
						InitialConnRetry: 10,
					},
					CORS: admin.CORSAppConfig{
						AllowOrigins: "*",
						AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
						AllowHeaders: "Authorization, Content-Type",
						MaxAge:       86400,
					},
				},
				GitHub: admin.GitHubConfig{
					OAuth: admin.GitHubOAuthConfig{
						SecretName: "github-oauth",
					},
				},
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "ghcr.io/tacokumo/portal",
					Tag:        "main",
				},
				Service: PortalServiceConfig{
					Type:       "ClusterIP",
					Port:       8080,
					TargetPort: 8080,
				},
				ServiceAccount: admin.ServiceAccountConfig{
					Name:   "tacokumo-portal",
					Create: true,
				},
			},
			wantErr: false,
		},
		{
			name: "zero replica count",
			config: PortalConfig{
				ReplicaCount: 0,
				Config: PortalAppConfig{
					Addr:       "0.0.0.0",
					Port:       "8080",
					BaseDomain: "example.com",
					LogLevel:   "info",
					PortalDB: PortalDBConfig{
						Port:             5432,
						InitialConnRetry: 10,
					},
					Auth: admin.AuthAppConfig{
						CallbackURL:  "https://portal.example.com/auth/callback",
						FrontendURL:  "https://portal.example.com",
						GitHubOrg:    "test-org",
						SessionTTL:   "24h",
						CookieSecure: true,
					},
					Redis: admin.RedisAppConfig{
						Port:             6379,
						DB:               0,
						InitialConnRetry: 10,
					},
					CORS: admin.CORSAppConfig{
						AllowOrigins: "*",
						AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
						AllowHeaders: "Authorization, Content-Type",
						MaxAge:       86400,
					},
				},
				GitHub: admin.GitHubConfig{
					OAuth: admin.GitHubOAuthConfig{
						SecretName: "github-oauth",
					},
				},
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "ghcr.io/tacokumo/portal",
					Tag:        "main",
				},
				Service: PortalServiceConfig{
					Type:       "ClusterIP",
					Port:       8080,
					TargetPort: 8080,
				},
				ServiceAccount: admin.ServiceAccountConfig{
					Name:   "tacokumo-portal",
					Create: true,
				},
			},
			wantErr: true, // ReplicaCount: 0 violates validate:"min=1"
		},
		{
			name: "missing GitHub OAuth secret name",
			config: PortalConfig{
				ReplicaCount: 1,
				Config: PortalAppConfig{
					Addr:       "0.0.0.0",
					Port:       "8080",
					BaseDomain: "example.com",
					LogLevel:   "info",
					PortalDB: PortalDBConfig{
						Port:             5432,
						InitialConnRetry: 10,
					},
					Auth: admin.AuthAppConfig{
						CallbackURL:  "https://portal.example.com/auth/callback",
						FrontendURL:  "https://portal.example.com",
						GitHubOrg:    "test-org",
						SessionTTL:   "24h",
						CookieSecure: true,
					},
					Redis: admin.RedisAppConfig{
						Port:             6379,
						DB:               0,
						InitialConnRetry: 10,
					},
					CORS: admin.CORSAppConfig{
						AllowOrigins: "*",
						AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
						AllowHeaders: "Authorization, Content-Type",
						MaxAge:       86400,
					},
				},
				GitHub: admin.GitHubConfig{
					OAuth: admin.GitHubOAuthConfig{},
				},
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "ghcr.io/tacokumo/portal",
					Tag:        "main",
				},
				Service: PortalServiceConfig{
					Type:       "ClusterIP",
					Port:       8080,
					TargetPort: 8080,
				},
				ServiceAccount: admin.ServiceAccountConfig{
					Name:   "tacokumo-portal",
					Create: true,
				},
			},
			wantErr: true, // SecretName is required field
		},
		{
			name: "missing service account name",
			config: PortalConfig{
				ReplicaCount: 1,
				Config: PortalAppConfig{
					Addr:       "0.0.0.0",
					Port:       "8080",
					BaseDomain: "example.com",
					LogLevel:   "info",
					PortalDB: PortalDBConfig{
						Port:             5432,
						InitialConnRetry: 10,
					},
					Auth: admin.AuthAppConfig{
						CallbackURL:  "https://portal.example.com/auth/callback",
						FrontendURL:  "https://portal.example.com",
						GitHubOrg:    "test-org",
						SessionTTL:   "24h",
						CookieSecure: true,
					},
					Redis: admin.RedisAppConfig{
						Port:             6379,
						DB:               0,
						InitialConnRetry: 10,
					},
					CORS: admin.CORSAppConfig{
						AllowOrigins: "*",
						AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
						AllowHeaders: "Authorization, Content-Type",
						MaxAge:       86400,
					},
				},
				GitHub: admin.GitHubConfig{
					OAuth: admin.GitHubOAuthConfig{
						SecretName: "github-oauth",
					},
				},
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "ghcr.io/tacokumo/portal",
					Tag:        "main",
				},
				Service: PortalServiceConfig{
					Type:       "ClusterIP",
					Port:       8080,
					TargetPort: 8080,
				},
				ServiceAccount: admin.ServiceAccountConfig{
					Create: true,
				},
			},
			wantErr: true, // Name field is required
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortalConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPortalAppConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PortalAppConfig
		wantErr bool
	}{
		{
			name: "valid portal app config",
			config: PortalAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				PortalDB: PortalDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: admin.AuthAppConfig{
					CallbackURL:  "https://portal.example.com/auth/callback",
					FrontendURL:  "https://portal.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: admin.RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: admin.CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
				TLS: admin.TLSAppConfig{
					Enabled: false,
				},
				OpenTelemetry: PortalOpenTelemetryAppConfig{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid IP address",
			config: PortalAppConfig{
				Addr:       "invalid-ip",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				PortalDB: PortalDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: admin.AuthAppConfig{
					CallbackURL:  "https://portal.example.com/auth/callback",
					FrontendURL:  "https://portal.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: admin.RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: admin.CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true, // Addr: "invalid-ip" violates validate:"required,ip"
		},
		{
			name: "invalid base domain",
			config: PortalAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "invalid_domain",
				LogLevel:   "info",
				PortalDB: PortalDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: admin.AuthAppConfig{
					CallbackURL:  "https://portal.example.com/auth/callback",
					FrontendURL:  "https://portal.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: admin.RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: admin.CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true, // BaseDomain: "invalid_domain" violates validate:"required,fqdn"
		},
		{
			name: "zero portal DB initial connection retry",
			config: PortalAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				PortalDB: PortalDBConfig{
					Port:             5432,
					InitialConnRetry: 0,
				},
				Auth: admin.AuthAppConfig{
					CallbackURL:  "https://portal.example.com/auth/callback",
					FrontendURL:  "https://portal.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: admin.RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: admin.CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true, // InitialConnRetry: 0 violates validate:"min=1"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortalAppConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPortalDBConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PortalDBConfig
		wantErr bool
	}{
		{
			name: "valid portal DB config",
			config: PortalDBConfig{
				Port:             5432,
				InitialConnRetry: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: PortalDBConfig{
				Port:             0,
				InitialConnRetry: 10,
			},
			wantErr: true, // Port: 0 violates validate:"min=1,max=65535"
		},
		{
			name: "invalid port range (too high)",
			config: PortalDBConfig{
				Port:             70000,
				InitialConnRetry: 10,
			},
			wantErr: true, // Port: 70000 violates validate:"min=1,max=65535"
		},
		{
			name: "zero initial connection retry",
			config: PortalDBConfig{
				Port:             5432,
				InitialConnRetry: 0,
			},
			wantErr: true, // InitialConnRetry: 0 violates validate:"min=1"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortalDBConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPortalOpenTelemetryAppConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PortalOpenTelemetryAppConfig
		wantErr bool
	}{
		{
			name: "OpenTelemetry disabled",
			config: PortalOpenTelemetryAppConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "valid OpenTelemetry config",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				ServiceName:    "tacokumo-portal",
				TracesExporter: "otlp",
				OTLPEndpoint:   "http://otel-collector:4317",
				OTLPProtocol:   "grpc",
			},
			wantErr: false,
		},
		{
			name: "OpenTelemetry enabled missing service name",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				TracesExporter: "otlp",
				OTLPEndpoint:   "http://otel-collector:4317",
				OTLPProtocol:   "grpc",
			},
			wantErr: false, // ServiceName is optional field, no validation error expected
		},
		{
			name: "OpenTelemetry enabled missing traces exporter",
			config: PortalOpenTelemetryAppConfig{
				Enabled:      true,
				ServiceName:  "tacokumo-portal",
				OTLPEndpoint: "http://otel-collector:4317",
				OTLPProtocol: "grpc",
			},
			wantErr: false, // TracesExporter is optional field, no validation error expected
		},
		{
			name: "OpenTelemetry enabled invalid traces exporter",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				ServiceName:    "tacokumo-portal",
				TracesExporter: "invalid",
				OTLPEndpoint:   "http://otel-collector:4317",
				OTLPProtocol:   "grpc",
			},
			wantErr: true, // TracesExporter: "invalid" violates validate:"omitempty,oneof=otlp jaeger zipkin console"
		},
		{
			name: "OpenTelemetry enabled missing OTLP endpoint",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				ServiceName:    "tacokumo-portal",
				TracesExporter: "otlp",
				OTLPProtocol:   "grpc",
			},
			wantErr: false, // OTLPEndpoint is optional field, no validation error expected
		},
		{
			name: "OpenTelemetry enabled invalid OTLP endpoint",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				ServiceName:    "tacokumo-portal",
				TracesExporter: "otlp",
				OTLPEndpoint:   "invalid-url",
				OTLPProtocol:   "grpc",
			},
			wantErr: true, // OTLPEndpoint: "invalid-url" violates validate:"omitempty,url"
		},
		{
			name: "OpenTelemetry enabled missing OTLP protocol",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				ServiceName:    "tacokumo-portal",
				TracesExporter: "otlp",
				OTLPEndpoint:   "http://otel-collector:4317",
			},
			wantErr: false, // OTLPProtocol is optional field, no validation error expected
		},
		{
			name: "OpenTelemetry enabled invalid OTLP protocol",
			config: PortalOpenTelemetryAppConfig{
				Enabled:        true,
				ServiceName:    "tacokumo-portal",
				TracesExporter: "otlp",
				OTLPEndpoint:   "http://otel-collector:4317",
				OTLPProtocol:   "invalid",
			},
			wantErr: true, // OTLPProtocol: "invalid" violates validate:"omitempty,oneof=grpc http"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortalOpenTelemetryAppConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPortalServiceConfigValidation(t *testing.T) {
	// PortalServiceConfig is identical to AdminServiceConfig, so we just test basic functionality
	tests := []struct {
		name    string
		service PortalServiceConfig
		wantErr bool
	}{
		{
			name: "valid ClusterIP service",
			service: PortalServiceConfig{
				Type:       "ClusterIP",
				Port:       8080,
				TargetPort: 8080,
			},
			wantErr: false,
		},
		{
			name: "invalid service type",
			service: PortalServiceConfig{
				Type:       "InvalidType",
				Port:       8080,
				TargetPort: 8080,
			},
			wantErr: true, // Type: "InvalidType" is not a valid service type
		},
		{
			name: "invalid port",
			service: PortalServiceConfig{
				Type:       "ClusterIP",
				Port:       0,
				TargetPort: 8080,
			},
			wantErr: true, // Port: 0 violates validation rules
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortalServiceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
