package portal_controller_kubernetes

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

	// Additional specific checks
	if values.Controller.ManagerContainer.Image.Repository == "" {
		t.Error("Manager container image repository should not be empty")
	}

	if values.Controller.ManagerContainer.Image.Tag == "" {
		t.Error("Manager container image tag should not be empty")
	}

	if values.Controller.TerminationGracePeriodSeconds < 0 {
		t.Error("TerminationGracePeriodSeconds should not be negative")
	}
}

func TestCRDsConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  CRDsConfig
		wantErr bool
	}{
		{
			name: "CRDs enabled",
			config: CRDsConfig{
				Install: true,
			},
			wantErr: false,
		},
		{
			name: "CRDs disabled",
			config: CRDsConfig{
				Install: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := Values{
				CRDs: tt.config,
				Controller: ControllerConfig{
					TerminationGracePeriodSeconds: 10,
					ManagerContainer: ManagerContainerConfig{
						Image: ContainerImage{
							Repository: "test/manager",
							Tag:        "v1.0.0",
						},
						LivenessProbe: HTTPProbeConfig{
							HTTPGet: HTTPGetAction{
								Path: "/healthz",
								Port: 8081,
							},
							InitialDelaySeconds: 15,
							PeriodSeconds:       20,
						},
						ReadinessProbe: HTTPProbeConfig{
							HTTPGet: HTTPGetAction{
								Path: "/readyz",
								Port: 8081,
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       10,
						},
					},
				},
			}

			err := values.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CRDsConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControllerConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ControllerConfig
		wantErr bool
	}{
		{
			name: "valid controller config",
			config: ControllerConfig{
				TerminationGracePeriodSeconds: 10,
				ManagerContainer: ManagerContainerConfig{
					Image: ContainerImage{
						Repository: "test/manager",
						Tag:        "v1.0.0",
					},
					LivenessProbe: HTTPProbeConfig{
						HTTPGet: HTTPGetAction{
							Path: "/healthz",
							Port: 8081,
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       20,
					},
					ReadinessProbe: HTTPProbeConfig{
						HTTPGet: HTTPGetAction{
							Path: "/readyz",
							Port: 8081,
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "negative termination grace period",
			config: ControllerConfig{
				TerminationGracePeriodSeconds: -1,
				ManagerContainer: ManagerContainerConfig{
					Image: ContainerImage{
						Repository: "test/manager",
						Tag:        "v1.0.0",
					},
					LivenessProbe: HTTPProbeConfig{
						HTTPGet: HTTPGetAction{
							Path: "/healthz",
							Port: 8081,
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       20,
					},
					ReadinessProbe: HTTPProbeConfig{
						HTTPGet: HTTPGetAction{
							Path: "/readyz",
							Port: 8081,
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing manager container",
			config: ControllerConfig{
				TerminationGracePeriodSeconds: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ControllerConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManagerContainerConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ManagerContainerConfig
		wantErr bool
	}{
		{
			name: "valid manager container config",
			config: ManagerContainerConfig{
				Image: ContainerImage{
					Repository: "test/manager",
					Tag:        "v1.0.0",
					PullPolicy: "IfNotPresent",
				},
				LivenessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/healthz",
						Port: 8081,
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
				ReadinessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/readyz",
						Port: 8081,
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
			},
			wantErr: false,
		},
		{
			name: "missing image repository",
			config: ManagerContainerConfig{
				Image: ContainerImage{
					Tag: "v1.0.0",
				},
				LivenessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/healthz",
						Port: 8081,
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
				ReadinessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/readyz",
						Port: 8081,
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "missing image tag",
			config: ManagerContainerConfig{
				Image: ContainerImage{
					Repository: "test/manager",
				},
				LivenessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/healthz",
						Port: 8081,
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
				ReadinessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/readyz",
						Port: 8081,
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid image pull policy",
			config: ManagerContainerConfig{
				Image: ContainerImage{
					Repository: "test/manager",
					Tag:        "v1.0.0",
					PullPolicy: "InvalidPolicy",
				},
				LivenessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/healthz",
						Port: 8081,
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
				ReadinessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/readyz",
						Port: 8081,
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ManagerContainerConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPProbeConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		probe   HTTPProbeConfig
		wantErr bool
	}{
		{
			name: "valid HTTP probe",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8081,
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
			},
			wantErr: false,
		},
		{
			name: "valid HTTP probe with optional fields",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path:   "/healthz",
					Port:   8081,
					Host:   "127.0.0.1",
					Scheme: "HTTP",
					HTTPHeaders: []HTTPHeader{
						{Name: "Custom-Header", Value: "test"},
					},
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
				TimeoutSeconds:      5,
				SuccessThreshold:    1,
				FailureThreshold:    3,
			},
			wantErr: false,
		},
		{
			name: "missing path",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Port: 8081,
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
			},
			wantErr: true,
		},
		{
			name: "missing port",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
			},
			wantErr: true,
		},
		{
			name: "invalid port range",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 70000,
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
			},
			wantErr: true,
		},
		{
			name: "negative initial delay",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8081,
				},
				InitialDelaySeconds: -1,
				PeriodSeconds:       20,
			},
			wantErr: true,
		},
		{
			name: "invalid period seconds",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8081,
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       0,
			},
			wantErr: true,
		},
		{
			name: "invalid scheme",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path:   "/healthz",
					Port:   8081,
					Scheme: "FTP",
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
			},
			wantErr: true,
		},
		{
			name: "invalid host format",
			probe: HTTPProbeConfig{
				HTTPGet: HTTPGetAction{
					Path: "/healthz",
					Port: 8081,
					Host: "invalid_host",
				},
				InitialDelaySeconds: 15,
				PeriodSeconds:       20,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.probe.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPProbeConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeccompProfileValidation(t *testing.T) {
	tests := []struct {
		name    string
		profile SeccompProfile
		wantErr bool
	}{
		{
			name: "valid RuntimeDefault profile",
			profile: SeccompProfile{
				Type: "RuntimeDefault",
			},
			wantErr: false,
		},
		{
			name: "valid Localhost profile",
			profile: SeccompProfile{
				Type:             "Localhost",
				LocalhostProfile: "/profiles/custom.json",
			},
			wantErr: false,
		},
		{
			name: "valid Unconfined profile",
			profile: SeccompProfile{
				Type: "Unconfined",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			profile: SeccompProfile{
				Type: "InvalidType",
			},
			wantErr: true,
		},
		{
			name: "Localhost type missing profile",
			profile: SeccompProfile{
				Type: "Localhost",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ContainerSecurityContext{
				SeccompProfile: &tt.profile,
			}

			container := ManagerContainerConfig{
				Image: ContainerImage{
					Repository: "test/manager",
					Tag:        "v1.0.0",
				},
				LivenessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/healthz",
						Port: 8081,
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
				ReadinessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/readyz",
						Port: 8081,
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				SecurityContext: config,
			}

			err := container.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SeccompProfile validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainerPortValidation(t *testing.T) {
	tests := []struct {
		name    string
		port    ContainerPort
		wantErr bool
	}{
		{
			name: "valid container port",
			port: ContainerPort{
				Name:          "http",
				ContainerPort: 8080,
				Protocol:      "TCP",
			},
			wantErr: false,
		},
		{
			name: "valid container port with host port",
			port: ContainerPort{
				ContainerPort: 8080,
				HostPort:      8080,
				HostIP:        "127.0.0.1",
			},
			wantErr: false,
		},
		{
			name: "missing container port",
			port: ContainerPort{
				Name: "http",
			},
			wantErr: true,
		},
		{
			name: "invalid container port range",
			port: ContainerPort{
				ContainerPort: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid container port range (too high)",
			port: ContainerPort{
				ContainerPort: 70000,
			},
			wantErr: true,
		},
		{
			name: "invalid protocol",
			port: ContainerPort{
				ContainerPort: 8080,
				Protocol:      "INVALID",
			},
			wantErr: true,
		},
		{
			name: "invalid host IP",
			port: ContainerPort{
				ContainerPort: 8080,
				HostIP:        "invalid-ip",
			},
			wantErr: true,
		},
		{
			name: "invalid host port range",
			port: ContainerPort{
				ContainerPort: 8080,
				HostPort:      -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ManagerContainerConfig{
				Image: ContainerImage{
					Repository: "test/manager",
					Tag:        "v1.0.0",
				},
				LivenessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/healthz",
						Port: 8081,
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
				ReadinessProbe: HTTPProbeConfig{
					HTTPGet: HTTPGetAction{
						Path: "/readyz",
						Port: 8081,
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Ports: []ContainerPort{tt.port},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContainerPort validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricsConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		metrics MetricsConfig
		wantErr bool
	}{
		{
			name: "metrics disabled",
			metrics: MetricsConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "metrics enabled with valid config",
			metrics: MetricsConfig{
				Enabled: true,
				Port:    9090,
				Path:    "/metrics",
			},
			wantErr: false,
		},
		{
			name: "metrics enabled missing port",
			metrics: MetricsConfig{
				Enabled: true,
				Path:    "/metrics",
			},
			wantErr: true,
		},
		{
			name: "metrics enabled missing path",
			metrics: MetricsConfig{
				Enabled: true,
				Port:    9090,
			},
			wantErr: true,
		},
		{
			name: "metrics enabled invalid port",
			metrics: MetricsConfig{
				Enabled: true,
				Port:    0,
				Path:    "/metrics",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ControllerConfig{
				TerminationGracePeriodSeconds: 10,
				ManagerContainer: ManagerContainerConfig{
					Image: ContainerImage{
						Repository: "test/manager",
						Tag:        "v1.0.0",
					},
					LivenessProbe: HTTPProbeConfig{
						HTTPGet: HTTPGetAction{
							Path: "/healthz",
							Port: 8081,
						},
						InitialDelaySeconds: 15,
						PeriodSeconds:       20,
					},
					ReadinessProbe: HTTPProbeConfig{
						HTTPGet: HTTPGetAction{
							Path: "/readyz",
							Port: 8081,
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
				},
				Metrics: tt.metrics,
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
