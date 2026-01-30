package helmcharts

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestImageValidation(t *testing.T) {
	v, err := GetValidatorWithCustomValidations()
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}

	tests := []struct {
		name    string
		image   Image
		wantErr bool
	}{
		{
			name: "valid image",
			image: Image{
				Repository: "nginx",
				Tag:        "1.21",
				PullPolicy: "IfNotPresent",
			},
			wantErr: false,
		},
		{
			name: "missing repository",
			image: Image{
				Tag:        "1.21",
				PullPolicy: "IfNotPresent",
			},
			wantErr: true,
		},
		{
			name: "missing tag",
			image: Image{
				Repository: "nginx",
				PullPolicy: "IfNotPresent",
			},
			wantErr: true,
		},
		{
			name: "invalid pull policy",
			image: Image{
				Repository: "nginx",
				Tag:        "1.21",
				PullPolicy: "InvalidPolicy",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourcesValidation(t *testing.T) {
	v, err := GetValidatorWithCustomValidations()
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}

	tests := []struct {
		name      string
		resources Resources
		wantErr   bool
	}{
		{
			name: "valid resources",
			resources: Resources{
				Requests: ResourceRequests{
					CPU:    "100m",
					Memory: "128Mi",
				},
				Limits: ResourceLimits{
					CPU:    "500m",
					Memory: "512Mi",
				},
			},
			wantErr: false,
		},
		{
			name: "valid decimal CPU",
			resources: Resources{
				Requests: ResourceRequests{
					CPU:    "0.1",
					Memory: "128Mi",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid CPU format",
			resources: Resources{
				Requests: ResourceRequests{
					CPU:    "invalid",
					Memory: "128Mi",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid memory format",
			resources: Resources{
				Requests: ResourceRequests{
					CPU:    "100m",
					Memory: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resources validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPProbeValidation(t *testing.T) {
	v, err := GetValidatorWithCustomValidations()
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}

	tests := []struct {
		name    string
		probe   HTTPProbe
		wantErr bool
	}{
		{
			name: "disabled probe",
			probe: HTTPProbe{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "valid enabled probe",
			probe: HTTPProbe{
				Enabled:             true,
				Path:                "/health",
				Port:                8080,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				TimeoutSeconds:      5,
				SuccessThreshold:    1,
				FailureThreshold:    3,
			},
			wantErr: false,
		},
		{
			name: "enabled probe missing path",
			probe: HTTPProbe{
				Enabled: true,
				Port:    8080,
			},
			wantErr: true,
		},
		{
			name: "enabled probe missing port",
			probe: HTTPProbe{
				Enabled: true,
				Path:    "/health",
			},
			wantErr: true,
		},
		{
			name: "invalid port range",
			probe: HTTPProbe{
				Enabled: true,
				Path:    "/health",
				Port:    70000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.probe)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPProbe validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceValidation(t *testing.T) {
	v, err := GetValidatorWithCustomValidations()
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}

	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "valid ClusterIP service",
			service: Service{
				Type: "ClusterIP",
				Port: 80,
			},
			wantErr: false,
		},
		{
			name: "valid NodePort service",
			service: Service{
				Type:     "NodePort",
				Port:     80,
				NodePort: 30080,
			},
			wantErr: false,
		},
		{
			name: "invalid service type",
			service: Service{
				Type: "InvalidType",
				Port: 80,
			},
			wantErr: true,
		},
		{
			name: "invalid port range",
			service: Service{
				Type: "ClusterIP",
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid NodePort range",
			service: Service{
				Type:     "NodePort",
				Port:     80,
				NodePort: 80,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIngressValidation(t *testing.T) {
	v, err := GetValidatorWithCustomValidations()
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}

	tests := []struct {
		name    string
		ingress Ingress
		wantErr bool
	}{
		{
			name: "disabled ingress",
			ingress: Ingress{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "valid enabled ingress",
			ingress: Ingress{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "example.com",
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
			ingress: Ingress{
				Enabled: true,
				Hosts: []IngressHost{
					{
						Host: "example.com",
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
			ingress: Ingress{
				Enabled:   true,
				ClassName: "nginx",
			},
			wantErr: true,
		},
		{
			name: "invalid host",
			ingress: Ingress{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "invalid_host",
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
			ingress: Ingress{
				Enabled:   true,
				ClassName: "nginx",
				Hosts: []IngressHost{
					{
						Host: "example.com",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.ingress)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ingress validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomValidators(t *testing.T) {
	v := validator.New()
	err := RegisterCustomValidators(v)
	if err != nil {
		t.Fatalf("Failed to register custom validators: %v", err)
	}

	// Test resource quantity validation
	resourceTests := []struct {
		name     string
		quantity string
		wantErr  bool
	}{
		{"valid CPU millicores", "100m", false},
		{"valid CPU cores", "0.5", false},
		{"valid memory Mi", "256Mi", false},
		{"valid memory Gi", "1Gi", false},
		{"invalid format", "invalid", true},
		{"empty string", "", false}, // Should be valid for omitempty
	}

	type TestResource struct {
		Quantity string `validate:"resource_quantity"`
	}

	for _, tt := range resourceTests {
		t.Run(tt.name, func(t *testing.T) {
			test := TestResource{Quantity: tt.quantity}
			err := v.Struct(test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource quantity validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Test duration validation
	durationTests := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{"valid duration 1h", "1h", false},
		{"valid duration 30m", "30m", false},
		{"valid duration 45s", "45s", false},
		{"valid duration 1h30m", "1h30m", false},
		{"invalid format", "invalid", true},
		{"empty string", "", false}, // Should be valid for omitempty
	}

	type TestDuration struct {
		Duration string `validate:"duration"`
	}

	for _, tt := range durationTests {
		t.Run(tt.name, func(t *testing.T) {
			test := TestDuration{Duration: tt.duration}
			err := v.Struct(test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Test filepath validation
	filepathTests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"valid absolute path", "/etc/ssl/certs/ca.pem", false},
		{"valid relative path", "./config.yaml", false},
		{"valid Windows path", "C:\\Program Files\\app\\config.yaml", false},
		{"empty string", "", false}, // Should be valid for omitempty
	}

	type TestFilepath struct {
		Filepath string `validate:"filepath"`
	}

	for _, tt := range filepathTests {
		t.Run(tt.name, func(t *testing.T) {
			test := TestFilepath{Filepath: tt.filepath}
			err := v.Struct(test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filepath validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
