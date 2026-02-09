package tacokumo_portal

import (
	"os"
	"path/filepath"
	"testing"

	helmcharts "github.com/tacokumo/helm-charts"
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
}

func TestAPIConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  APIConfig
		wantErr bool
	}{
		{
			name: "valid minimal config",
			config: APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
					PullPolicy: "IfNotPresent",
				},
			},
			wantErr: false,
		},
		{
			name: "missing portal name",
			config: APIConfig{
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
					PullPolicy: "IfNotPresent",
				},
			},
			wantErr: true,
		},
		{
			name: "missing image",
			config: APIConfig{
				PortalName: "test-portal",
			},
			wantErr: true,
		},
		{
			name: "valid with log level debug",
			config: APIConfig{
				PortalName: "test-portal",
				LogLevel:   "debug",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
			},
			wantErr: false,
		},
		{
			name: "valid with log level info",
			config: APIConfig{
				PortalName: "test-portal",
				LogLevel:   "info",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
			},
			wantErr: false,
		},
		{
			name: "valid with log level warn",
			config: APIConfig{
				PortalName: "test-portal",
				LogLevel:   "warn",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
			},
			wantErr: false,
		},
		{
			name: "valid with log level error",
			config: APIConfig{
				PortalName: "test-portal",
				LogLevel:   "error",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: APIConfig{
				PortalName: "test-portal",
				LogLevel:   "invalid",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid image pull policy",
			config: APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
					PullPolicy: "InvalidPolicy",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("APIConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHPAConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  HPAConfig
		wantErr bool
	}{
		{
			name: "valid HPA config",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       1,
				MaxReplicas:                       3,
				TargetMemoryUtilizationPercentage: 80,
			},
			wantErr: false,
		},
		{
			name: "disabled HPA (no validation required)",
			config: HPAConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "zero min replicas",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       0,
				MaxReplicas:                       3,
				TargetMemoryUtilizationPercentage: 80,
			},
			wantErr: true,
		},
		{
			name: "max replicas less than min replicas",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       5,
				MaxReplicas:                       2,
				TargetMemoryUtilizationPercentage: 80,
			},
			wantErr: true,
		},
		{
			name: "memory utilization percentage too low",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       1,
				MaxReplicas:                       3,
				TargetMemoryUtilizationPercentage: 0,
			},
			wantErr: true,
		},
		{
			name: "memory utilization percentage too high",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       1,
				MaxReplicas:                       3,
				TargetMemoryUtilizationPercentage: 101,
			},
			wantErr: true,
		},
		{
			name: "valid with equal min and max replicas",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       2,
				MaxReplicas:                       2,
				TargetMemoryUtilizationPercentage: 80,
			},
			wantErr: false,
		},
		{
			name: "valid with 100 percent memory utilization",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       1,
				MaxReplicas:                       3,
				TargetMemoryUtilizationPercentage: 100,
			},
			wantErr: false,
		},
		{
			name: "valid with 1 percent memory utilization",
			config: HPAConfig{
				Enabled:                           true,
				MinReplicas:                       1,
				MaxReplicas:                       3,
				TargetMemoryUtilizationPercentage: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HPAConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ServiceConfig
		wantErr bool
	}{
		{
			name: "valid ClusterIP service",
			config: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Port:    1323,
			},
			wantErr: false,
		},
		{
			name: "valid NodePort service",
			config: ServiceConfig{
				Enabled: true,
				Type:    "NodePort",
				Port:    1323,
			},
			wantErr: false,
		},
		{
			name: "valid LoadBalancer service",
			config: ServiceConfig{
				Enabled: true,
				Type:    "LoadBalancer",
				Port:    1323,
			},
			wantErr: false,
		},
		{
			name: "disabled service (no validation required)",
			config: ServiceConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "enabled service missing port",
			config: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Port:    0,
			},
			wantErr: true,
		},
		{
			name: "invalid service type",
			config: ServiceConfig{
				Enabled: true,
				Type:    "InvalidType",
				Port:    1323,
			},
			wantErr: true,
		},
		{
			name: "port too high",
			config: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Port:    65536,
			},
			wantErr: true,
		},
		{
			name: "valid max port",
			config: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Port:    65535,
			},
			wantErr: false,
		},
		{
			name: "valid min port",
			config: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Port:    1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProbeConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		probe   ProbeConfig
		wantErr bool
	}{
		{
			name:    "empty probe (disabled)",
			probe:   ProbeConfig{},
			wantErr: false,
		},
		{
			name: "valid HTTP probe",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Path: "/health/liveness",
					Port: 1323,
				},
				InitialDelaySeconds: intPtr(15),
				PeriodSeconds:       intPtr(20),
			},
			wantErr: false,
		},
		{
			name: "valid TCP probe",
			probe: ProbeConfig{
				TCPSocket: &TCPSocketAction{
					Port: 1323,
				},
			},
			wantErr: false,
		},
		{
			name: "valid exec probe",
			probe: ProbeConfig{
				Exec: &ExecAction{
					Command: []string{"/bin/sh", "-c", "test -f /healthy"},
				},
			},
			wantErr: false,
		},
		{
			name: "HTTP probe missing path",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Port: 1323,
				},
			},
			wantErr: true,
		},
		{
			name: "HTTP probe missing port",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Path: "/health",
				},
			},
			wantErr: true,
		},
		{
			name: "HTTP probe invalid port",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Path: "/health",
					Port: 70000,
				},
			},
			wantErr: true,
		},
		{
			name: "TCP probe missing port",
			probe: ProbeConfig{
				TCPSocket: &TCPSocketAction{},
			},
			wantErr: true,
		},
		{
			name: "exec probe empty command",
			probe: ProbeConfig{
				Exec: &ExecAction{
					Command: []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid period seconds",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Path: "/health",
					Port: 1323,
				},
				PeriodSeconds: intPtr(0),
			},
			wantErr: true,
		},
		{
			name: "valid HTTP probe with HTTPS scheme",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Path:   "/health",
					Port:   1323,
					Scheme: "HTTPS",
				},
			},
			wantErr: false,
		},
		{
			name: "HTTP probe invalid scheme",
			probe: ProbeConfig{
				HTTPGet: &HTTPGetAction{
					Path:   "/health",
					Port:   1323,
					Scheme: "INVALID",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.probe.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProbeConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvVarValidation(t *testing.T) {
	tests := []struct {
		name    string
		env     EnvVar
		wantErr bool
	}{
		{
			name: "valid env with value",
			env: EnvVar{
				Name:  "MY_VAR",
				Value: "my-value",
			},
			wantErr: false,
		},
		{
			name: "valid env with configMapKeyRef",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					ConfigMapKeyRef: &ConfigMapKeySelector{
						Name: "my-configmap",
						Key:  "my-key",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid env with secretKeyRef",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					SecretKeyRef: &SecretKeySelector{
						Name: "my-secret",
						Key:  "my-key",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid env with fieldRef",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					FieldRef: &ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "missing name",
			env:     EnvVar{},
			wantErr: true,
		},
		{
			name: "configMapKeyRef missing name",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					ConfigMapKeyRef: &ConfigMapKeySelector{
						Key: "my-key",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "configMapKeyRef missing key",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					ConfigMapKeyRef: &ConfigMapKeySelector{
						Name: "my-configmap",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "secretKeyRef missing name",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					SecretKeyRef: &SecretKeySelector{
						Key: "my-key",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "secretKeyRef missing key",
			env: EnvVar{
				Name: "MY_VAR",
				ValueFrom: &EnvVarSource{
					SecretKeyRef: &SecretKeySelector{
						Name: "my-secret",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
				Env: []EnvVar{tt.env},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvVar validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvFromSourceValidation(t *testing.T) {
	tests := []struct {
		name    string
		envFrom EnvFromSource
		wantErr bool
	}{
		{
			name: "valid ConfigMap reference",
			envFrom: EnvFromSource{
				ConfigMapRef: &ConfigMapEnvSource{
					Name: "app-config",
				},
			},
			wantErr: false,
		},
		{
			name: "valid Secret reference",
			envFrom: EnvFromSource{
				SecretRef: &SecretEnvSource{
					Name: "app-secrets",
				},
			},
			wantErr: false,
		},
		{
			name: "valid ConfigMap reference with prefix",
			envFrom: EnvFromSource{
				Prefix: "CONFIG_",
				ConfigMapRef: &ConfigMapEnvSource{
					Name: "app-config",
				},
			},
			wantErr: false,
		},
		{
			name: "ConfigMap reference missing name",
			envFrom: EnvFromSource{
				ConfigMapRef: &ConfigMapEnvSource{},
			},
			wantErr: true,
		},
		{
			name: "Secret reference missing name",
			envFrom: EnvFromSource{
				SecretRef: &SecretEnvSource{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
				EnvFrom: []EnvFromSource{tt.envFrom},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvFromSource validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImagePullSecretValidation(t *testing.T) {
	tests := []struct {
		name    string
		secret  ImagePullSecret
		wantErr bool
	}{
		{
			name: "valid secret",
			secret: ImagePullSecret{
				Name: "regcred",
			},
			wantErr: false,
		},
		{
			name:    "missing name",
			secret:  ImagePullSecret{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
				ImagePullSecrets: []ImagePullSecret{tt.secret},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ImagePullSecret validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		resources ResourceConfig
		wantErr   bool
	}{
		{
			name:      "empty resources (default)",
			resources: ResourceConfig{},
			wantErr:   false,
		},
		{
			name: "limits only",
			resources: ResourceConfig{
				Limits: ResourceSpec{
					CPU:    "500m",
					Memory: "256Mi",
				},
			},
			wantErr: false,
		},
		{
			name: "requests only",
			resources: ResourceConfig{
				Requests: ResourceSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
			},
			wantErr: false,
		},
		{
			name: "both limits and requests",
			resources: ResourceConfig{
				Limits: ResourceSpec{
					CPU:    "500m",
					Memory: "256Mi",
				},
				Requests: ResourceSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
				Resources: tt.resources,
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecurityContextValidation(t *testing.T) {
	tests := []struct {
		name            string
		securityContext SecurityContext
		wantErr         bool
	}{
		{
			name:            "empty security context",
			securityContext: SecurityContext{},
			wantErr:         false,
		},
		{
			name: "valid with runAsUser",
			securityContext: SecurityContext{
				RunAsUser: int64Ptr(65532),
			},
			wantErr: false,
		},
		{
			name: "valid with runAsNonRoot",
			securityContext: SecurityContext{
				RunAsNonRoot: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "valid with readOnlyRootFilesystem",
			securityContext: SecurityContext{
				ReadOnlyRootFilesystem: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "valid with allowPrivilegeEscalation false",
			securityContext: SecurityContext{
				AllowPrivilegeEscalation: boolPtr(false),
			},
			wantErr: false,
		},
		{
			name: "valid with capabilities drop ALL",
			securityContext: SecurityContext{
				Capabilities: &Capabilities{
					Drop: []string{"ALL"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid seccomp profile RuntimeDefault",
			securityContext: SecurityContext{
				SeccompProfile: &SeccompProfile{
					Type: "RuntimeDefault",
				},
			},
			wantErr: false,
		},
		{
			name: "valid seccomp profile Unconfined",
			securityContext: SecurityContext{
				SeccompProfile: &SeccompProfile{
					Type: "Unconfined",
				},
			},
			wantErr: false,
		},
		{
			name: "valid seccomp profile Localhost with profile",
			securityContext: SecurityContext{
				SeccompProfile: &SeccompProfile{
					Type:             "Localhost",
					LocalhostProfile: "my-profile.json",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid seccomp profile type",
			securityContext: SecurityContext{
				SeccompProfile: &SeccompProfile{
					Type: "InvalidType",
				},
			},
			wantErr: true,
		},
		{
			name: "seccomp Localhost missing profile",
			securityContext: SecurityContext{
				SeccompProfile: &SeccompProfile{
					Type: "Localhost",
				},
			},
			wantErr: true,
		},
		{
			name: "negative runAsUser",
			securityContext: SecurityContext{
				RunAsUser: int64Ptr(-1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
				SecurityContext: tt.securityContext,
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SecurityContext validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceAccountConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ServiceAccountConfig
		wantErr bool
	}{
		{
			name: "valid with create true and name",
			config: ServiceAccountConfig{
				Create: true,
				Name:   "portal-api",
			},
			wantErr: false,
		},
		{
			name: "valid with create false and no name",
			config: ServiceAccountConfig{
				Create: false,
			},
			wantErr: false,
		},
		{
			name: "valid with annotations",
			config: ServiceAccountConfig{
				Create: true,
				Name:   "portal-api",
				Annotations: map[string]string{
					"eks.amazonaws.com/role-arn": "arn:aws:iam::123456789012:role/my-role",
				},
			},
			wantErr: false,
		},
		{
			name: "create true missing name",
			config: ServiceAccountConfig{
				Create: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := APIConfig{
				PortalName: "test-portal",
				Image: helmcharts.Image{
					Repository: "ghcr.io/tacokumo/portal-api",
					Tag:        "latest",
				},
				ServiceAccount: tt.config,
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceAccountConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create int pointers for test cases
func intPtr(i int) *int {
	return &i
}

// Helper function to create int64 pointers for test cases
func int64Ptr(i int64) *int64 {
	return &i
}

// Helper function to create bool pointers for test cases
func boolPtr(b bool) *bool {
	return &b
}
