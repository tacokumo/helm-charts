package tacokumo_portal

import (
	"fmt"

	helmcharts "github.com/tacokumo/helm-charts"
)

// Values represents the root configuration for tacokumo-portal Helm chart
type Values struct {
	API APIConfig `yaml:"api" validate:"required"`
}

// APIConfig represents the API service configuration
type APIConfig struct {
	// PortalName is the portal namespace name (REQUIRED, maps to PORTAL_NAME env var)
	PortalName string `yaml:"portalName" validate:"required"`

	// LogLevel is the logging level (optional: debug, info, warn, error)
	LogLevel string `yaml:"logLevel" validate:"omitempty,oneof=debug info warn error"`

	// Container image configuration
	Image helmcharts.Image `yaml:"image" validate:"required"`

	// HPA configuration
	HPA HPAConfig `yaml:"hpa"`

	// Service configuration
	Service ServiceConfig `yaml:"service"`

	// Resource limits and requests
	Resources ResourceConfig `yaml:"resources,omitempty"`

	// Health check probes
	LivenessProbe  ProbeConfig `yaml:"livenessProbe"`
	ReadinessProbe ProbeConfig `yaml:"readinessProbe"`

	// Pod-level security context
	SecurityContext SecurityContext `yaml:"securityContext"`

	// Container-level security context
	ContainerSecurityContext SecurityContext `yaml:"containerSecurityContext"`

	// RBAC configuration
	RBAC RBACConfig `yaml:"rbac"`

	// ServiceAccount configuration
	ServiceAccount ServiceAccountConfig `yaml:"serviceAccount"`

	// TerminationGracePeriodSeconds for the pod
	TerminationGracePeriodSeconds int64 `yaml:"terminationGracePeriodSeconds" validate:"omitempty,min=0"`

	// Annotations for the deployment
	Annotations map[string]string `yaml:"annotations,omitempty"`

	// Annotations for pods
	PodAnnotations map[string]string `yaml:"podAnnotations,omitempty"`

	// Labels for the deployment
	Labels map[string]string `yaml:"labels,omitempty"`

	// Labels for pods
	PodLabels map[string]string `yaml:"podLabels,omitempty"`

	// NodeSelector for pod scheduling
	NodeSelector map[string]string `yaml:"nodeSelector,omitempty"`

	// Tolerations for pod scheduling
	Tolerations helmcharts.Tolerations `yaml:"tolerations,omitempty"`

	// Affinity for pod scheduling
	Affinity *helmcharts.Affinity `yaml:"affinity,omitempty"`

	// ImagePullSecrets for pulling container images
	ImagePullSecrets []ImagePullSecret `yaml:"imagePullSecrets,omitempty" validate:"dive"`

	// Additional environment variables
	Env []EnvVar `yaml:"env,omitempty" validate:"dive"`

	// Environment variables from ConfigMaps or Secrets
	EnvFrom []EnvFromSource `yaml:"envFrom,omitempty" validate:"dive"`
}

// HPAConfig represents HorizontalPodAutoscaler configuration
type HPAConfig struct {
	Enabled                           bool `yaml:"enabled"`
	MinReplicas                       int  `yaml:"minReplicas" validate:"required_if=Enabled true,omitempty,min=1"`
	MaxReplicas                       int  `yaml:"maxReplicas" validate:"required_if=Enabled true,omitempty,min=1,gtefield=MinReplicas"`
	TargetMemoryUtilizationPercentage int  `yaml:"targetMemoryUtilizationPercentage" validate:"required_if=Enabled true,omitempty,min=1,max=100"`
}

// ServiceConfig represents Kubernetes Service configuration
type ServiceConfig struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type,omitempty" validate:"omitempty,oneof=ClusterIP NodePort LoadBalancer"`
	Port    int    `yaml:"port" validate:"required_if=Enabled true,omitempty,min=1,max=65535"`
}

// ResourceConfig represents container resource limits and requests
type ResourceConfig struct {
	Limits   ResourceSpec `yaml:"limits,omitempty"`
	Requests ResourceSpec `yaml:"requests,omitempty"`
}

// ResourceSpec represents CPU and memory resource specifications
type ResourceSpec struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// ProbeConfig represents health check probe configuration
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
	Host        string       `yaml:"host,omitempty"`
	Scheme      string       `yaml:"scheme,omitempty" validate:"omitempty,oneof=HTTP HTTPS"`
	HTTPHeaders []HTTPHeader `yaml:"httpHeaders,omitempty" validate:"dive"`
}

