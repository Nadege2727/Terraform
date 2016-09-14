package terraform

import (
	"github.com/hashicorp/terraform/config/module"
)

// ApplyGraphBuilder implements GraphBuilder and is responsible for building
// a graph for applying a Terraform diff.
//
// Because the graph is built from the diff (vs. the config or state),
// this helps ensure that the apply-time graph doesn't modify any resources
// that aren't explicitly in the diff. There are other scenarios where the
// diff can be deviated, so this is just one layer of protection.
type ApplyGraphBuilder struct {
	// Config is the root module for the graph to build.
	Config *module.Tree

	// Diff is the diff to apply.
	Diff *Diff

	// State is the current state
	State *State

	// Providers is the list of providers supported.
	Providers []string

	// Provisioners is the list of provisioners supported.
	Provisioners []string
}

// See GraphBuilder
func (b *ApplyGraphBuilder) Build(path []string) (*Graph, error) {
	return (&BasicGraphBuilder{
		Steps:    b.Steps(),
		Validate: true,
	}).Build(path)
}

// See GraphBuilder
func (b *ApplyGraphBuilder) Steps() []GraphTransformer {
	steps := []GraphTransformer{
		// Creates all the nodes represented in the diff.
		&DiffTransformer{
			Diff:   b.Diff,
			Config: b.Config,
			State:  b.State,
		},

		// Create all the providers
		&MissingProviderTransformer{Providers: b.Providers},
		&ProviderTransformer{},

		// Single root
		&RootTransformer{},
	}

	return steps
}
