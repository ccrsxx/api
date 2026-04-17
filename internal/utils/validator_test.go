package utils

import (
	"errors"
	"slices"
	"testing"
)

// Each struct is minimal and named for what it tests.
// This avoids the noise of setting up unrelated fields to pass validation.

type requiredOnlyStruct struct {
	Title string `validate:"required"`
}

type minStruct struct {
	Age int `validate:"required,min=18"`
}

type maxStruct struct {
	Score int `validate:"required,max=100"`
}

type oneofStruct struct {
	Type string `validate:"required,oneof=blog project"`
}

type defaultParamStruct struct {
	Code string `validate:"required,len=5"`
}

type contentTypeStruct struct {
	ContentType string `validate:"content_type"`
}

// multiErrorStruct is used only for the multiple-errors case.
// Fields are ordered to match the expected error output.
type multiErrorStruct struct {
	Age         int    `validate:"required,min=18"`
	Type        string `validate:"required,oneof=blog project"`
	ContentType string `validate:"content_type"`
}

func TestFormatValidationError(t *testing.T) {
	tests := []struct {
		name            string
		errSetup        func() error
		expectedSummary string
		expectedDetails []string
	}{
		{
			name: "required field missing",
			errSetup: func() error {
				return Validate.Struct(requiredOnlyStruct{Title: ""})
			},
			expectedSummary: "invalid input: field 'Title' failed on 'required' validation",
			expectedDetails: []string{
				"field 'Title' failed on 'required' validation",
			},
		},
		{
			name: "min constraint violated",
			errSetup: func() error {
				return Validate.Struct(minStruct{Age: 15})
			},
			expectedSummary: "invalid input: field 'Age' failed on 'min' validation (minimum: 18)",
			expectedDetails: []string{
				"field 'Age' failed on 'min' validation (minimum: 18)",
			},
		},
		{
			name: "max constraint violated",
			errSetup: func() error {
				return Validate.Struct(maxStruct{Score: 105})
			},
			expectedSummary: "invalid input: field 'Score' failed on 'max' validation (maximum: 100)",
			expectedDetails: []string{
				"field 'Score' failed on 'max' validation (maximum: 100)",
			},
		},
		{
			name: "oneof constraint violated",
			errSetup: func() error {
				return Validate.Struct(oneofStruct{Type: "video"})
			},
			expectedSummary: "invalid input: field 'Type' failed on 'oneof' validation (must be one of: blog, project)",
			expectedDetails: []string{
				"field 'Type' failed on 'oneof' validation (must be one of: blog, project)",
			},
		},
		{
			name: "unknown param falls back to rule label",
			errSetup: func() error {
				return Validate.Struct(defaultParamStruct{Code: "12"})
			},
			expectedSummary: "invalid input: field 'Code' failed on 'len' validation (rule: 5)",
			expectedDetails: []string{
				"field 'Code' failed on 'len' validation (rule: 5)",
			},
		},
		{
			name: "custom content_type constraint violated",
			errSetup: func() error {
				return Validate.Struct(contentTypeStruct{ContentType: "video"})
			},
			expectedSummary: "invalid input: field 'ContentType' failed on 'content_type' validation (must be one of: blog, project)",
			expectedDetails: []string{
				"field 'ContentType' failed on 'content_type' validation (must be one of: blog, project)",
			},
		},
		{
			name: "multiple errors are joined in field-definition order",
			errSetup: func() error {
				return Validate.Struct(multiErrorStruct{
					Age:         15,      // fails min=18
					Type:        "video", // fails oneof=blog project
					ContentType: "video", // fails content_type
				})
			},
			expectedSummary: "invalid input: " +
				"field 'Age' failed on 'min' validation (minimum: 18), " +
				"field 'Type' failed on 'oneof' validation (must be one of: blog, project), " +
				"field 'ContentType' failed on 'content_type' validation (must be one of: blog, project)",
			expectedDetails: []string{
				"field 'Age' failed on 'min' validation (minimum: 18)",
				"field 'Type' failed on 'oneof' validation (must be one of: blog, project)",
				"field 'ContentType' failed on 'content_type' validation (must be one of: blog, project)",
			},
		},
		{
			name: "non-validation error returns generic message",
			errSetup: func() error {
				return errors.New("some unexpected error")
			},
			expectedSummary: "invalid input format",
			expectedDetails: nil,
		},
		{
			name: "nil error returns generic message",
			errSetup: func() error {
				return nil
			},
			expectedSummary: "invalid input format",
			expectedDetails: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errSetup()
			summary, details := FormatValidationError(err)

			if summary != tt.expectedSummary {
				t.Errorf("summary\ngot:  %s\nwant: %s", summary, tt.expectedSummary)
			}

			if !slices.Equal(details, tt.expectedDetails) {
				t.Errorf("details\ngot:  %v\nwant: %v", details, tt.expectedDetails)
			}
		})
	}
}
