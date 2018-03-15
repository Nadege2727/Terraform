package terraform

// NodeApplyableProvider represents a provider during an apply.
type NodeApplyableProvider struct {
	*NodeAbstractProvider
}

// GraphNodeEvaluable
func (n *NodeApplyableProvider) EvalTree() EvalNode {
	return ProviderEvalTree(n, n.ProviderConfig())
}
