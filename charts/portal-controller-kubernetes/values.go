package portal_controller_kubernetes

import (
	helmcharts "github.com/tacokumo/helm-charts"
)

// Values represents the root configuration for portal-controller-kubernetes Helm chart
type Values struct {
	CRDs       CRDsConfig       `yaml:"crds" validate:"required"`
	Controller ControllerConfig `yaml:"controller" validate:"required"`
}

// CRDsConfig represents Custom Resource Definitions configuration
type CRDsConfig struct {
	Install bool `yaml:"install"`
}

// ControllerConfig represents the controller deployment configuration
type ControllerConfig struct {
	// Pod configuration
	TerminationGracePeriodSeconds int64 `yaml:"terminationGracePeriodSeconds" validate:"min=0"`

	// Kubernetes resource metadata
	Affinity         *helmcharts.Affinity `yaml:"affinity,omitempty"`
	Labels           map[string]string    `yaml:"labels,omitempty"`
	Annotations      map[string]string    `yaml:"annotations,omitempty"`
	PodAnnotations   map[string]string    `yaml:"podAnnotations,omitempty"`
	PodLabels        map[string]string    `yaml:"podLabels,omitempty"`
	ImagePullSecrets []ImagePullSecret    `yaml:"imagePullSecrets,omitempty" validate:"dive"`

	// Security context for the pod
	SecurityContext PodSecurityContext `yaml:"securityContext"`

	// Manager container configuration
	ManagerContainer ManagerContainerConfig `yaml:"managerContainer" validate:"required"`

	// Additional configurations
	NodeSelector      map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations       helmcharts.Tolerations `yaml:"tolerations,omitempty"`
	PriorityClassName string                 `yaml:"priorityClassName,omitempty"`

	// Service account
	ServiceAccount ServiceAccountConfig `yaml:"serviceAccount"`

	// Replica configuration
	ReplicaCount int `yaml:"replicaCount,omitempty" validate:"omitempty,min=1"`

	// Pod Disruption Budget
	PodDisruptionBudget helmcharts.PodDisruptionBudget `yaml:"podDisruptionBudget"`

	// Metrics configuration
	Metrics MetricsConfig `yaml:"metrics"`
}

// ManagerContainerConfig represents the manager container configuration
type ManagerContainerConfig struct {
	// Container image
	Image ContainerImage `yaml:"image" validate:"required"`

	// Command line arguments
	ExtraArgs []string `yaml:"extraArgs,omitempty"`

	// Health check probes
	LivenessProbe  HTTPProbeConfig `yaml:"livenessProbe"`
	ReadinessProbe HTTPProbeConfig `yaml:"readinessProbe"`

	// Security context for the container
	SecurityContext ContainerSecurityContext `yaml:"securityContext"`

	// Resource limits and requests
	Resources helmcharts.Resources `yaml:"resources"`

	// Environment variables
	Env     []EnvVar        `yaml:"env,omitempty" validate:"dive"`
	EnvFrom []EnvFromSource `yaml:"envFrom,omitempty" validate:"dive"`

	// Volume mounts
	VolumeMounts []VolumeMount `yaml:"volumeMounts,omitempty" validate:"dive"`

	// Container ports
	Ports []ContainerPort `yaml:"ports,omitempty" validate:"dive"`

	// Working directory
	WorkingDir string `yaml:"workingDir,omitempty" validate:"omitempty,filepath"`

	// Command and args
	Command []string `yaml:"command,omitempty"`
	Args    []string `yaml:"args,omitempty"`
}

// ContainerImage represents container image configuration
type ContainerImage struct {
	Repository string `yaml:"repository" validate:"required"`
	Tag        string `yaml:"tag" validate:"required"`
	PullPolicy string `yaml:"pullPolicy,omitempty" validate:"omitempty,oneof=Always IfNotPresent Never"`
}

// HTTPProbeConfig represents HTTP health check probe configuration
type HTTPProbeConfig struct {
	HTTPGet             HTTPGetAction `yaml:"httpGet" validate:"required"`
	InitialDelaySeconds int           `yaml:"initialDelaySeconds" validate:"min=0"`
	PeriodSeconds       int           `yaml:"periodSeconds" validate:"min=1"`
	TimeoutSeconds      int           `yaml:"timeoutSeconds,omitempty" validate:"omitempty,min=1"`
	SuccessThreshold    int           `yaml:"successThreshold,omitempty" validate:"omitempty,min=1"`
	FailureThreshold    int           `yaml:"failureThreshold,omitempty" validate:"omitempty,min=1"`
}

// HTTPGetAction represents HTTP GET action for probes
type HTTPGetAction struct {
	Path        string       `yaml:"path" validate:"required"`
	Port        int          `yaml:"port" validate:"required,min=1,max=65535"`
	Host        string       `yaml:"host,omitempty" validate:"omitempty,fqdn|ip"`
	Scheme      string       `yaml:"scheme,omitempty" validate:"omitempty,oneof=HTTP HTTPS"`
	HTTPHeaders []HTTPHeader `yaml:"httpHeaders,omitempty" validate:"dive"`
}

// HTTPHeader represents HTTP header for probes
type HTTPHeader struct {
	Name  string `yaml:"name" validate:"required"`
	Value string `yaml:"value" validate:"required"`
}

