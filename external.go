package helmcharts

import (
	"github.com/go-playground/validator/v10"
)

// PostgreSQLConfig represents PostgreSQL database configuration
type PostgreSQLConfig struct {
	Host     string `yaml:"host" validate:"required,fqdn|ip"`
	Port     int    `yaml:"port" validate:"required,min=1,max=65535"`
	Database string `yaml:"database" validate:"required"`
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	SSLMode  string `yaml:"sslMode,omitempty" validate:"omitempty,oneof=disable allow prefer require verify-ca verify-full"`

	// Connection pool settings
	MaxOpenConns    int    `yaml:"maxOpenConns,omitempty" validate:"omitempty,min=1"`
	MaxIdleConns    int    `yaml:"maxIdleConns,omitempty" validate:"omitempty,min=1"`
	ConnMaxLifetime string `yaml:"connMaxLifetime,omitempty" validate:"omitempty,duration"`
	ConnMaxIdleTime string `yaml:"connMaxIdleTime,omitempty" validate:"omitempty,duration"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host" validate:"required,fqdn|ip"`
	Port     int    `yaml:"port" validate:"required,min=1,max=65535"`
	Password string `yaml:"password,omitempty"`
	DB       int    `yaml:"db" validate:"min=0,max=15"`

	// Connection settings
	MaxRetries      int    `yaml:"maxRetries,omitempty" validate:"omitempty,min=0"`
	MinRetryBackoff string `yaml:"minRetryBackoff,omitempty" validate:"omitempty,duration"`
	MaxRetryBackoff string `yaml:"maxRetryBackoff,omitempty" validate:"omitempty,duration"`
	DialTimeout     string `yaml:"dialTimeout,omitempty" validate:"omitempty,duration"`
	ReadTimeout     string `yaml:"readTimeout,omitempty" validate:"omitempty,duration"`
	WriteTimeout    string `yaml:"writeTimeout,omitempty" validate:"omitempty,duration"`

	// Pool settings
	PoolSize           int    `yaml:"poolSize,omitempty" validate:"omitempty,min=1"`
	MinIdleConns       int    `yaml:"minIdleConns,omitempty" validate:"omitempty,min=0"`
	MaxConnAge         string `yaml:"maxConnAge,omitempty" validate:"omitempty,duration"`
	PoolTimeout        string `yaml:"poolTimeout,omitempty" validate:"omitempty,duration"`
	IdleTimeout        string `yaml:"idleTimeout,omitempty" validate:"omitempty,duration"`
	IdleCheckFrequency string `yaml:"idleCheckFrequency,omitempty" validate:"omitempty,duration"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Enabled bool `yaml:"enabled"`

	// OAuth settings
	OAuthClientID     string   `yaml:"oauthClientId" validate:"required_if=Enabled true"`
	OAuthClientSecret string   `yaml:"oauthClientSecret" validate:"required_if=Enabled true"`
	OAuthAuthURL      string   `yaml:"oauthAuthUrl" validate:"required_if=Enabled true,url"`
	OAuthTokenURL     string   `yaml:"oauthTokenUrl" validate:"required_if=Enabled true,url"`
	OAuthCallbackURL  string   `yaml:"oauthCallbackUrl" validate:"required_if=Enabled true,url"`
	OAuthScopes       []string `yaml:"oauthScopes,omitempty"`

	// Session settings
	SessionSecret   string `yaml:"sessionSecret" validate:"required_if=Enabled true,min=32"`
	SessionName     string `yaml:"sessionName" validate:"required_if=Enabled true"`
	SessionTTL      string `yaml:"sessionTtl" validate:"required_if=Enabled true,duration"`
	SessionSecure   bool   `yaml:"sessionSecure"`
	SessionSameSite string `yaml:"sessionSameSite" validate:"omitempty,oneof=Strict Lax None"`
	SessionDomain   string `yaml:"sessionDomain,omitempty" validate:"omitempty,fqdn"`

	// JWT settings
	JWTSecret            string `yaml:"jwtSecret,omitempty" validate:"omitempty,min=32"`
	JWTExpiration        string `yaml:"jwtExpiration,omitempty" validate:"omitempty,duration"`
	JWTRefreshExpiration string `yaml:"jwtRefreshExpiration,omitempty" validate:"omitempty,duration"`
}

// CORSConfig represents CORS policy configuration
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowedOrigins   []string `yaml:"allowedOrigins" validate:"required_if=Enabled true,dive,url"`
	AllowedMethods   []string `yaml:"allowedMethods" validate:"required_if=Enabled true,dive,oneof=GET POST PUT DELETE PATCH HEAD OPTIONS"`
	AllowedHeaders   []string `yaml:"allowedHeaders,omitempty"`
	ExposedHeaders   []string `yaml:"exposedHeaders,omitempty"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge,omitempty" validate:"omitempty,min=0"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CertFile   string `yaml:"certFile" validate:"required_if=Enabled true,filepath"`
	KeyFile    string `yaml:"keyFile" validate:"required_if=Enabled true,filepath"`
	CAFile     string `yaml:"caFile,omitempty" validate:"omitempty,filepath"`
	MinVersion string `yaml:"minVersion,omitempty" validate:"omitempty,oneof=1.0 1.1 1.2 1.3"`
	MaxVersion string `yaml:"maxVersion,omitempty" validate:"omitempty,oneof=1.0 1.1 1.2 1.3"`
}

