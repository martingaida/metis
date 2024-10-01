package main

import (
	"strings"
	"testing"
)

func TestCallLLMMicroservice(t *testing.T) {
	// Assuming callLLMMicroservice takes a string parameter
	response, err := callLLMMicroservice("sample input", "Basic", false)

	if err != nil {
		if strings.Contains(err.Error(), "unsupported protocol scheme") {
			t.Skip("Skipping test due to network configuration issue. Error: ", err)
		} else {
			t.Fatalf("callLLMMicroservice returned an unexpected error: %v", err)
		}
	}

	if response.Topics == nil {
		t.Error("Expected Topics to be initialized, got nil")
	} else if len(response.Topics) == 0 {
		t.Error("Expected at least one topic, got none")
	}

	if response.MainTakeaway == "" {
		t.Error("Expected a non-empty main takeaway, got an empty string")
	}

	if response.MainTakeaway != "Sample takeaway" {
		t.Errorf("Expected main takeaway 'Sample takeaway', got '%s'", response.MainTakeaway)
	}
}
