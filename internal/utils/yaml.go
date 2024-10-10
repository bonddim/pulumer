package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

// YamlUnmarshal reads the content of a YAML file and unmarshals it into the given struct.
func YamlUnmarshal(filepath string, out interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}
