package tacokumo_portal_proxy

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

	// Additional checks specific to the parsed values
	if values.PortalProxy.BaseDomain == "" {
		t.Error("BaseDomain should not be empty")
	}

	if values.PortalProxy.ReplicaCount < 1 {
		t.Error("ReplicaCount should be at least 1")
	}

	if values.PortalProxy.Service.HTTPPort <= 0 || values.PortalProxy.Service.HTTPPort > 65535 {
		t.Error("HTTPPort should be between 1 and 65535")
	}

	if values.PortalProxy.Service.MetricsPort <= 0 || values.PortalProxy.Service.MetricsPort > 65535 {
		t.Error("MetricsPort should be between 1 and 65535")
	}
}

func TestPortalProxyConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PortalProxyConfig
		wantErr bool
	}{
		{
			name: "valid minimal config",
			config: PortalProxyConfig{
				ReplicaCount: 1,
				BaseDomain:   "example.com",
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "caddy",
					Tag:        "2.11",
					PullPolicy: "IfNotPresent",
				},
				Service: ProxyServiceConfig{
					Type:        "ClusterIP",
					HTTPPort:    80,
					MetricsPort: 2019,
				},
			},
			wantErr: false,
		},
		{
			name: "zero replica count",
			config: PortalProxyConfig{
				ReplicaCount: 0,
				BaseDomain:   "example.com",
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "caddy",
					Tag:        "2.11",
				},
				Service: ProxyServiceConfig{
					Type:        "ClusterIP",
					HTTPPort:    80,
					MetricsPort: 2019,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid base domain",
			config: PortalProxyConfig{
				ReplicaCount: 1,
				BaseDomain:   "invalid_domain",
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "caddy",
					Tag:        "2.11",
				},
				Service: ProxyServiceConfig{
					Type:        "ClusterIP",
					HTTPPort:    80,
					MetricsPort: 2019,
				},
			},
			wantErr: true,
		},
		{
			name: "missing base domain",
			config: PortalProxyConfig{
				ReplicaCount: 1,
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "caddy",
					Tag:        "2.11",
				},
				Service: ProxyServiceConfig{
					Type:        "ClusterIP",
					HTTPPort:    80,
					MetricsPort: 2019,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortalProxyConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyServiceConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		service ProxyServiceConfig
		wantErr bool
	}{
		{
			name: "valid service config",
			service: ProxyServiceConfig{
				Type:        "ClusterIP",
				HTTPPort:    80,
				MetricsPort: 2019,
			},
			wantErr: false,
		},
		{
			name: "valid LoadBalancer service",
			service: ProxyServiceConfig{
				Type:                     "LoadBalancer",
				HTTPPort:                 80,
				MetricsPort:              2019,
				LoadBalancerIP:           "192.168.1.100",
				LoadBalancerSourceRanges: []string{"10.0.0.0/8", "172.16.0.0/12"},
				ExternalTrafficPolicy:    "Local",
			},
			wantErr: false,
		},
		{
			name: "invalid service type",
			service: ProxyServiceConfig{
				Type:        "InvalidType",
				HTTPPort:    80,
				MetricsPort: 2019,
			},
			wantErr: true,
		},
		{
			name: "invalid HTTP port",
			service: ProxyServiceConfig{
				Type:        "ClusterIP",
				HTTPPort:    0,
				MetricsPort: 2019,
			},
			wantErr: true,
		},
		{
			name: "invalid metrics port",
			service: ProxyServiceConfig{
				Type:        "ClusterIP",
				HTTPPort:    80,
				MetricsPort: 70000,
			},
			wantErr: true,
		},
		{
			name: "invalid load balancer IP",
			service: ProxyServiceConfig{
				Type:           "LoadBalancer",
				HTTPPort:       80,
				MetricsPort:    2019,
				LoadBalancerIP: "invalid-ip",
			},
			wantErr: true,
		},
		{
			name: "invalid external traffic policy",
			service: ProxyServiceConfig{
				Type:                  "LoadBalancer",
				HTTPPort:              80,
				MetricsPort:           2019,
				ExternalTrafficPolicy: "Invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProxyServiceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePortValidation(t *testing.T) {
	tests := []struct {
		name    string
		port    ServicePort
		wantErr bool
	}{
		{
			name: "valid service port",
			port: ServicePort{
				Name: "http",
				Port: 8080,
			},
			wantErr: false,
		},
		{
			name: "valid service port with target port",
			port: ServicePort{
				Name:       "http",
				Port:       8080,
				TargetPort: 80,
				Protocol:   "TCP",
			},
			wantErr: false,
		},
		{
			name: "valid NodePort",
			port: ServicePort{
				Name:     "http",
				Port:     8080,
				NodePort: 30080,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			port: ServicePort{
				Port: 8080,
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			port: ServicePort{
				Name: "http",
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid target port",
			port: ServicePort{
				Name:       "http",
				Port:       8080,
				TargetPort: 70000,
			},
			wantErr: true,
		},
		{
			name: "invalid NodePort range",
			port: ServicePort{
				Name:     "http",
				Port:     8080,
				NodePort: 80,
			},
			wantErr: true,
		},
		{
			name: "invalid protocol",
			port: ServicePort{
				Name:     "http",
				Port:     8080,
				Protocol: "INVALID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ProxyServiceConfig{
				Type:        "ClusterIP",
				HTTPPort:    80,
				MetricsPort: 2019,
				ExtraPorts:  []ServicePort{tt.port},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServicePort validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvVarValidation(t *testing.T) {
	tests := []struct {
		name    string
		envVar  EnvVar
		wantErr bool
	}{
		{
			name: "valid environment variable with value",
			envVar: EnvVar{
				Name:  "MY_VAR",
				Value: "my-value",
			},
			wantErr: false,
		},
		{
			name: "valid environment variable with field reference",
			envVar: EnvVar{
				Name: "NODE_NAME",
				ValueFrom: &EnvVarSource{
					FieldRef: &ObjectFieldSelector{
						FieldPath: "spec.nodeName",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid environment variable with ConfigMap reference",
			envVar: EnvVar{
				Name: "DB_HOST",
				ValueFrom: &EnvVarSource{
					ConfigMapKeyRef: &ConfigMapKeySelector{
						Name: "app-config",
						Key:  "database.host",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid environment variable with Secret reference",
			envVar: EnvVar{
				Name: "DB_PASSWORD",
				ValueFrom: &EnvVarSource{
					SecretKeyRef: &SecretKeySelector{
						Name: "app-secrets",
						Key:  "database.password",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			envVar: EnvVar{
				Value: "my-value",
			},
			wantErr: true,
		},
		{
			name: "field reference missing field path",
			envVar: EnvVar{
				Name: "NODE_NAME",
				ValueFrom: &EnvVarSource{
					FieldRef: &ObjectFieldSelector{},
				},
			},
			wantErr: true,
		},
		{
			name: "ConfigMap reference missing name",
			envVar: EnvVar{
				Name: "DB_HOST",
				ValueFrom: &EnvVarSource{
					ConfigMapKeyRef: &ConfigMapKeySelector{
						Key: "database.host",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ConfigMap reference missing key",
			envVar: EnvVar{
				Name: "DB_HOST",
				ValueFrom: &EnvVarSource{
					ConfigMapKeyRef: &ConfigMapKeySelector{
						Name: "app-config",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Secret reference missing name",
			envVar: EnvVar{
				Name: "DB_PASSWORD",
				ValueFrom: &EnvVarSource{
					SecretKeyRef: &SecretKeySelector{
						Key: "database.password",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Secret reference missing key",
			envVar: EnvVar{
				Name: "DB_PASSWORD",
				ValueFrom: &EnvVarSource{
					SecretKeyRef: &SecretKeySelector{
						Name: "app-secrets",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := PortalProxyConfig{
				ReplicaCount: 1,
				BaseDomain:   "example.com",
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "caddy",
					Tag:        "2.11",
				},
				Service: ProxyServiceConfig{
					Type:        "ClusterIP",
					HTTPPort:    80,
					MetricsPort: 2019,
				},
				Env: []EnvVar{tt.envVar},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvVar validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVolumeMountValidation(t *testing.T) {
	tests := []struct {
		name        string
		volumeMount VolumeMount
		wantErr     bool
	}{
		{
			name: "valid volume mount",
			volumeMount: VolumeMount{
				Name:      "config",
				MountPath: "/etc/config",
			},
			wantErr: false,
		},
		{
			name: "valid volume mount with sub path",
			volumeMount: VolumeMount{
				Name:      "config",
				MountPath: "/etc/config",
				SubPath:   "app.conf",
				ReadOnly:  true,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			volumeMount: VolumeMount{
				MountPath: "/etc/config",
			},
			wantErr: true,
		},
		{
			name: "missing mount path",
			volumeMount: VolumeMount{
				Name: "config",
			},
			wantErr: true,
		},
		{
			name: "invalid mount propagation",
			volumeMount: VolumeMount{
				Name:             "config",
				MountPath:        "/etc/config",
				MountPropagation: "Invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := PortalProxyConfig{
				ReplicaCount: 1,
				BaseDomain:   "example.com",
				Image: struct {
					Repository string `yaml:"repository" validate:"required"`
					Tag        string `yaml:"tag" validate:"required"`
					PullPolicy string `yaml:"pullPolicy" validate:"omitempty,oneof=Always IfNotPresent Never"`
				}{
					Repository: "caddy",
					Tag:        "2.11",
				},
				Service: ProxyServiceConfig{
					Type:        "ClusterIP",
					HTTPPort:    80,
					MetricsPort: 2019,
				},
				VolumeMounts: []VolumeMount{tt.volumeMount},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("VolumeMount validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
