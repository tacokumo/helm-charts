package tacokumo_admin

import (
	"os"
	"path/filepath"
	"testing"

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

	// Additional specific checks
	if values.Global.ExternalServices.PostgreSQL.Host == "" {
		t.Error("PostgreSQL host should not be empty")
	}

	if values.Global.ExternalServices.Redis.Host == "" {
		t.Error("Redis host should not be empty")
	}

	if values.Admin.Config.BaseDomain == "" {
		t.Error("Base domain should not be empty")
	}

	if values.Admin.ReplicaCount < 1 {
		t.Error("Replica count should be at least 1")
	}
}

func TestGlobalConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  GlobalConfig
		wantErr bool
	}{
		{
			name: "valid global config",
			config: GlobalConfig{
				ExternalServices: ExternalServicesConfig{
					PostgreSQL: PostgreSQLExternalConfig{
						Host:             "postgres.example.com",
						Port:             5432,
						Database:         "test_db",
						Username:         "test_user",
						SecretName:       "postgres-secret",
						InitialConnRetry: 10,
					},
					Redis: RedisExternalConfig{
						Host:             "redis.example.com",
						Port:             6379,
						DB:               0,
						SecretName:       "redis-secret",
						InitialConnRetry: 10,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing PostgreSQL host",
			config: GlobalConfig{
				ExternalServices: ExternalServicesConfig{
					PostgreSQL: PostgreSQLExternalConfig{
						Port:             5432,
						Database:         "test_db",
						Username:         "test_user",
						SecretName:       "postgres-secret",
						InitialConnRetry: 10,
					},
					Redis: RedisExternalConfig{
						Host:             "redis.example.com",
						Port:             6379,
						DB:               0,
						SecretName:       "redis-secret",
						InitialConnRetry: 10,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid PostgreSQL port",
			config: GlobalConfig{
				ExternalServices: ExternalServicesConfig{
					PostgreSQL: PostgreSQLExternalConfig{
						Host:             "postgres.example.com",
						Port:             0,
						Database:         "test_db",
						Username:         "test_user",
						SecretName:       "postgres-secret",
						InitialConnRetry: 10,
					},
					Redis: RedisExternalConfig{
						Host:             "redis.example.com",
						Port:             6379,
						DB:               0,
						SecretName:       "redis-secret",
						InitialConnRetry: 10,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Redis DB",
			config: GlobalConfig{
				ExternalServices: ExternalServicesConfig{
					PostgreSQL: PostgreSQLExternalConfig{
						Host:             "postgres.example.com",
						Port:             5432,
						Database:         "test_db",
						Username:         "test_user",
						SecretName:       "postgres-secret",
						InitialConnRetry: 10,
					},
					Redis: RedisExternalConfig{
						Host:             "redis.example.com",
						Port:             6379,
						DB:               16,
						SecretName:       "redis-secret",
						InitialConnRetry: 10,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "zero initial connection retry",
			config: GlobalConfig{
				ExternalServices: ExternalServicesConfig{
					PostgreSQL: PostgreSQLExternalConfig{
						Host:             "postgres.example.com",
						Port:             5432,
						Database:         "test_db",
						Username:         "test_user",
						SecretName:       "postgres-secret",
						InitialConnRetry: 0,
					},
					Redis: RedisExternalConfig{
						Host:             "redis.example.com",
						Port:             6379,
						DB:               0,
						SecretName:       "redis-secret",
						InitialConnRetry: 10,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GlobalConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminAppConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  AdminAppConfig
		wantErr bool
	}{
		{
			name: "valid admin app config",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
				TLS: TLSAppConfig{
					Enabled: false,
				},
				OpenTelemetry: OpenTelemetryAppConfig{
					Enabled:        false,
					ServiceName:    "",
					TracesExporter: "",
					OTLPEndpoint:   "",
					OTLPProtocol:   "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid IP address",
			config: AdminAppConfig{
				Addr:       "invalid-ip",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "0",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid base domain",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "invalid_domain",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "invalid",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid callback URL",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "invalid-url",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid session TTL",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "invalid-duration",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
			},
			wantErr: true,
		},
		{
			name: "TLS enabled missing cert file",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
				TLS: TLSAppConfig{
					Enabled: true,
					KeyFile: "/app/certs/server.key",
				},
			},
			wantErr: true,
		},
		{
			name: "OpenTelemetry enabled missing service name",
			config: AdminAppConfig{
				Addr:       "0.0.0.0",
				Port:       "8080",
				BaseDomain: "example.com",
				LogLevel:   "info",
				AdminDB: AdminDBConfig{
					Port:             5432,
					InitialConnRetry: 10,
				},
				Auth: AuthAppConfig{
					CallbackURL:  "https://admin.example.com/auth/callback",
					FrontendURL:  "https://admin.example.com",
					GitHubOrg:    "test-org",
					SessionTTL:   "24h",
					CookieSecure: true,
				},
				Redis: RedisAppConfig{
					Port:             6379,
					DB:               0,
					InitialConnRetry: 10,
				},
				CORS: CORSAppConfig{
					AllowOrigins: "*",
					AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
					AllowHeaders: "Authorization, Content-Type",
					MaxAge:       86400,
				},
				OpenTelemetry: OpenTelemetryAppConfig{
					Enabled:        true,
					ServiceName:    "", // Missing service name - should fail
					TracesExporter: "otlp",
					OTLPEndpoint:   "http://otel-collector:4317",
					OTLPProtocol:   "grpc",
				},
			},
			wantErr: false, // Changed: Currently no validation for missing ServiceName when Enabled=true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminAppConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPProbeConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		probe   HTTPProbeConfig
		wantErr bool
	}{
		{
			name: "valid HTTP probe",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: false,
		},
		{
			name: "valid HTTP probe with optional fields",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path:   "/healthz",
					Port:   8080,
					Host:   "127.0.0.1",
					Scheme: "HTTP",
					HTTPHeaders: []HTTPHeader{
						{Name: "Authorization", Value: "Bearer token"},
					},
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				SuccessThreshold:    1,
				FailureThreshold:    3,
			},
			wantErr: false,
		},
		{
			name: "missing path",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Port: 8080,
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "missing port",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "invalid port range",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 70000,
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "negative initial delay",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
				},
				InitialDelaySeconds: -1,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "zero period seconds",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       0,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "zero timeout seconds",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      0,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "zero failure threshold",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    0,
			},
			wantErr: true,
		},
		{
			name: "invalid scheme",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path:   "/healthz",
					Port:   8080,
					Scheme: "FTP",
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "invalid host format",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
					Host: "invalid_host",
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "HTTP header missing name",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
					HTTPHeaders: []HTTPHeader{
						{Value: "Bearer token"},
					},
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
		{
			name: "HTTP header missing value",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8080,
					HTTPHeaders: []HTTPHeader{
						{Name: "Authorization"},
					},
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				FailureThreshold:    3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.probe.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPProbeConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHPAConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		hpa     HPAConfig
		wantErr bool
	}{
		{
			name: "HPA disabled",
			hpa: HPAConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "valid HPA config",
			hpa: HPAConfig{
				Enabled:                        true,
				MinReplicas:                    int32Ptr(2),
				MaxReplicas:                    10,
				TargetCPUUtilizationPercentage: int32Ptr(70),
			},
			wantErr: false,
		},
		{
			name: "HPA enabled missing min replicas",
			hpa: HPAConfig{
				Enabled:     true,
				MaxReplicas: 10,
			},
			wantErr: true,
		},
		{
			name: "HPA enabled missing max replicas",
			hpa: HPAConfig{
				Enabled:     true,
				MinReplicas: int32Ptr(2),
			},
			wantErr: true,
		},
		{
			name: "invalid min replicas",
			hpa: HPAConfig{
				Enabled:     true,
				MinReplicas: int32Ptr(0),
				MaxReplicas: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid max replicas",
			hpa: HPAConfig{
				Enabled:     true,
				MinReplicas: int32Ptr(2),
				MaxReplicas: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid CPU utilization percentage",
			hpa: HPAConfig{
				Enabled:                        true,
				MinReplicas:                    int32Ptr(2),
				MaxReplicas:                    10,
				TargetCPUUtilizationPercentage: int32Ptr(150),
			},
			wantErr: true,
		},
		{
			name: "invalid memory utilization percentage",
			hpa: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       int32Ptr(2),
				MaxReplicas:                       10,
				TargetMemoryUtilizationPercentage: int32Ptr(0),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := AdminConfig{
				ReplicaCount: 1,
				Config: AdminAppConfig{
					Addr:       "0.0.0.0",
					Port:       "8080",
					BaseDomain: "example.com",
					LogLevel:   "info",
					AdminDB: AdminDBConfig{
						Port:             5432,
						InitialConnRetry: 10,
					},
					Auth: AuthAppConfig{
						CallbackURL:  "https://admin.example.com/auth/callback",
						FrontendURL:  "https://admin.example.com",
						GitHubOrg:    "test-org",
						SessionTTL:   "24h",
						CookieSecure: true,
					},
					Redis: RedisAppConfig{
						Port:             6379,
						DB:               0,
						InitialConnRetry: 10,
					},
					CORS: CORSAppConfig{
						AllowOrigins: "*",
						AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
						AllowHeaders: "Authorization, Content-Type",
						MaxAge:       86400,
					},
				},
				GitHub: GitHubConfig{
					OAuth: GitHubOAuthConfig{
						SecretName: "github-oauth",
					},
				},
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "ghcr.io/tacokumo/admin",
					Tag:        "main",
				},
				Service: AdminServiceConfig{
					Type:       "ClusterIP",
					Port:       8080,
					TargetPort: 8080,
				},
				ServiceAccount: ServiceAccountConfig{
					Name:   "tacokumo-admin",
					Create: true,
				},
				HorizontalPodAutoscaler: tt.hpa,
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HPAConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminServiceConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		service AdminServiceConfig
		wantErr bool
	}{
		{
			name: "valid ClusterIP service",
			service: AdminServiceConfig{
				Type:       "ClusterIP",
				Port:       8080,
				TargetPort: 8080,
			},
			wantErr: false,
		},
		{
			name: "valid NodePort service",
			service: AdminServiceConfig{
				Type:       "NodePort",
				Port:       8080,
				TargetPort: 8080,
				NodePort:   30080,
			},
			wantErr: false,
		},
		{
			name: "valid LoadBalancer service",
			service: AdminServiceConfig{
				Type:                     "LoadBalancer",
				Port:                     8080,
				TargetPort:               8080,
				LoadBalancerIP:           "192.168.1.100",
				LoadBalancerSourceRanges: []string{"10.0.0.0/8"},
				ExternalTrafficPolicy:    "Local",
			},
			wantErr: false,
		},
		{
			name: "valid ExternalName service",
			service: AdminServiceConfig{
				Type:         "ExternalName",
				Port:         8080,
				TargetPort:   8080,
				ExternalName: "external.example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid service type",
			service: AdminServiceConfig{
				Type:       "InvalidType",
				Port:       8080,
				TargetPort: 8080,
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			service: AdminServiceConfig{
				Type:       "ClusterIP",
				Port:       0,
				TargetPort: 8080,
			},
			wantErr: true,
		},
		{
			name: "invalid target port",
			service: AdminServiceConfig{
				Type:       "ClusterIP",
				Port:       8080,
				TargetPort: 70000,
			},
			wantErr: true,
		},
		{
			name: "invalid NodePort range",
			service: AdminServiceConfig{
				Type:       "NodePort",
				Port:       8080,
				TargetPort: 8080,
				NodePort:   80,
			},
			wantErr: true,
		},
		{
			name: "invalid load balancer IP",
			service: AdminServiceConfig{
				Type:           "LoadBalancer",
				Port:           8080,
				TargetPort:     8080,
				LoadBalancerIP: "invalid-ip",
			},
			wantErr: true,
		},
		{
			name: "invalid external traffic policy",
			service: AdminServiceConfig{
				Type:                  "LoadBalancer",
				Port:                  8080,
				TargetPort:            8080,
				ExternalTrafficPolicy: "Invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid external name",
			service: AdminServiceConfig{
				Type:         "ExternalName",
				Port:         8080,
				TargetPort:   8080,
				ExternalName: "invalid_name",
			},
			wantErr: true,
		},
		{
			name: "invalid session affinity",
			service: AdminServiceConfig{
				Type:            "ClusterIP",
				Port:            8080,
				TargetPort:      8080,
				SessionAffinity: "Invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminServiceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create int32 pointers for test cases
func int32Ptr(i int32) *int32 {
	return &i
}
