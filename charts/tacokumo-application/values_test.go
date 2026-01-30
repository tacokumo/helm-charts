package tacokumo_application

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
}

func TestMainConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  MainConfig
		wantErr bool
	}{
		{
			name: "valid minimal config",
			config: MainConfig{
				ApplicationName: "test-app",
				ReplicaCount:    1,
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
			},
			wantErr: false,
		},
		{
			name: "missing application name",
			config: MainConfig{
				ReplicaCount:    1,
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
			},
			wantErr: true,
		},
		{
			name: "missing image",
			config: MainConfig{
				ApplicationName: "test-app",
				ReplicaCount:    1,
				ImagePullPolicy: "IfNotPresent",
			},
			wantErr: true,
		},
		{
			name: "zero replica count",
			config: MainConfig{
				ApplicationName: "test-app",
				ReplicaCount:    0,
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
			},
			wantErr: true,
		},
		{
			name: "invalid image pull policy",
			config: MainConfig{
				ApplicationName: "test-app",
				ReplicaCount:    1,
				Image:           "nginx:latest",
				ImagePullPolicy: "InvalidPolicy",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MainConfig validation error = %v, wantErr %v", err, tt.wantErr)
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
					Path: "/health",
					Port: 8080,
				},
				InitialDelaySeconds: intPtr(10),
				PeriodSeconds:       intPtr(10),
				TimeoutSeconds:      intPtr(5),
				SuccessThreshold:    intPtr(1),
				FailureThreshold:    intPtr(3),
			},
			wantErr: false,
		},
		{
			name: "valid TCP probe",
			probe: ProbeConfig{
				TCPSocket: &TCPSocketAction{
					Port: 8080,
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
					Port: 8080,
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
			name: "TCP probe invalid port",
			probe: ProbeConfig{
				TCPSocket: &TCPSocketAction{
					Port: 0,
				},
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
					Port: 8080,
				},
				PeriodSeconds: intPtr(0),
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
			config := MainConfig{
				ApplicationName: "test-app",
				ReplicaCount:    1,
				Image:           "nginx:latest",
				EnvFrom:         []EnvFromSource{tt.envFrom},
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
			config := MainConfig{
				ApplicationName:  "test-app",
				ReplicaCount:     1,
				Image:            "nginx:latest",
				ImagePullSecrets: []ImagePullSecret{tt.secret},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ImagePullSecret validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create int pointers for test cases
func intPtr(i int) *int {
	return &i
}
