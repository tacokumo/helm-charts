package helmcharts

import (
	"github.com/go-playground/validator/v10"
)

// Image represents container image configuration
type Image struct {
	Repository string `yaml:"repository" validate:"required"`
	Tag        string `yaml:"tag" validate:"required"`
	PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
}

// Resources represents Kubernetes resource requests and limits
type Resources struct {
	Requests ResourceRequests `yaml:"requests,omitempty"`
	Limits   ResourceLimits   `yaml:"limits,omitempty"`
}

// ResourceRequests represents resource requests
type ResourceRequests struct {
	CPU    string `yaml:"cpu,omitempty" validate:"omitempty,resource_quantity"`
	Memory string `yaml:"memory,omitempty" validate:"omitempty,resource_quantity"`
}

// ResourceLimits represents resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu,omitempty" validate:"omitempty,resource_quantity"`
	Memory string `yaml:"memory,omitempty" validate:"omitempty,resource_quantity"`
}

// HTTPProbe represents HTTP health check configuration
type HTTPProbe struct {
	Enabled             bool   `yaml:"enabled"`
	Path                string `yaml:"path" validate:"required_if=Enabled true"`
	Port                int    `yaml:"port" validate:"required_if=Enabled true,omitempty,min=1,max=65535"`
	InitialDelaySeconds int    `yaml:"initialDelaySeconds" validate:"omitempty,min=0"`
	PeriodSeconds       int    `yaml:"periodSeconds" validate:"required_if=Enabled true,omitempty,min=1"`
	TimeoutSeconds      int    `yaml:"timeoutSeconds" validate:"required_if=Enabled true,omitempty,min=1"`
	SuccessThreshold    int    `yaml:"successThreshold" validate:"omitempty,min=1"`
	FailureThreshold    int    `yaml:"failureThreshold" validate:"required_if=Enabled true,omitempty,min=1"`
}

