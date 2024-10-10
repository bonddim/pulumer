package stacks

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type Stacks []stack

func NewStacks() Stacks {
	return Stacks{}
}

// GetStacks returns a list of stacks found in the working directory
func (s *Stacks) New(paths []string) error {
	stacks := make(Stacks, 0, len(paths))

	for _, path := range paths {
		stack, err := newStackFromFile(path)
		if err != nil {
			return err
		}

		stacks = append(stacks, stack)
	}

	*s = stacks
	return nil
}

func (s *Stacks) Changed(diff []string) {
	stacks := *s
	for i := 0; i < len(stacks); i++ {
		if !stacks[i].isChanged(diff) {
			stacks = append(stacks[:i], stacks[i+1:]...)
			i--
		}
	}
	*s = stacks
}

// Sort sorts the stacks based on their dependencies
func (s *Stacks) Sort() error {
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
func (s *Stacks) Reverse() {
	stacks := *s

	for i, j := 0, len(stacks)-1; i < j; i, j = i+1, j-1 {
		stacks[i], stacks[j] = stacks[j], stacks[i]
	}
	*s = stacks
}
