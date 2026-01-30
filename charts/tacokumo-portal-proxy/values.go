package tacokumo_portal_proxy

import (
	helmcharts "github.com/tacokumo/helm-charts"
)

// Values represents the root configuration for tacokumo-portal-proxy Helm chart
type Values struct {
	PortalProxy PortalProxyConfig `yaml:"portalProxy" validate:"required"`
}

// PortalProxyConfig represents the portal proxy service configuration
type PortalProxyConfig struct {
	ReplicaCount int    `yaml:"replicaCount" validate:"min=1"`
	BaseDomain   string `yaml:"baseDomain" validate:"required,fqdn"`

	// Container image configuration
	Image helmcharts.Image `yaml:"image" validate:"required"`

	// Service configuration
	Service ProxyServiceConfig `yaml:"service" validate:"required"`

	// Resource limits and requests
	Resources helmcharts.Resources `yaml:"resources"`

	// Health check probes
	LivenessProbe  helmcharts.HTTPProbe `yaml:"livenessProbe"`
	ReadinessProbe helmcharts.HTTPProbe `yaml:"readinessProbe"`

	// Security context
	SecurityContext helmcharts.SecurityContext `yaml:"securityContext"`

	// Additional configurations
	Annotations      map[string]string      `yaml:"annotations,omitempty"`
	PodAnnotations   map[string]string      `yaml:"podAnnotations,omitempty"`
	NodeSelector     map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations      helmcharts.Tolerations `yaml:"tolerations,omitempty"`
	Affinity         *helmcharts.Affinity   `yaml:"affinity,omitempty"`
	ImagePullSecrets []ImagePullSecret      `yaml:"imagePullSecrets,omitempty" validate:"dive"`

	// Pod Disruption Budget
	PodDisruptionBudget helmcharts.PodDisruptionBudget `yaml:"podDisruptionBudget"`

	// Environment variables
	Env     []EnvVar        `yaml:"env,omitempty" validate:"dive"`
	EnvFrom []EnvFromSource `yaml:"envFrom,omitempty" validate:"dive"`

	// Volume mounts
	VolumeMounts []VolumeMount `yaml:"volumeMounts,omitempty" validate:"dive"`
	Volumes      []Volume      `yaml:"volumes,omitempty" validate:"dive"`
}

// ProxyServiceConfig represents proxy service specific configuration
type ProxyServiceConfig struct {
	Type        string            `yaml:"type" validate:"oneof=ClusterIP NodePort LoadBalancer ExternalName"`
	HTTPPort    int               `yaml:"httpPort" validate:"min=1,max=65535"`
	MetricsPort int               `yaml:"metricsPort" validate:"min=1,max=65535"`
	Annotations map[string]string `yaml:"annotations,omitempty"`

	// Additional service ports
	ExtraPorts []ServicePort `yaml:"extraPorts,omitempty" validate:"dive"`

	// Load balancer configuration (for LoadBalancer type)
	LoadBalancerIP           string   `yaml:"loadBalancerIP,omitempty" validate:"omitempty,ip"`
	LoadBalancerSourceRanges []string `yaml:"loadBalancerSourceRanges,omitempty" validate:"dive,cidr"`

	// External traffic policy (for NodePort and LoadBalancer)
	ExternalTrafficPolicy string `yaml:"externalTrafficPolicy,omitempty" validate:"omitempty,oneof=Cluster Local"`
}

// ServicePort represents additional service port configuration
type ServicePort struct {
	Name       string `yaml:"name" validate:"required"`
	Port       int    `yaml:"port" validate:"required,min=1,max=65535"`
	TargetPort int    `yaml:"targetPort,omitempty" validate:"omitempty,min=1,max=65535"`
	Protocol   string `yaml:"protocol,omitempty" validate:"omitempty,oneof=TCP UDP SCTP"`
	NodePort   int    `yaml:"nodePort,omitempty" validate:"omitempty,min=30000,max=32767"`
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
}