// PodSecurityContext represents pod-level security context
type PodSecurityContext struct {
	RunAsUser           *int64          `yaml:"runAsUser,omitempty" validate:"omitempty,min=0"`
	RunAsGroup          *int64          `yaml:"runAsGroup,omitempty" validate:"omitempty,min=0"`
	RunAsNonRoot        *bool           `yaml:"runAsNonRoot,omitempty"`
	FSGroup             *int64          `yaml:"fsGroup,omitempty" validate:"omitempty,min=0"`
	FSGroupChangePolicy *string         `yaml:"fsGroupChangePolicy,omitempty" validate:"omitempty,oneof=Always OnRootMismatch"`
	SeccompProfile      *SeccompProfile `yaml:"seccompProfile,omitempty"`
	SELinuxOptions      *SELinuxOptions `yaml:"seLinuxOptions,omitempty"`
	SupplementalGroups  []int64         `yaml:"supplementalGroups,omitempty" validate:"dive,min=0"`
	Sysctls             []Sysctl        `yaml:"sysctls,omitempty" validate:"dive"`
	WindowsOptions      *WindowsOptions `yaml:"windowsOptions,omitempty"`
}

// ContainerSecurityContext represents container-level security context
type ContainerSecurityContext struct {
	RunAsUser                *int64          `yaml:"runAsUser,omitempty" validate:"omitempty,min=0"`
	RunAsGroup               *int64          `yaml:"runAsGroup,omitempty" validate:"omitempty,min=0"`
	RunAsNonRoot             *bool           `yaml:"runAsNonRoot,omitempty"`
	ReadOnlyRootFilesystem   *bool           `yaml:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation *bool           `yaml:"allowPrivilegeEscalation,omitempty"`
	Privileged               *bool           `yaml:"privileged,omitempty"`
	Capabilities             *Capabilities   `yaml:"capabilities,omitempty"`
	SeccompProfile           *SeccompProfile `yaml:"seccompProfile,omitempty"`
	SELinuxOptions           *SELinuxOptions `yaml:"seLinuxOptions,omitempty"`
	WindowsOptions           *WindowsOptions `yaml:"windowsOptions,omitempty"`
}

// Capabilities represents security capabilities
type Capabilities struct {
	Add  []string `yaml:"add,omitempty"`
	Drop []string `yaml:"drop,omitempty"`
}

// SeccompProfile represents seccomp profile
type SeccompProfile struct {
	Type             string `yaml:"type" validate:"required,oneof=RuntimeDefault Localhost Unconfined"`
	LocalhostProfile string `yaml:"localhostProfile,omitempty" validate:"required_if=Type Localhost,filepath"`
}

// SELinuxOptions represents SELinux options
type SELinuxOptions struct {
	User  string `yaml:"user,omitempty"`
	Role  string `yaml:"role,omitempty"`
	Type  string `yaml:"type,omitempty"`
	Level string `yaml:"level,omitempty"`
}

// Sysctl represents a sysctl setting
type Sysctl struct {
	Name  string `yaml:"name" validate:"required"`
	Value string `yaml:"value" validate:"required"`
}

// WindowsOptions represents Windows-specific options
type WindowsOptions struct {
	GMSACredentialSpecName string `yaml:"gmsaCredentialSpecName,omitempty"`
	GMSACredentialSpec     string `yaml:"gmsaCredentialSpec,omitempty"`
	RunAsUserName          string `yaml:"runAsUserName,omitempty"`
	HostProcess            *bool  `yaml:"hostProcess,omitempty"`
}

// ImagePullSecret represents image pull secret configuration
type ImagePullSecret struct {
	Name string `yaml:"name" validate:"required"`
}

// ServiceAccountConfig represents service account configuration
type ServiceAccountConfig struct {
	Create      bool              `yaml:"create"`
	Name        string            `yaml:"name,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Port        int               `yaml:"port" validate:"required_if=Enabled true,omitempty,min=1,max=65535"`
	Path        string            `yaml:"path" validate:"required_if=Enabled true"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
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
	Divisor       string `yaml:"divisor,omitempty" validate:"omitempty,resource_quantity"`
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

// VolumeMount represents volume mount configuration
type VolumeMount struct {
	Name             string `yaml:"name" validate:"required"`
	MountPath        string `yaml:"mountPath" validate:"required,filepath"`
	SubPath          string `yaml:"subPath,omitempty"`
	ReadOnly         bool   `yaml:"readOnly,omitempty"`
	MountPropagation string `yaml:"mountPropagation,omitempty" validate:"omitempty,oneof=None HostToContainer Bidirectional"`
	SubPathExpr      string `yaml:"subPathExpr,omitempty"`
}

// ContainerPort represents container port configuration
type ContainerPort struct {
	Name          string `yaml:"name,omitempty"`
	ContainerPort int    `yaml:"containerPort" validate:"required,min=1,max=65535"`
	Protocol      string `yaml:"protocol,omitempty" validate:"omitempty,oneof=TCP UDP SCTP"`
	HostIP        string `yaml:"hostIP,omitempty" validate:"omitempty,ip"`
	HostPort      int    `yaml:"hostPort,omitempty" validate:"omitempty,min=0,max=65535"`
}

// Validate validates the entire Values configuration
func (v *Values) Validate() error {
	return helmcharts.ValidateStruct(v)
}

// Validate validates the ControllerConfig
func (c *ControllerConfig) Validate() error {
	return helmcharts.ValidateStruct(c)
}

// Validate validates the ManagerContainerConfig
func (m *ManagerContainerConfig) Validate() error {
	return helmcharts.ValidateStruct(m)
}

// Validate validates the HTTPProbeConfig
func (h *HTTPProbeConfig) Validate() error {
	return helmcharts.ValidateStruct(h)
}
