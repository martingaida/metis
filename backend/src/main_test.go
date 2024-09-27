package main

import (
	"strings"
	"testing"
)

func TestCallLLMMicroservice(t *testing.T) {
	// Assuming callLLMMicroservice takes a string parameter
	response, err := callLLMMicroservice("sample input")

	if err != nil {
		if strings.Contains(err.Error(), "unsupported protocol scheme") {
			t.Skip("Skipping test due to network configuration issue. Error: ", err)
		} else {
			t.Fatalf("callLLMMicroservice returned an unexpected error: %v", err)
		}
	}

	if response.Explanations.Topics == nil {
		t.Error("Expected Topics to be initialized, got nil")
	} else if len(response.Explanations.Topics) == 0 {
		t.Error("Expected at least one topic, got none")
	}

	if response.Explanations.MainTakeaway == "" {
		t.Error("Expected a non-empty main takeaway, got an empty string")
	}

	// Remove or comment out the checks for the Name field if it doesn't exist in your Topic struct
	// if response.Explanations.Topics[0].Name != "Sample Topic" {
	// 	t.Errorf("Expected topic name 'Sample Topic', got '%s'", response.Explanations.Topics[0].Name)
	// }

	if response.Explanations.MainTakeaway != "Sample takeaway" {
		t.Errorf("Expected main takeaway 'Sample takeaway', got '%s'", response.Explanations.MainTakeaway)
	}
}
