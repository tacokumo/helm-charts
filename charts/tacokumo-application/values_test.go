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
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			},
			wantErr: false,
		},
		{
			name: "missing application name",
			config: MainConfig{
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			},
			wantErr: true,
		},
		{
			name: "missing image",
			config: MainConfig{
				ApplicationName: "test-app",
				ImagePullPolicy: "IfNotPresent",
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			},
			wantErr: true,
		},
		{
			name: "zero min replicas",
			config: MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
				HPA: HPAConfig{
					MinReplicas:                       0,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			},
			wantErr: true,
		},
		{
			name: "max replicas less than min replicas",
			config: MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
				HPA: HPAConfig{
					MinReplicas:                       3,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid memory utilization percentage",
			config: MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 101,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid image pull policy",
			config: MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				ImagePullPolicy: "InvalidPolicy",
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
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
				Image:           "nginx:latest",
				EnvFrom:         []EnvFromSource{tt.envFrom},
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
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
				Image:            "nginx:latest",
				ImagePullSecrets: []ImagePullSecret{tt.secret},
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
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
					Memory: "512Mi",
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
					Memory: "512Mi",
				},
				Requests: ResourceSpec{
					CPU:    "100m",
					Memory: "128Mi",
				},
			},
			wantErr: false,
		},
		{
			name: "partial limits (cpu only)",
			resources: ResourceConfig{
				Limits: ResourceSpec{
					CPU: "500m",
				},
			},
			wantErr: false,
		},
		{
			name: "partial limits (memory only)",
			resources: ResourceConfig{
				Limits: ResourceSpec{
					Memory: "512Mi",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
				Resources:       tt.resources,
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		service ServiceConfig
		wantErr bool
	}{
		{
			name: "disabled service (no ports required)",
			service: ServiceConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "enabled service with valid port",
			service: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Ports: []ServicePortConfig{
					{
						Name: "http",
						Port: 80,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "enabled service without ports",
			service: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Ports:   []ServicePortConfig{},
			},
			wantErr: true,
		},
		{
			name: "invalid service type",
			service: ServiceConfig{
				Enabled: true,
				Type:    "InvalidType",
				Ports: []ServicePortConfig{
					{Port: 80},
				},
			},
			wantErr: true,
		},
		{
			name: "NodePort type",
			service: ServiceConfig{
				Enabled: true,
				Type:    "NodePort",
				Ports: []ServicePortConfig{
					{
						Port:     80,
						NodePort: 30080,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "LoadBalancer type",
			service: ServiceConfig{
				Enabled: true,
				Type:    "LoadBalancer",
				Ports: []ServicePortConfig{
					{Port: 80},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple ports",
			service: ServiceConfig{
				Enabled: true,
				Type:    "ClusterIP",
				Ports: []ServicePortConfig{
					{Name: "http", Port: 80, TargetPort: 8080},
					{Name: "https", Port: 443, TargetPort: 8443},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				Service:         tt.service,
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePortConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		port    ServicePortConfig
		wantErr bool
	}{
		{
			name: "valid minimal port",
			port: ServicePortConfig{
				Port: 80,
			},
			wantErr: false,
		},
		{
			name: "valid full port config",
			port: ServicePortConfig{
				Name:       "http",
				Port:       80,
				TargetPort: 8080,
				Protocol:   "TCP",
			},
			wantErr: false,
		},
		{
			name: "valid UDP protocol",
			port: ServicePortConfig{
				Name:     "dns",
				Port:     53,
				Protocol: "UDP",
			},
			wantErr: false,
		},
		{
			name: "valid SCTP protocol",
			port: ServicePortConfig{
				Port:     3868,
				Protocol: "SCTP",
			},
			wantErr: false,
		},
		{
			name: "missing port",
			port: ServicePortConfig{
				Name: "http",
			},
			wantErr: true,
		},
		{
			name: "port too low",
			port: ServicePortConfig{
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "port too high",
			port: ServicePortConfig{
				Port: 65536,
			},
			wantErr: true,
		},
		{
			name: "targetPort too high",
			port: ServicePortConfig{
				Port:       80,
				TargetPort: 70000,
			},
			wantErr: true,
		},
		{
			name: "invalid protocol",
			port: ServicePortConfig{
				Port:     80,
				Protocol: "HTTP",
			},
			wantErr: true,
		},
		{
			name: "valid nodePort",
			port: ServicePortConfig{
				Port:     80,
				NodePort: 30080,
			},
			wantErr: false,
		},
		{
			name: "nodePort too low",
			port: ServicePortConfig{
				Port:     80,
				NodePort: 29999,
			},
			wantErr: true,
		},
		{
			name: "nodePort too high",
			port: ServicePortConfig{
				Port:     80,
				NodePort: 32768,
			},
			wantErr: true,
		},
		{
			name: "max valid port",
			port: ServicePortConfig{
				Port: 65535,
			},
			wantErr: false,
		},
		{
			name: "min valid port",
			port: ServicePortConfig{
				Port: 1,
			},
			wantErr: false,
		},
		{
			name: "min valid nodePort",
			port: ServicePortConfig{
				Port:     80,
				NodePort: 30000,
			},
			wantErr: false,
		},
		{
			name: "max valid nodePort",
			port: ServicePortConfig{
				Port:     80,
				NodePort: 32767,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				Service: ServiceConfig{
					Enabled: true,
					Type:    "ClusterIP",
					Ports:   []ServicePortConfig{tt.port},
				},
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServicePortConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create int pointers for test cases
func intPtr(i int) *int {
	return &i
}

func TestIngressConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		ingress IngressConfig
		wantErr bool
	}{
		{
			name: "disabled ingress",
			ingress: IngressConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "valid enabled ingress",
			ingress: IngressConfig{
				Enabled:   true,
				ClassName: "nginx",
				Annotations: map[string]string{
					"nginx.ingress.kubernetes.io/rewrite-target": "/",
				},
				Hosts: []IngressHost{
					{
						Host: "app.example.com",
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "Prefix",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "enabled ingress missing className",
			ingress: IngressConfig{
				Enabled: true,
				Hosts: []IngressHost{
					{
						Host: "app.example.com",
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "Prefix",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "enabled ingress missing hosts",
			ingress: IngressConfig{
				Enabled:   true,
				ClassName: "nginx",
			},
			wantErr: true,
		},
		{
			name: "enabled ingress with TLS",
			ingress: IngressConfig{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "app.example.com",
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "Prefix",
							},
						},
					},
				},
				TLS: []IngressTLS{
					{
						SecretName: "app-tls",
						Hosts:      []string{"app.example.com"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid host format",
			ingress: IngressConfig{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "invalid-host",
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "Prefix",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid path type",
			ingress: IngressConfig{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "app.example.com",
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "InvalidType",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "TLS missing secret name",
			ingress: IngressConfig{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "app.example.com",
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "Prefix",
							},
						},
					},
				},
				TLS: []IngressTLS{
					{
						Hosts: []string{"app.example.com"},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				Ingress:         tt.ingress,
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IngressConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPRouteConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		route   RouteConfig
		wantErr bool
	}{
		{
			name: "disabled HTTPRoute",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "valid enabled HTTPRoute",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "PathPrefix",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "enabled HTTPRoute missing parentRefs",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled:   true,
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "PathPrefix",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "enabled HTTPRoute missing hostnames",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "PathPrefix",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "enabled HTTPRoute missing rules",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid hostname format",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"invalid-hostname"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "PathPrefix",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "parentRef missing name",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "PathPrefix",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "parentRef missing namespace",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name: "default-gateway",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "PathPrefix",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid path type",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "InvalidType",
										Value: "/",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "path missing value",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type: "PathPrefix",
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Exact path type",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "Exact",
										Value: "/api/v1/health",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "RegularExpression path type",
			route: RouteConfig{
				HTTP: HTTPRouteConfig{
					Enabled: true,
					ParentRefs: []HTTPRouteParentRef{
						{
							Name:      "default-gateway",
							Namespace: "gateway-system",
						},
					},
					Hostnames: []string{"app.example.com"},
					Rules: []HTTPRouteRule{
						{
							Matches: []HTTPRouteMatch{
								{
									Path: &HTTPRoutePath{
										Type:  "RegularExpression",
										Value: "/api/.*/health",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MainConfig{
				ApplicationName: "test-app",
				Image:           "nginx:latest",
				Route:           tt.route,
				HPA: HPAConfig{
					MinReplicas:                       1,
					MaxReplicas:                       1,
					TargetMemoryUtilizationPercentage: 80,
				},
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPRouteConfig validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