// SecurityContext represents pod security context
type SecurityContext struct {
	RunAsUser                *int64        `yaml:"runAsUser,omitempty" validate:"omitempty,min=0"`
	RunAsGroup               *int64        `yaml:"runAsGroup,omitempty" validate:"omitempty,min=0"`
	RunAsNonRoot             *bool         `yaml:"runAsNonRoot,omitempty"`
	ReadOnlyRootFilesystem   *bool         `yaml:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation *bool         `yaml:"allowPrivilegeEscalation,omitempty"`
	Capabilities             *Capabilities `yaml:"capabilities,omitempty"`
}

// Capabilities represents security capabilities
type Capabilities struct {
	Add  []string `yaml:"add,omitempty"`
	Drop []string `yaml:"drop,omitempty"`
}

// Service represents Kubernetes service configuration
type Service struct {
	Type        string            `yaml:"type" validate:"oneof=ClusterIP NodePort LoadBalancer ExternalName"`
	Port        int               `yaml:"port" validate:"min=1,max=65535"`
	TargetPort  int               `yaml:"targetPort,omitempty" validate:"omitempty,min=1,max=65535"`
	NodePort    int               `yaml:"nodePort,omitempty" validate:"omitempty,min=30000,max=32767"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Ingress represents Kubernetes ingress configuration
type Ingress struct {
	Enabled     bool              `yaml:"enabled"`
	ClassName   string            `yaml:"className,omitempty" validate:"required_if=Enabled true"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Hosts       []IngressHost     `yaml:"hosts,omitempty" validate:"required_if=Enabled true,dive"`
	TLS         []IngressTLS      `yaml:"tls,omitempty" validate:"dive"`
}

// IngressHost represents ingress host configuration
type IngressHost struct {
	Host  string        `yaml:"host" validate:"required,fqdn"`
	Paths []IngressPath `yaml:"paths" validate:"required,dive"`
}

// IngressPath represents ingress path configuration
type IngressPath struct {
	Path     string `yaml:"path" validate:"required"`
	PathType string `yaml:"pathType" validate:"required,oneof=Exact Prefix ImplementationSpecific"`
}

// IngressTLS represents ingress TLS configuration
type IngressTLS struct {
	SecretName string   `yaml:"secretName" validate:"required"`
	Hosts      []string `yaml:"hosts" validate:"required,dive,fqdn"`
}

// PodDisruptionBudget represents PDB configuration
type PodDisruptionBudget struct {
	Enabled        bool `yaml:"enabled"`
	MinAvailable   *int `yaml:"minAvailable,omitempty" validate:"omitempty,min=1"`
	MaxUnavailable *int `yaml:"maxUnavailable,omitempty" validate:"omitempty,min=1"`
}

// Affinity represents pod affinity configuration
type Affinity struct {
	NodeAffinity    *NodeAffinity    `yaml:"nodeAffinity,omitempty"`
	PodAffinity     *PodAffinity     `yaml:"podAffinity,omitempty"`
	PodAntiAffinity *PodAntiAffinity `yaml:"podAntiAffinity,omitempty"`
}

// NodeAffinity represents node affinity rules
type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  *NodeSelector             `yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `yaml:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// NodeSelector represents node selector terms
type NodeSelector struct {
	NodeSelectorTerms []NodeSelectorTerm `yaml:"nodeSelectorTerms" validate:"required,dive"`
}

// NodeSelectorTerm represents node selector term
type NodeSelectorTerm struct {
	MatchExpressions []NodeSelectorRequirement `yaml:"matchExpressions,omitempty" validate:"dive"`
	MatchFields      []NodeSelectorRequirement `yaml:"matchFields,omitempty" validate:"dive"`
}

// NodeSelectorRequirement represents node selector requirement
type NodeSelectorRequirement struct {
	Key      string   `yaml:"key" validate:"required"`
	Operator string   `yaml:"operator" validate:"required,oneof=In NotIn Exists DoesNotExist Gt Lt"`
	Values   []string `yaml:"values,omitempty"`
}

// PreferredSchedulingTerm represents preferred scheduling term
type PreferredSchedulingTerm struct {
	Weight     int32            `yaml:"weight" validate:"required,min=1,max=100"`
	Preference NodeSelectorTerm `yaml:"preference" validate:"required"`
}

// PodAffinity represents pod affinity rules
type PodAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTerm         `yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" validate:"dive"`
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `yaml:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" validate:"dive"`
}

// PodAntiAffinity represents pod anti-affinity rules
type PodAntiAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTerm         `yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" validate:"dive"`
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `yaml:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" validate:"dive"`
}

// PodAffinityTerm represents pod affinity term
type PodAffinityTerm struct {
	LabelSelector *LabelSelector `yaml:"labelSelector,omitempty"`
	Namespaces    []string       `yaml:"namespaces,omitempty"`
	TopologyKey   string         `yaml:"topologyKey" validate:"required"`
}

// WeightedPodAffinityTerm represents weighted pod affinity term
type WeightedPodAffinityTerm struct {
	Weight          int32           `yaml:"weight" validate:"required,min=1,max=100"`
	PodAffinityTerm PodAffinityTerm `yaml:"podAffinityTerm" validate:"required"`
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

// Tolerations represents pod tolerations
type Tolerations []Toleration

// Toleration represents pod toleration
type Toleration struct {
	Key               string `yaml:"key,omitempty"`
	Operator          string `yaml:"operator,omitempty" validate:"omitempty,oneof=Exists Equal"`
	Value             string `yaml:"value,omitempty"`
	Effect            string `yaml:"effect,omitempty" validate:"omitempty,oneof=NoSchedule PreferNoSchedule NoExecute"`
	TolerationSeconds *int64 `yaml:"tolerationSeconds,omitempty" validate:"omitempty,min=0"`
}

// Validate validates the common structs using the validator
func (i *Image) Validate() error {
	return validator.New().Struct(i)
}

func (r *Resources) Validate() error {
	return validator.New().Struct(r)
}

func (p *HTTPProbe) Validate() error {
	return validator.New().Struct(p)
}

func (s *SecurityContext) Validate() error {
	return validator.New().Struct(s)
}

func (s *Service) Validate() error {
	return validator.New().Struct(s)
}

func (i *Ingress) Validate() error {
	return validator.New().Struct(i)
}
