package lib

import (
	"encoding/json"
	"testing"
)

func TestStringOrSlice_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StringOrSlice
	}{
		{
			name:     "single string",
			input:    `{"digest":"sha256:abc123"}`,
			expected: StringOrSlice{"sha256:abc123"},
		},
		{
			name:     "empty string",
			input:    `{"digest":""}`,
			expected: StringOrSlice{""},
		},
		{
			name:     "array of strings",
			input:    `{"digest":["sha256:abc123","sha256:def456"]}`,
			expected: StringOrSlice{"sha256:abc123", "sha256:def456"},
		},
		{
			name:     "empty array",
			input:    `{"digest":[]}`,
			expected: StringOrSlice{},
		},
		{
			name:     "null value",
			input:    `{"digest":null}`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data
			if err := json.Unmarshal([]byte(tt.input), &data); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(data.Digest) != len(tt.expected) {
				t.Fatalf("expected %d elements, got %d", len(tt.expected), len(data.Digest))
			}

			for i, v := range tt.expected {
				if data.Digest[i] != v {
					t.Errorf("element %d: expected %q, got %q", i, v, data.Digest[i])
				}
			}
		})
	}
}

func TestResourceStateConstants(t *testing.T) {
	if ResourceStateActive != "active" {
		t.Errorf("expected 'active', got %q", ResourceStateActive)
	}
	if ResourceStateInactive != "inactive" {
		t.Errorf("expected 'inactive', got %q", ResourceStateInactive)
	}
}

func TestResourceStateJSON(t *testing.T) {
	req := UpdateComponentTypeRequest{Name: "test", State: ResourceStateActive}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded UpdateComponentTypeRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded.State != ResourceStateActive {
		t.Errorf("expected %q, got %q", ResourceStateActive, decoded.State)
	}
}

func TestStringOrSlice_UnmarshalJSON_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "number value",
			input: `{"digest":42}`,
		},
		{
			name:  "boolean value",
			input: `{"digest":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data
			if err := json.Unmarshal([]byte(tt.input), &data); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
