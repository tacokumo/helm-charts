package tacokumo_application

import (
	helmcharts "github.com/tacokumo/helm-charts"
)

// Values represents the root configuration for tacokumo-application Helm chart
type Values struct {
	Main MainConfig `yaml:"main" validate:"required"`
}

// MainConfig represents the main application configuration
type MainConfig struct {
	ApplicationName  string            `yaml:"applicationName" validate:"required"`
	ReplicaCount     int               `yaml:"replicaCount" validate:"min=1"`
	Image            string            `yaml:"image" validate:"required"`
	ImagePullSecrets []ImagePullSecret `yaml:"imagePullSecrets,omitempty" validate:"dive"`
	ImagePullPolicy  string            `yaml:"imagePullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`

	// Annotations for various Kubernetes resources
	Annotations    map[string]string `yaml:"annotations,omitempty"`
	PodAnnotations map[string]string `yaml:"podAnnotations,omitempty"`

	// Environment configuration
	EnvFrom []EnvFromSource `yaml:"envFrom,omitempty" validate:"dive"`

	// Health check probes
	LivenessProbe  ProbeConfig `yaml:"livenessProbe"`
	ReadinessProbe ProbeConfig `yaml:"readinessProbe"`
	StartupProbe   ProbeConfig `yaml:"startupProbe"`
}

// ImagePullSecret represents image pull secret configuration
type ImagePullSecret struct {
	Name string `yaml:"name" validate:"required"`
}

// EnvFromSource represents environment variable source configuration
type EnvFromSource struct {
	ConfigMapRef *ConfigMapEnvSource `yaml:"configMapRef,omitempty"`
	SecretRef    *SecretEnvSource    `yaml:"secretRef,omitempty"`
}

// ConfigMapEnvSource represents ConfigMap environment source
type ConfigMapEnvSource struct {
	Name     string `yaml:"name" validate:"required"`
	Optional *bool  `yaml:"optional,omitempty"`
}

// SecretEnvSource represents Secret environment source
type SecretEnvSource struct {
	Name     string `yaml:"name" validate:"required"`
	Optional *bool  `yaml:"optional,omitempty"`
}

// ProbeConfig represents health check probe configuration
// This is a flexible configuration that can be empty (disabled) or contain probe settings
type ProbeConfig struct {
	// HTTP probe configuration
	HTTPGet *HTTPGetAction `yaml:"httpGet,omitempty"`

	// TCP probe configuration
	TCPSocket *TCPSocketAction `yaml:"tcpSocket,omitempty"`

	// Exec probe configuration
	Exec *ExecAction `yaml:"exec,omitempty"`

	// Common probe settings
	InitialDelaySeconds *int `yaml:"initialDelaySeconds,omitempty" validate:"omitempty,min=0"`
	PeriodSeconds       *int `yaml:"periodSeconds,omitempty" validate:"omitempty,min=1"`
	TimeoutSeconds      *int `yaml:"timeoutSeconds,omitempty" validate:"omitempty,min=1"`
	SuccessThreshold    *int `yaml:"successThreshold,omitempty" validate:"omitempty,min=1"`
	FailureThreshold    *int `yaml:"failureThreshold,omitempty" validate:"omitempty,min=1"`
}

// HTTPGetAction represents HTTP GET action for probes
type HTTPGetAction struct {
	Path        string       `yaml:"path" validate:"required"`
	Port        int          `yaml:"port" validate:"required,min=1,max=65535"`
	Host        string       `yaml:"host,omitempty" validate:"omitempty,fqdn|ip"`
	Scheme      string       `yaml:"scheme,omitempty" validate:"omitempty,oneof=HTTP HTTPS"`
	HTTPHeaders []HTTPHeader `yaml:"httpHeaders,omitempty" validate:"dive"`
}

// TCPSocketAction represents TCP socket action for probes
type TCPSocketAction struct {
	Port int    `yaml:"port" validate:"required,min=1,max=65535"`
	Host string `yaml:"host,omitempty" validate:"omitempty,fqdn|ip"`
}

// ExecAction represents exec action for probes
type ExecAction struct {
	Command []string `yaml:"command" validate:"required,min=1"`
}

// HTTPHeader represents HTTP header for probes
type HTTPHeader struct {
	Name  string `yaml:"name" validate:"required"`
	Value string `yaml:"value" validate:"required"`
}

// Validate validates the entire Values configuration
func (v *Values) Validate() error {
	return helmcharts.ValidateStruct(v)
}

// Validate validates the MainConfig
func (m *MainConfig) Validate() error {
	return helmcharts.ValidateStruct(m)
}

// Validate validates the ProbeConfig
func (p *ProbeConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}
