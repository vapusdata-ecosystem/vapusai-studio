package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteTomlFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Define test data
	data := struct {
		Name string
		Age  int
	}{
		Name: "John Doe",
		Age:  30,
	}

	// Define the expected file path
	expectedFilePath := filepath.Join(tempDir, "test.toml")

	// Call the WriteTomlFile function
	err := WriteTomlFile(data, "test", tempDir)
	if err != nil {
		t.Errorf("Error writing to TOML file: %v", err)
	}

	// Check if the file exists
	_, err = os.Stat(expectedFilePath)
	if err != nil {
		t.Errorf("Expected file %s does not exist", expectedFilePath)
	}

	// Read the contents of the file
	var readData struct {
		Name string
		Age  int
	}
	err = ReadTomlFile(&readData, "test", tempDir)
	if err != nil {
		t.Errorf("Error reading from TOML file: %v", err)
	}

	// Check if the read data matches the original data
	if !reflect.DeepEqual(data, readData) {
		t.Errorf("Read data does not match original data")
	}
}

func TestReadTomlFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Define test data
	data := struct {
		Name string
		Age  int
	}{
		Name: "John Doe",
		Age:  30,
	}

	// Write the test data to a TOML file
	err := WriteTomlFile(data, "test", tempDir)
	if err != nil {
		t.Errorf("Error writing to TOML file: %v", err)
	}

	// Read the contents of the file
	var readData struct {
		Name string
		Age  int
	}
	err = ReadTomlFile(&readData, "test", tempDir)
	if err != nil {
		t.Errorf("Error reading from TOML file: %v", err)
	}

	// Check if the read data matches the original data
	if !reflect.DeepEqual(data, readData) {
		t.Errorf("Read data does not match original data")
	}
}
