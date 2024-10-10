package stacks

import (
	"fmt"
	"path/filepath"
	"strings"
)

type stack struct {
	Id        string
	Name      string
	Project   string
	Workspace string
	DependsOn []string `yaml:"dependsOn"`
	Watch     []string `yaml:"watch"`
}

// newStackFromFile creates a new stack from a Pulumi stack file
func newStackFromFile(path string) (stack, error) {
	var s stack

	// Get stack configuration
	if err := yamlUnmarshal(path, &s); err != nil {
		return s, fmt.Errorf("failed unmarshalling stack file %s\n%v", path, err)
	}

	s.Name = strings.Split(filepath.Base(path), ".")[1]
	s.Workspace = filepath.Dir(path)
	s.Watch = append(s.Watch, s.Workspace)

	// Get project configuration
	if err := s.mergeProjectConfig(); err != nil {
		return stack{}, err
	}

	return s, nil
}

// mergeProjectConfig applies the project configuration to the stack
func (s *stack) mergeProjectConfig() error {
	var p struct {
		Name      string   `yaml:"name"`
		DependsOn []string `yaml:"dependsOn"`
		Watch     []string `yaml:"watch"`
	}

	path := filepath.Join(s.Workspace, "Pulumi.yaml")
	if err := yamlUnmarshal(path, &p); err != nil {
		return fmt.Errorf("failed unmarshalling project file %s\n%v", path, err)
	}

	// Set stack project name
	s.Project = p.Name
	// Set stack ID
	s.Id = fmt.Sprintf("%s/%s", p.Name, s.Name)
	// Append project dependencies to stack dependencies
	s.DependsOn = append(s.DependsOn, p.DependsOn...)
	// Append project watch paths to stack watch paths
	s.Watch = append(s.Watch, p.Watch...)

	return nil
}

// isChanged checks if the stack has changes
func (s *stack) isChanged(diff []string) bool {
	for _, file := range s.Watch {
		for _, changedFile := range diff {
			if strings.Contains(changedFile, file) {
				return true
			}
		}
	}
	return false
}