// OpenTelemetryConfig represents OpenTelemetry configuration
type OpenTelemetryConfig struct {
	Enabled        bool   `yaml:"enabled"`
	ServiceName    string `yaml:"serviceName" validate:"required_if=Enabled true"`
	ServiceVersion string `yaml:"serviceVersion,omitempty"`

	// Tracing
	TracingEnabled  bool    `yaml:"tracingEnabled"`
	TracingEndpoint string  `yaml:"tracingEndpoint" validate:"required_if=TracingEnabled true,url"`
	TracingSampling float64 `yaml:"tracingSampling" validate:"min=0,max=1"`

	// Metrics
	MetricsEnabled  bool   `yaml:"metricsEnabled"`
	MetricsEndpoint string `yaml:"metricsEndpoint" validate:"required_if=MetricsEnabled true,url"`
	MetricsInterval string `yaml:"metricsInterval" validate:"omitempty,duration"`

	// Logging
	LoggingEnabled bool   `yaml:"loggingEnabled"`
	LogLevel       string `yaml:"logLevel" validate:"omitempty,oneof=debug info warn error"`
	LogFormat      string `yaml:"logFormat" validate:"omitempty,oneof=json text"`

	// Resource attributes
	ResourceAttributes map[string]string `yaml:"resourceAttributes,omitempty"`

	// Headers for authentication
	Headers map[string]string `yaml:"headers,omitempty"`
}

// ExternalServiceConfig represents configuration for external service dependencies
type ExternalServiceConfig struct {
	// Database
	Database PostgreSQLConfig `yaml:"database" validate:"required"`

	// Cache
	Redis RedisConfig `yaml:"redis" validate:"required"`

	// Authentication
	Auth AuthConfig `yaml:"auth"`

	// CORS policy
	CORS CORSConfig `yaml:"cors"`

	// TLS configuration
	TLS TLSConfig `yaml:"tls"`

	// Observability
	OpenTelemetry OpenTelemetryConfig `yaml:"openTelemetry"`
}

// Validate validates the external service configurations
func (p *PostgreSQLConfig) Validate() error {
	return validator.New().Struct(p)
}

func (r *RedisConfig) Validate() error {
	return validator.New().Struct(r)
}

func (a *AuthConfig) Validate() error {
	return validator.New().Struct(a)
}

func (c *CORSConfig) Validate() error {
	return validator.New().Struct(c)
}

func (t *TLSConfig) Validate() error {
	return validator.New().Struct(t)
}

func (o *OpenTelemetryConfig) Validate() error {
	return validator.New().Struct(o)
}

func (e *ExternalServiceConfig) Validate() error {
	return validator.New().Struct(e)
}
