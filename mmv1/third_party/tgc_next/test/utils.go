package test

import (
	"encoding/json"
	"fmt"
	"os"
)

// Writes the data into a JSON file
func writeJSONFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling data for %s: %v\n", filename, err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file %s: %v\n", filename, err)
	}
	return nil
}
