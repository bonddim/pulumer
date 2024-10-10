package stacks

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bonddim/pulumer/internal/config"
	"github.com/bonddim/pulumer/internal/utils"
)

type stack struct {
	Id        string
	Name      string
	Project   string
	Workspace string
	DependsOn []string `yaml:"dependsOn"`
	Watch     []string `yaml:"watch"`
	Changed   bool
	Config    config.Config
}

// newStackFromFile creates a new stack from a Pulumi stack file
func newStackFromFile(path string) (stack, error) {
	var s stack

	// Get stack configuration
	if err := utils.YamlUnmarshal(path, &s); err != nil {
		return stack{}, fmt.Errorf("failed unmarshalling stack file %s\n%v", path, err)
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
	if err := utils.YamlUnmarshal(path, &p); err != nil {
		return fmt.Errorf("failed unmarshalling project file %s\n%v", path, err)
	}

	// Set stack project name
	s.Project = p.Name
	// Set stack ID
	s.Id = fmt.Sprintf("%s/%s", p.Name, s.Name)
	// Append project dependencies to stack dependencies
	s.DependsOn = append(s.DependsOn, p.DependsOn...)
	s.Watch = append(s.Watch, p.Watch...)

	return nil
}

// isChanged checks if the stack has changes
func (s *stack) isChanged(diff []string) bool {
	stackFiles := append(s.Watch, s.Workspace)

	for _, file := range stackFiles {
		for _, changedFile := range diff {
			if strings.Contains(changedFile, file) {
				return true
			}
		}
	}
	return false
}
