package tacokumo_admin

import (
	helmcharts "github.com/tacokumo/helm-charts"
)

// Values represents the root configuration for tacokumo-admin Helm chart
type Values struct {
	Global  GlobalConfig       `yaml:"global" validate:"required"`
	Admin   AdminConfig        `yaml:"admin" validate:"required"`
	Ingress helmcharts.Ingress `yaml:"ingress"`
}

// GlobalConfig represents global configuration shared across components
type GlobalConfig struct {
	ExternalServices ExternalServicesConfig `yaml:"externalServices" validate:"required"`
}

// ExternalServicesConfig represents configuration for external service dependencies
type ExternalServicesConfig struct {
	PostgreSQL PostgreSQLExternalConfig `yaml:"postgresql" validate:"required"`
	Redis      RedisExternalConfig      `yaml:"redis" validate:"required"`
}

// PostgreSQLExternalConfig represents PostgreSQL external service configuration
type PostgreSQLExternalConfig struct {
	Host             string `yaml:"host" validate:"required,fqdn|ip"`
	Port             int    `yaml:"port" validate:"required,min=1,max=65535"`
	Database         string `yaml:"database" validate:"required"`
	Username         string `yaml:"username" validate:"required"`
	SecretName       string `yaml:"secretName" validate:"required"`
	InitialConnRetry int    `yaml:"initialConnRetry" validate:"min=1"`
}

// RedisExternalConfig represents Redis external service configuration
type RedisExternalConfig struct {
	Host             string `yaml:"host" validate:"required,fqdn|ip"`
	Port             int    `yaml:"port" validate:"required,min=1,max=65535"`
	DB               int    `yaml:"db" validate:"min=0,max=15"`
	SecretName       string `yaml:"secretName" validate:"required"`
	InitialConnRetry int    `yaml:"initialConnRetry" validate:"min=1"`
}

// AdminConfig represents the admin application configuration
type AdminConfig struct {
	ReplicaCount int `yaml:"replicaCount" validate:"min=1"`

	// Application configuration
	Config AdminAppConfig `yaml:"config" validate:"required"`

	// GitHub OAuth configuration
	GitHub GitHubConfig `yaml:"github" validate:"required"`

	// Container image configuration
	Image helmcharts.Image `yaml:"image" validate:"required"`

	// Service configuration
	Service AdminServiceConfig `yaml:"service" validate:"required"`

	// Service account configuration
	ServiceAccount ServiceAccountConfig `yaml:"serviceAccount" validate:"required"`

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
	PodSecurityContext *PodSecurityContext         `yaml:"podSecurityContext,omitempty"`

	// Health checks
	LivenessProbe  *HTTPProbeConfig `yaml:"livenessProbe,omitempty"`
	ReadinessProbe *HTTPProbeConfig `yaml:"readinessProbe,omitempty"`

	// Resource limits and requests
	Resources helmcharts.Resources `yaml:"resources"`

	// Additional configurations
	ImagePullSecrets        []ImagePullSecret              `yaml:"imagePullSecrets,omitempty" validate:"dive"`
	PodDisruptionBudget     helmcharts.PodDisruptionBudget `yaml:"podDisruptionBudget"`
	HorizontalPodAutoscaler HPAConfig                      `yaml:"horizontalPodAutoscaler"`

	// Environment variables and volumes
	Env          []EnvVar        `yaml:"env,omitempty" validate:"dive"`
	EnvFrom      []EnvFromSource `yaml:"envFrom,omitempty" validate:"dive"`
	VolumeMounts []VolumeMount   `yaml:"volumeMounts,omitempty" validate:"dive"`
	Volumes      []Volume        `yaml:"volumes,omitempty" validate:"dive"`
}

// AdminAppConfig represents the application-specific configuration
type AdminAppConfig struct {
	// Server settings
	Addr       string `yaml:"addr" validate:"required,ip"`
	Port       string `yaml:"port" validate:"required,port_string"`
	BaseDomain string `yaml:"baseDomain" validate:"required,fqdn"`
	LogLevel   string `yaml:"logLevel" validate:"required,oneof=debug info warn error"`

	// Database configuration
	AdminDB AdminDBConfig `yaml:"adminDb" validate:"required"`

	// Authentication configuration
	Auth AuthAppConfig `yaml:"auth" validate:"required"`

	// Redis configuration
	Redis RedisAppConfig `yaml:"redis" validate:"required"`

	// CORS configuration
	CORS CORSAppConfig `yaml:"cors" validate:"required"`

	// TLS configuration
	TLS TLSAppConfig `yaml:"tls"`

	// OpenTelemetry configuration
	OpenTelemetry OpenTelemetryAppConfig `yaml:"opentelemetry"`
}