// TCPSocketAction represents TCP socket action for probes
type TCPSocketAction struct {
	Port int    `yaml:"port" validate:"required,min=1,max=65535"`
	Host string `yaml:"host,omitempty"`
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

// SecurityContext represents pod or container security context
type SecurityContext struct {
	RunAsUser                *int64        `yaml:"runAsUser,omitempty" validate:"omitempty,min=0"`
	RunAsGroup               *int64        `yaml:"runAsGroup,omitempty" validate:"omitempty,min=0"`
	RunAsNonRoot             *bool         `yaml:"runAsNonRoot,omitempty"`
	ReadOnlyRootFilesystem   *bool         `yaml:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation *bool         `yaml:"allowPrivilegeEscalation,omitempty"`
	Capabilities             *Capabilities `yaml:"capabilities,omitempty"`
	SeccompProfile           *SeccompProfile `yaml:"seccompProfile,omitempty"`
}

// Capabilities represents security capabilities
type Capabilities struct {
	Add  []string `yaml:"add,omitempty"`
	Drop []string `yaml:"drop,omitempty"`
}

// SeccompProfile represents seccomp profile configuration
type SeccompProfile struct {
	Type             string `yaml:"type" validate:"required,oneof=Unconfined RuntimeDefault Localhost"`
	LocalhostProfile string `yaml:"localhostProfile,omitempty" validate:"required_if=Type Localhost"`
}

// RBACConfig represents RBAC configuration
type RBACConfig struct {
	Create bool `yaml:"create"`
}

// ServiceAccountConfig represents ServiceAccount configuration
type ServiceAccountConfig struct {
	Create      bool              `yaml:"create"`
	Name        string            `yaml:"name" validate:"required_if=Create true"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// ImagePullSecret represents image pull secret configuration
type ImagePullSecret struct {
	Name string `yaml:"name" validate:"required"`
}

// EnvVar represents environment variable configuration
type EnvVar struct {
	Name      string        `yaml:"name" validate:"required"`
	Value     string        `yaml:"value,omitempty"`
	ValueFrom *EnvVarSource `yaml:"valueFrom,omitempty"`
}

// EnvVarSource represents environment variable source
type EnvVarSource struct {
	FieldRef         *ObjectFieldSelector   `yaml:"fieldRef,omitempty"`
	ResourceFieldRef *ResourceFieldSelector `yaml:"resourceFieldRef,omitempty"`
	ConfigMapKeyRef  *ConfigMapKeySelector  `yaml:"configMapKeyRef,omitempty"`
	SecretKeyRef     *SecretKeySelector     `yaml:"secretKeyRef,omitempty"`
}

// ObjectFieldSelector represents object field selector
type ObjectFieldSelector struct {
	APIVersion string `yaml:"apiVersion,omitempty"`
	FieldPath  string `yaml:"fieldPath" validate:"required"`
}

// ResourceFieldSelector represents resource field selector
type ResourceFieldSelector struct {
	ContainerName string `yaml:"containerName,omitempty"`
	Resource      string `yaml:"resource" validate:"required"`
	Divisor       string `yaml:"divisor,omitempty"`
}

// ConfigMapKeySelector represents ConfigMap key selector
type ConfigMapKeySelector struct {
	Name     string `yaml:"name" validate:"required"`
	Key      string `yaml:"key" validate:"required"`
	Optional *bool  `yaml:"optional,omitempty"`
}

// SecretKeySelector represents Secret key selector
type SecretKeySelector struct {
	Name     string `yaml:"name" validate:"required"`
	Key      string `yaml:"key" validate:"required"`
	Optional *bool  `yaml:"optional,omitempty"`
}

// EnvFromSource represents environment variable source
type EnvFromSource struct {
	Prefix       string              `yaml:"prefix,omitempty"`
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

// Validate validates the entire Values configuration
func (v *Values) Validate() error {
	return helmcharts.ValidateStruct(v)
}

// Validate validates the APIConfig
func (a *APIConfig) Validate() error {
	if err := helmcharts.ValidateStruct(a); err != nil {
		return err
	}
	// Validate nested ServiceConfig with custom validation
	if err := a.Service.Validate(); err != nil {
		return err
	}
	// Validate nested HPAConfig with custom validation
	if err := a.HPA.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates the HPAConfig
func (h *HPAConfig) Validate() error {
	if err := helmcharts.ValidateStruct(h); err != nil {
		return err
	}
	// Additional validation: when enabled, all fields must be properly set
	if h.Enabled {
		if h.MinReplicas < 1 {
			return fmt.Errorf("HPA.MinReplicas: must be at least 1 when HPA is enabled")
		}
		if h.MaxReplicas < h.MinReplicas {
			return fmt.Errorf("HPA.MaxReplicas: must be greater than or equal to MinReplicas")
		}
		if h.TargetMemoryUtilizationPercentage < 1 || h.TargetMemoryUtilizationPercentage > 100 {
			return fmt.Errorf("HPA.TargetMemoryUtilizationPercentage: must be between 1 and 100")
		}
	}
	return nil
}

// Validate validates the ServiceConfig
func (s *ServiceConfig) Validate() error {
	if err := helmcharts.ValidateStruct(s); err != nil {
		return err
	}
	// Additional validation: when enabled, port must be valid
	if s.Enabled && s.Port == 0 {
		return fmt.Errorf("Service.Port: port is required when service is enabled")
	}
	return nil
}

// Validate validates the ProbeConfig
func (p *ProbeConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}
