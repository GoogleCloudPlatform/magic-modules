package cache

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadCache(cachePath string) (map[string][]int, error) {
	b, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var touchedFiles map[string][]int
	err = yaml.Unmarshal(b, &touchedFiles)
	return touchedFiles, err
}

func WriteCache(cachePath string, touchedFiles map[string][]int) error {
	b, err := yaml.Marshal(touchedFiles)
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, b, 0666)
}
