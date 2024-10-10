package stacks

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/bonddim/pulumer/internal/config"
	"github.com/bonddim/pulumer/internal/utils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type Stacks []stack

var pattern = regexp.MustCompile(`^Pulumi\..*\.yaml$`)

// GetStacks returns a list of stacks found in the working directory
func GetStacks(cfg config.Config) (Stacks, error) {
	paths := utils.SearchFiles(cfg.WorkDir, cfg.MaxDepth, pattern)
	stacks := make(Stacks, 0, len(paths))

	for _, path := range paths {
		s, err := newStackFromFile(path)
		if err != nil {
			return nil, err
		}
		s.Config = cfg
		s.Workspace, err = filepath.Rel(cfg.WorkDir, s.Workspace)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path for %s\n%v", s.Workspace, err)
		}
		stacks = append(stacks, s)
	}

	return stacks, nil
}

// Sorted sorts the stacks based on their dependencies
func (s *Stacks) Sorted() error {
	// Create a new directed graph
	g := simple.NewDirectedGraph()

	// Create a map to hold nodes corresponding to stacks
	nodeMap := make(map[string]graph.Node)

	// Create nodes and add them to the graph
	for _, stack := range *s {
		n := g.NewNode()
		g.AddNode(n)
		nodeMap[stack.Id] = n
		nodeMap[stack.Project] = n
	}

	for _, stack := range *s {
		// Add directed edges between the nodes based on DependsOn criteria
		for _, dep := range stack.DependsOn {
			g.SetEdge(g.NewEdge(nodeMap[dep], nodeMap[stack.Id]))
		}
	}

	// Perform a topological sort using SortStabilized
	order, err := topo.SortStabilized(g, nil)
	if err != nil {
		return fmt.Errorf("performing topological sort: %v", err)
	}

	// Create a slice to hold the sorted stacks
	sortedStacks := make([]stack, len(order))
	for i, n := range order {
		for _, stack := range *s {
			if nodeMap[stack.Id].ID() == n.ID() {
				sortedStacks[i] = stack
				break
			}
		}
	}

	//Update the original slice with the sorted stacks
	*s = sortedStacks

	return nil
}

// Reverse reverses the order of stacks
func (s Stacks) Reverse() Stacks {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s Stacks) FilterChanged(diff []string) Stacks {
	for i := 0; i < len(s); i++ {
		if !s[i].isChanged(diff) {
			s = append(s[:i], s[i+1:]...)
			i--
		}
	}
	return s
}
