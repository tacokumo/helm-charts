package helmcharts

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers all custom validators with the given validator instance
func RegisterCustomValidators(v *validator.Validate) error {
	// Resource quantity validator (CPU: 100m, 0.1, 1, Memory: 128Mi, 1Gi)
	if err := v.RegisterValidation("resource_quantity", validateResourceQuantity); err != nil {
		return err
	}

	// Duration validator (Go duration format: 1h30m, 24h, 30s)
	if err := v.RegisterValidation("duration", validateDuration); err != nil {
		return err
	}

	// File path validator
	if err := v.RegisterValidation("filepath", validateFilePath); err != nil {
		return err
	}

	// Port string validator (validates string representation of port numbers)
	if err := v.RegisterValidation("port_string", validatePortString); err != nil {
		return err
	}

	return nil
}

// validateResourceQuantity validates Kubernetes resource quantities
func validateResourceQuantity(fl validator.FieldLevel) bool {
	quantity := fl.Field().String()
	if quantity == "" {
		return true // Allow empty values for omitempty
	}

	// CPU resource patterns: 100m, 0.1, 1
	cpuPattern := regexp.MustCompile(`^([0-9]+(\.[0-9]*)?|\.[0-9]+)([m]?)$`)

	// Memory resource patterns: 128Mi, 1Gi, 512Ki, 1000000000 (bytes)
	memoryPattern := regexp.MustCompile(`^([0-9]+(\.[0-9]*)?|\.[0-9]+)([KMGTPE]i?|[kmgtpe])?$`)

	return cpuPattern.MatchString(quantity) || memoryPattern.MatchString(quantity)
}

// validateDuration validates Go duration strings
func validateDuration(fl validator.FieldLevel) bool {
	duration := fl.Field().String()
	if duration == "" {
		return true // Allow empty values for omitempty
	}

	_, err := time.ParseDuration(duration)
	return err == nil
}

// validateFilePath validates file paths
func validateFilePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true // Allow empty values for omitempty
	}

	// Check if path is valid (not necessarily existing)
	_, err := filepath.Abs(path)
	return err == nil
}

// validatePortString validates port number strings (1-65535)
func validatePortString(fl validator.FieldLevel) bool {
	portStr := fl.Field().String()
	if portStr == "" {
		return true // Allow empty values for omitempty
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false // Not a valid integer
	}

	return port >= 1 && port <= 65535
}

// GetValidatorWithCustomValidations returns a validator instance with all custom validations registered
func GetValidatorWithCustomValidations() (*validator.Validate, error) {
	v := validator.New()

	if err := RegisterCustomValidators(v); err != nil {
		return nil, fmt.Errorf("failed to register custom validators: %w", err)
	}

	return v, nil
}

// ValidateStruct validates a struct using the custom validator
func ValidateStruct(s interface{}) error {
	v, err := GetValidatorWithCustomValidations()
	if err != nil {
		return err
	}

	return v.Struct(s)
}