// Volume represents volume configuration
type Volume struct {
	Name                  string                       `yaml:"name" validate:"required"`
	HostPath              *HostPathVolumeSource        `yaml:"hostPath,omitempty"`
	EmptyDir              *EmptyDirVolumeSource        `yaml:"emptyDir,omitempty"`
	Secret                *SecretVolumeSource          `yaml:"secret,omitempty"`
	ConfigMap             *ConfigMapVolumeSource       `yaml:"configMap,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaimSource `yaml:"persistentVolumeClaim,omitempty"`
	Projected             *ProjectedVolumeSource       `yaml:"projected,omitempty"`
}

// HostPathVolumeSource represents host path volume source
type HostPathVolumeSource struct {
	Path string  `yaml:"path" validate:"required,filepath"`
	Type *string `yaml:"type,omitempty" validate:"omitempty,oneof='' DirectoryOrCreate Directory FileOrCreate File Socket CharDevice BlockDevice"`
}

// EmptyDirVolumeSource represents empty dir volume source
type EmptyDirVolumeSource struct {
	Medium    string `yaml:"medium,omitempty" validate:"omitempty,oneof='' Memory"`
	SizeLimit string `yaml:"sizeLimit,omitempty" validate:"omitempty,resource_quantity"`
}

// SecretVolumeSource represents secret volume source
type SecretVolumeSource struct {
	SecretName  string      `yaml:"secretName" validate:"required"`
	Items       []KeyToPath `yaml:"items,omitempty" validate:"dive"`
	DefaultMode *int32      `yaml:"defaultMode,omitempty" validate:"omitempty,min=0,max=511"`
	Optional    *bool       `yaml:"optional,omitempty"`
}

// ConfigMapVolumeSource represents config map volume source
type ConfigMapVolumeSource struct {
	Name        string      `yaml:"name" validate:"required"`
	Items       []KeyToPath `yaml:"items,omitempty" validate:"dive"`
	DefaultMode *int32      `yaml:"defaultMode,omitempty" validate:"omitempty,min=0,max=511"`
	Optional    *bool       `yaml:"optional,omitempty"`
}

// PersistentVolumeClaimSource represents PVC source
type PersistentVolumeClaimSource struct {
	ClaimName string `yaml:"claimName" validate:"required"`
	ReadOnly  bool   `yaml:"readOnly,omitempty"`
}

// ProjectedVolumeSource represents projected volume source
type ProjectedVolumeSource struct {
	Sources     []VolumeProjection `yaml:"sources" validate:"required,dive"`
	DefaultMode *int32             `yaml:"defaultMode,omitempty" validate:"omitempty,min=0,max=511"`
}

// VolumeProjection represents volume projection
type VolumeProjection struct {
	Secret              *SecretProjection              `yaml:"secret,omitempty"`
	ConfigMap           *ConfigMapProjection           `yaml:"configMap,omitempty"`
	DownwardAPI         *DownwardAPIProjection         `yaml:"downwardAPI,omitempty"`
	ServiceAccountToken *ServiceAccountTokenProjection `yaml:"serviceAccountToken,omitempty"`
}

// SecretProjection represents secret projection
type SecretProjection struct {
	Name     string      `yaml:"name" validate:"required"`
	Items    []KeyToPath `yaml:"items,omitempty" validate:"dive"`
	Optional *bool       `yaml:"optional,omitempty"`
}

// ConfigMapProjection represents config map projection
type ConfigMapProjection struct {
	Name     string      `yaml:"name" validate:"required"`
	Items    []KeyToPath `yaml:"items,omitempty" validate:"dive"`
	Optional *bool       `yaml:"optional,omitempty"`
}

// DownwardAPIProjection represents downward API projection
type DownwardAPIProjection struct {
	Items []DownwardAPIVolumeFile `yaml:"items" validate:"required,dive"`
}

// ServiceAccountTokenProjection represents service account token projection
type ServiceAccountTokenProjection struct {
	Audience          string `yaml:"audience,omitempty"`
	ExpirationSeconds *int64 `yaml:"expirationSeconds,omitempty" validate:"omitempty,min=600"`
	Path              string `yaml:"path" validate:"required,filepath"`
}

// KeyToPath represents key to path mapping
type KeyToPath struct {
	Key  string `yaml:"key" validate:"required"`
	Path string `yaml:"path" validate:"required,filepath"`
	Mode *int32 `yaml:"mode,omitempty" validate:"omitempty,min=0,max=511"`
}

// DownwardAPIVolumeFile represents downward API volume file
type DownwardAPIVolumeFile struct {
	Path             string                 `yaml:"path" validate:"required,filepath"`
	FieldRef         *ObjectFieldSelector   `yaml:"fieldRef,omitempty"`
	ResourceFieldRef *ResourceFieldSelector `yaml:"resourceFieldRef,omitempty"`
	Mode             *int32                 `yaml:"mode,omitempty" validate:"omitempty,min=0,max=511"`
}

// Validate validates the entire Values configuration
func (v *Values) Validate() error {
	return helmcharts.ValidateStruct(v)
}

// Validate validates the PortalProxyConfig
func (p *PortalProxyConfig) Validate() error {
	return helmcharts.ValidateStruct(p)
}

// Validate validates the ProxyServiceConfig
func (s *ProxyServiceConfig) Validate() error {
	return helmcharts.ValidateStruct(s)
}