// AdminDBConfig represents database configuration in the application config
type AdminDBConfig struct {
	Host             string `yaml:"host"` // Will be overridden by env var
	Port             int    `yaml:"port" validate:"min=1,max=65535"`
	User             string `yaml:"user"`     // Will be overridden by env var
	Password         string `yaml:"password"` // Will be overridden by env var
	DBName           string `yaml:"dbName"`   // Will be overridden by env var
	InitialConnRetry int    `yaml:"initialConnRetry" validate:"min=1"`
}

// AuthAppConfig represents authentication configuration in the application config
type AuthAppConfig struct {
	ClientID     string `yaml:"clientId"`     // Will be overridden by env var
	ClientSecret string `yaml:"clientSecret"` // Will be overridden by env var
	CallbackURL  string `yaml:"callbackUrl" validate:"required,url"`
	FrontendURL  string `yaml:"frontendUrl" validate:"required,url"`
	GitHubOrg    string `yaml:"githubOrg" validate:"required"`
	SessionTTL   string `yaml:"sessionTtl" validate:"required,duration"`
	CookieSecure bool   `yaml:"cookieSecure"`
}

// RedisAppConfig represents Redis configuration in the application config
type RedisAppConfig struct {
	Host             string `yaml:"host"` // Will be overridden by env var
	Port             int    `yaml:"port" validate:"min=1,max=65535"`
	Password         string `yaml:"password"` // Will be overridden by env var
	DB               int    `yaml:"db" validate:"min=0,max=15"`
	InitialConnRetry int    `yaml:"initialConnRetry" validate:"min=1"`
}

// CORSAppConfig represents CORS configuration in the application config
type CORSAppConfig struct {
	AllowOrigins     string `yaml:"allowOrigins" validate:"required"`
	AllowMethods     string `yaml:"allowMethods" validate:"required"`
	AllowHeaders     string `yaml:"allowHeaders" validate:"required"`
	ExposeHeaders    string `yaml:"exposeHeaders,omitempty"`
	AllowCredentials bool   `yaml:"allowCredentials"`
	MaxAge           int    `yaml:"maxAge" validate:"min=0"`
}

// TLSAppConfig represents TLS configuration in the application config
type TLSAppConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"certFile" validate:"required_if=Enabled true,filepath"`
	KeyFile  string `yaml:"keyFile" validate:"required_if=Enabled true,filepath"`
}

// OpenTelemetryAppConfig represents OpenTelemetry configuration in the application config
type OpenTelemetryAppConfig struct {
	Enabled        bool   `yaml:"enabled"`
	ServiceName    string `yaml:"serviceName,omitempty"`
	TracesExporter string `yaml:"tracesExporter,omitempty" validate:"omitempty,oneof=otlp jaeger zipkin console"`
	OTLPEndpoint   string `yaml:"otlpEndpoint,omitempty" validate:"omitempty,url"`
	OTLPProtocol   string `yaml:"otlpProtocol,omitempty" validate:"omitempty,oneof=grpc http"`
}

// GitHubConfig represents GitHub OAuth configuration
type GitHubConfig struct {
	OAuth GitHubOAuthConfig `yaml:"oauth" validate:"required"`
}

// GitHubOAuthConfig represents GitHub OAuth secret configuration
type GitHubOAuthConfig struct {
	SecretName string `yaml:"secretName" validate:"required"`
}

// AdminServiceConfig represents admin service configuration
type AdminServiceConfig struct {
	Type        string            `yaml:"type" validate:"oneof=ClusterIP NodePort LoadBalancer ExternalName"`
	Port        int               `yaml:"port" validate:"min=1,max=65535"`
	TargetPort  int               `yaml:"targetPort" validate:"min=1,max=65535"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	NodePort    int               `yaml:"nodePort,omitempty" validate:"omitempty,min=30000,max=32767"`

	// Additional service configuration
	SessionAffinity          string                 `yaml:"sessionAffinity,omitempty" validate:"omitempty,oneof=None ClientIP"`
	SessionAffinityConfig    *SessionAffinityConfig `yaml:"sessionAffinityConfig,omitempty"`
	LoadBalancerIP           string                 `yaml:"loadBalancerIP,omitempty" validate:"omitempty,ip"`
	LoadBalancerSourceRanges []string               `yaml:"loadBalancerSourceRanges,omitempty" validate:"dive,cidr"`
	ExternalTrafficPolicy    string                 `yaml:"externalTrafficPolicy,omitempty" validate:"omitempty,oneof=Cluster Local"`
	ExternalName             string                 `yaml:"externalName,omitempty" validate:"omitempty,fqdn"`
}

