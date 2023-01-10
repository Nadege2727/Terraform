package differ

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/internal/command/jsonformat/collections"
	"github.com/hashicorp/terraform/internal/command/jsonformat/computed"
	"github.com/hashicorp/terraform/internal/command/jsonformat/computed/renderers"
	"github.com/hashicorp/terraform/internal/command/jsonprovider"
	"github.com/hashicorp/terraform/internal/plans"
)

func (change Change) computeAttributeDiffAsMap(elementType cty.Type) computed.Diff {
	mapValue := change.asMap()
	elements, current := collections.TransformMap(mapValue.Before, mapValue.After, func(key string) computed.Diff {
		return mapValue.getChild(key).computeDiffForType(elementType)
	})
	return computed.NewDiff(renderers.Map(elements), current, change.replacePath())
}

func (change Change) computeAttributeDiffAsNestedMap(attributes map[string]*jsonprovider.Attribute) computed.Diff {
	mapValue := change.asMap()
	elements, current := collections.TransformMap(mapValue.Before, mapValue.After, func(key string) computed.Diff {
		return mapValue.getChild(key).computeDiffForNestedAttribute(&jsonprovider.NestedType{
			Attributes:  attributes,
			NestingMode: "single",
		})
	})
	return computed.NewDiff(renderers.Map(elements), current, change.replacePath())
}

func (change Change) computeBlockDiffsAsMap(block *jsonprovider.Block) ([]computed.Diff, plans.Action) {
	mapValue := change.asMap()
	elements, action := collections.TransformMap(mapValue.Before, mapValue.After, func(key string) computed.Diff {
		return mapValue.getChild(key).ComputeDiffForBlock(block)
	})

	var ret []computed.Diff
	for _, element := range elements {
		ret = append(ret, element)
	}
	return ret, action
}

func (change Change) processMap(process func(key string, value Change)) {
	mapValue := change.asMap()

	handled := make(map[string]bool)
	for key := range mapValue.Before {
		handled[key] = true
		process(key, mapValue.getChild(key))
	}
	for key := range mapValue.After {
		if _, ok := handled[key]; ok {
			continue
		}
		process(key, mapValue.getChild(key))
	}
}
