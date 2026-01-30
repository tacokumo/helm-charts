package tacokumo_portal

import (
	helmcharts "github.com/tacokumo/helm-charts"
	admin "github.com/tacokumo/helm-charts/charts/tacokumo-admin"
)

// Values represents the root configuration for tacokumo-portal Helm chart
// Similar structure to admin but with portal-specific naming
type Values struct {
	Global  GlobalConfig       `yaml:"global" validate:"required"`
	Portal  PortalConfig       `yaml:"portal" validate:"required"`
	Ingress helmcharts.Ingress `yaml:"ingress"`
}

// GlobalConfig represents global configuration shared across components
// Reuses admin structure with appropriate naming
type GlobalConfig = admin.GlobalConfig

// PortalConfig represents the portal application configuration
// Similar to AdminConfig but with portal-specific field names
type PortalConfig struct {
	ReplicaCount int `yaml:"replicaCount" validate:"min=1"`

	// Application configuration
	Config PortalAppConfig `yaml:"config" validate:"required"`

	// GitHub OAuth configuration
	GitHub admin.GitHubConfig `yaml:"github" validate:"required"`

	// Container image configuration
	Image helmcharts.Image `yaml:"image" validate:"required"`

	// Service configuration
	Service PortalServiceConfig `yaml:"service" validate:"required"`

	// Service account configuration
	ServiceAccount admin.ServiceAccountConfig `yaml:"serviceAccount" validate:"required"`

	// Pod configuration
	TerminationGracePeriodSeconds int64                  `yaml:"terminationGracePeriodSeconds" validate:"min=0"`
	Affinity                      *helmcharts.Affinity   `yaml:"affinity,omitempty"`
	NodeSelector                  map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations                   helmcharts.Tolerations `yaml:"tolerations,omitempty"`
	Labels                        map[string]string      `yaml:"labels,omitempty"`
	Annotations                   map[string]string      `yaml:"annotations,omitempty"`
	PodAnnotations                map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels                     map[string]string      `yaml:"podLabels,omitempty"`

	// Security contexts
	SecurityContext    *helmcharts.SecurityContext `yaml:"securityContext,omitempty"`
	PodSecurityContext *admin.PodSecurityContext   `yaml:"podSecurityContext,omitempty"`

	// Health checks
	LivenessProbe  *admin.HTTPProbeConfig `yaml:"livenessProbe,omitempty"`
	ReadinessProbe *admin.HTTPProbeConfig `yaml:"readinessProbe,omitempty"`

	// Resource limits and requests
	Resources helmcharts.Resources `yaml:"resources"`

	// Additional configurations
	ImagePullSecrets        []admin.ImagePullSecret        `yaml:"imagePullSecrets,omitempty" validate:"dive"`
	PodDisruptionBudget     helmcharts.PodDisruptionBudget `yaml:"podDisruptionBudget"`
	HorizontalPodAutoscaler admin.HPAConfig                `yaml:"horizontalPodAutoscaler"`

	// Environment variables and volumes
	Env          []admin.EnvVar        `yaml:"env,omitempty" validate:"dive"`
	EnvFrom      []admin.EnvFromSource `yaml:"envFrom,omitempty" validate:"dive"`
	VolumeMounts []admin.VolumeMount   `yaml:"volumeMounts,omitempty" validate:"dive"`
	Volumes      []admin.Volume        `yaml:"volumes,omitempty" validate:"dive"`
}

// PortalAppConfig represents the portal application-specific configuration
// Similar to AdminAppConfig but with portal-specific field names
type PortalAppConfig struct {
	// Server settings
	Addr       string `yaml:"addr" validate:"required,ip"`
	Port       string `yaml:"port" validate:"required,port_string"`
	BaseDomain string `yaml:"baseDomain" validate:"required,fqdn"`
	LogLevel   string `yaml:"logLevel" validate:"required,oneof=debug info warn error"`

	// Database configuration (portal-specific naming)
	PortalDB PortalDBConfig `yaml:"portalDb" validate:"required"`

	// Authentication configuration
	Auth admin.AuthAppConfig `yaml:"auth" validate:"required"`

	// Redis configuration
	Redis admin.RedisAppConfig `yaml:"redis" validate:"required"`

	// CORS configuration
	CORS admin.CORSAppConfig `yaml:"cors" validate:"required"`

	// TLS configuration
	TLS admin.TLSAppConfig `yaml:"tls"`

	// OpenTelemetry configuration
	OpenTelemetry PortalOpenTelemetryAppConfig `yaml:"opentelemetry"`
}

// PortalDBConfig represents database configuration in the portal application config
// Similar to AdminDBConfig but with portal-specific naming
type PortalDBConfig struct {
	Host             string `yaml:"host"` // Will be overridden by env var
	Port             int    `yaml:"port" validate:"min=1,max=65535"`
	User             string `yaml:"user"`     // Will be overridden by env var
	Password         string `yaml:"password"` // Will be overridden by env var
	DBName           string `yaml:"dbName"`   // Will be overridden by env var
	InitialConnRetry int    `yaml:"initialConnRetry" validate:"min=1"`
}

// PortalOpenTelemetryAppConfig represents OpenTelemetry configuration for portal
// Similar to admin but with portal-specific service name defaults
type PortalOpenTelemetryAppConfig struct {
	Enabled        bool   `yaml:"enabled"`
	ServiceName    string `yaml:"serviceName,omitempty"`
	TracesExporter string `yaml:"tracesExporter,omitempty" validate:"omitempty,oneof=otlp jaeger zipkin console"`
	OTLPEndpoint   string `yaml:"otlpEndpoint,omitempty" validate:"omitempty,url"`
	OTLPProtocol   string `yaml:"otlpProtocol,omitempty" validate:"omitempty,oneof=grpc http"`
}

// PortalServiceConfig represents portal service configuration
// Identical to AdminServiceConfig but with portal context
type PortalServiceConfig = admin.AdminServiceConfig

// Validation methods

// Validate validates the entire Values configuration
func (v *Values) Validate() error {
	return helmcharts.ValidateStruct(v)
}

// Validate validates the PortalConfig
func (p *PortalConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}

// Validate validates the PortalAppConfig
func (p *PortalAppConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}

// Validate validates the PortalDBConfig
func (p *PortalDBConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}

// Validate validates the PortalOpenTelemetryAppConfig
func (p *PortalOpenTelemetryAppConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}