// SessionAffinityConfig represents session affinity configuration
type SessionAffinityConfig struct {
	ClientIP *ClientIPConfig `yaml:"clientIP,omitempty"`
}

// ClientIPConfig represents client IP configuration for session affinity
type ClientIPConfig struct {
	TimeoutSeconds *int32 `yaml:"timeoutSeconds,omitempty" validate:"omitempty,min=0"`
}

// ServiceAccountConfig represents service account configuration
type ServiceAccountConfig struct {
	Name        string            `yaml:"name" validate:"required"`
	Create      bool              `yaml:"create"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
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

// HTTPProbeConfig represents HTTP health check probe configuration
type HTTPProbeConfig struct {
	HTTPGet             HTTPGetAction `yaml:"httpGet" validate:"required"`
	InitialDelaySeconds int           `yaml:"initialDelaySeconds" validate:"min=0"`
	PeriodSeconds       int           `yaml:"periodSeconds" validate:"min=1"`
	TimeoutSeconds      int           `yaml:"timeoutSeconds" validate:"min=1"`
	SuccessThreshold    int           `yaml:"successThreshold,omitempty" validate:"omitempty,min=1"`
	FailureThreshold    int           `yaml:"failureThreshold" validate:"min=1"`
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

// HPAConfig represents Horizontal Pod Autoscaler configuration
type HPAConfig struct {
	Enabled                           bool         `yaml:"enabled"`
	MinReplicas                       *int32       `yaml:"minReplicas" validate:"required_if=Enabled true,omitempty,min=1"`
	MaxReplicas                       int32        `yaml:"maxReplicas" validate:"required_if=Enabled true,omitempty,min=1"`
	TargetCPUUtilizationPercentage    *int32       `yaml:"targetCPUUtilizationPercentage" validate:"omitempty,min=1,max=100"`
	TargetMemoryUtilizationPercentage *int32       `yaml:"targetMemoryUtilizationPercentage" validate:"omitempty,min=1,max=100"`
	Metrics                           []MetricSpec `yaml:"metrics,omitempty" validate:"dive"`
	Behavior                          *HPABehavior `yaml:"behavior,omitempty"`
}

// MetricSpec represents HPA metric specification
type MetricSpec struct {
	Type              string                         `yaml:"type" validate:"required,oneof=Resource Pods Object External ContainerResource"`
	Resource          *ResourceMetricSource          `yaml:"resource,omitempty"`
	Pods              *PodsMetricSource              `yaml:"pods,omitempty"`
	Object            *ObjectMetricSource            `yaml:"object,omitempty"`
	External          *ExternalMetricSource          `yaml:"external,omitempty"`
	ContainerResource *ContainerResourceMetricSource `yaml:"containerResource,omitempty"`
}

// ResourceMetricSource represents resource metric source
type ResourceMetricSource struct {
	Name   string       `yaml:"name" validate:"required"`
	Target MetricTarget `yaml:"target" validate:"required"`
}

// PodsMetricSource represents pods metric source
type PodsMetricSource struct {
	Metric MetricIdentifier `yaml:"metric" validate:"required"`
	Target MetricTarget     `yaml:"target" validate:"required"`
}

// ObjectMetricSource represents object metric source
type ObjectMetricSource struct {
	DescribedObject CrossVersionObjectReference `yaml:"describedObject" validate:"required"`
	Metric          MetricIdentifier            `yaml:"metric" validate:"required"`
	Target          MetricTarget                `yaml:"target" validate:"required"`
}

// ExternalMetricSource represents external metric source
type ExternalMetricSource struct {
	Metric MetricIdentifier `yaml:"metric" validate:"required"`
	Target MetricTarget     `yaml:"target" validate:"required"`
}

// ContainerResourceMetricSource represents container resource metric source
type ContainerResourceMetricSource struct {
	Name      string       `yaml:"name" validate:"required"`
	Container string       `yaml:"container" validate:"required"`
	Target    MetricTarget `yaml:"target" validate:"required"`
}

// MetricTarget represents metric target
type MetricTarget struct {
	Type               string  `yaml:"type" validate:"required,oneof=Utilization Value AverageValue"`
	Value              *string `yaml:"value,omitempty" validate:"omitempty,resource_quantity"`
	AverageValue       *string `yaml:"averageValue,omitempty" validate:"omitempty,resource_quantity"`
	AverageUtilization *int32  `yaml:"averageUtilization,omitempty" validate:"omitempty,min=1,max=100"`
}

// MetricIdentifier represents metric identifier
type MetricIdentifier struct {
	Name     string         `yaml:"name" validate:"required"`
	Selector *LabelSelector `yaml:"selector,omitempty"`
}

// CrossVersionObjectReference represents cross version object reference
type CrossVersionObjectReference struct {
	Kind       string `yaml:"kind" validate:"required"`
	Name       string `yaml:"name" validate:"required"`
	APIVersion string `yaml:"apiVersion" validate:"required"`
}

// LabelSelector represents label selector
type LabelSelector struct {
	MatchLabels      map[string]string          `yaml:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `yaml:"matchExpressions,omitempty" validate:"dive"`
}

// LabelSelectorRequirement represents label selector requirement
type LabelSelectorRequirement struct {
	Key      string   `yaml:"key" validate:"required"`
	Operator string   `yaml:"operator" validate:"required,oneof=In NotIn Exists DoesNotExist"`
	Values   []string `yaml:"values,omitempty"`
}

// HPABehavior represents HPA behavior configuration
type HPABehavior struct {
	ScaleUp   *HPAScalingRules `yaml:"scaleUp,omitempty"`
	ScaleDown *HPAScalingRules `yaml:"scaleDown,omitempty"`
}

// HPAScalingRules represents HPA scaling rules
type HPAScalingRules struct {
	StabilizationWindowSeconds *int32             `yaml:"stabilizationWindowSeconds,omitempty" validate:"omitempty,min=0"`
	SelectPolicy               string             `yaml:"selectPolicy,omitempty" validate:"omitempty,oneof=Max Min Disabled"`
	Policies                   []HPAScalingPolicy `yaml:"policies,omitempty" validate:"dive"`
}

// HPAScalingPolicy represents HPA scaling policy
type HPAScalingPolicy struct {
	Type          string `yaml:"type" validate:"required,oneof=Pods Percent"`
	Value         int32  `yaml:"value" validate:"required,min=1"`
	PeriodSeconds int32  `yaml:"periodSeconds" validate:"required,min=1"`
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
	SubPathExpr      string `yaml:"subPathExpr,omitempty"`
}

// Volume represents volume configuration
type Volume struct {
	Name                  string                       `yaml:"name" validate:"required"`
	HostPath              *HostPathVolumeSource        `yaml:"hostPath,omitempty"`
	EmptyDir              *EmptyDirVolumeSource        `yaml:"emptyDir,omitempty"`
	Secret                *SecretVolumeSource          `yaml:"secret,omitempty"`
	ConfigMap             *ConfigMapVolumeSource       `yaml:"configMap,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaimSource `yaml:"persistentVolumeClaim,omitempty"`
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

// KeyToPath represents key to path mapping
type KeyToPath struct {
	Key  string `yaml:"key" validate:"required"`
	Path string `yaml:"path" validate:"required,filepath"`
	Mode *int32 `yaml:"mode,omitempty" validate:"omitempty,min=0,max=511"`
}

// Validate validates the entire Values configuration
func (v *Values) Validate() error {
	return helmcharts.ValidateStruct(v)
}

// Validate validates the GlobalConfig
func (g *GlobalConfig) Validate() error {
	return helmcharts.ValidateStruct(g)
}

// Validate validates the AdminConfig
func (a *AdminConfig) Validate() error {
	return helmcharts.ValidateStruct(a)
}

// Validate validates the AdminAppConfig
func (a *AdminAppConfig) Validate() error {
	return helmcharts.ValidateStruct(a)
}

// Validate validates the HTTPProbeConfig
func (h *HTTPProbeConfig) Validate() error {
	return helmcharts.ValidateStruct(h)
}

// Validate validates the AdminServiceConfig
func (s *AdminServiceConfig) Validate() error {
	return helmcharts.ValidateStruct(s)
}
