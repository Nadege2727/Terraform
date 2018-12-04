package jsonplan

import (
	"encoding/json"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/terraform"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// stateValues is the common representation of resolved values for both the
// prior state (which is always complete) and the planned new state.
type stateValues struct {
	Outputs    map[string]output `json:"outputs,omitempty"`
	RootModule module            `json:"root_module,omitempty"`
}

// marshalPlannedValues takes the current state, planned changes, and schemas
// and populates the PlannedValues. Any unknown values will be omitted.
func (p *plan) marshalPlannedValues(
	changes *plans.Changes,
	s *states.State,
	schemas *terraform.Schemas,
) error {

	// var curr, planned stateValues
	var curr stateValues
	// marshal the current state into a stateValues
	err := curr.MarshalState(s, schemas)
	if err != nil {
		return err
	}

	// // marshal the changes into a stateValues
	// planned.MarshalChanges(changes, schemas)

	// // merge the planned stateValues into the current stateValues
	// plannedChanges := curr.Merge(planned)

	// p.PlannedValues.RootModule = plannedChanges

	outputs, err := marshalPlannedOutputs(changes, s)
	if err != nil {
		return err
	}
	p.PlannedValues.Outputs = outputs

	return nil
}

// attributeValues is the JSON representation of the attribute values of the
// resource, whose structure depends on the resource type schema.
type attributeValues map[string]interface{}

func marshalAttributeValues() attributeValues {
	return attributeValues{}
}

func marshalPlannedOutputs(changes *plans.Changes, s *states.State) (map[string]output, error) {
	ret := make(map[string]output)

	// add the current state's outputs to the map
	if !s.Empty() {
		for k, v := range s.RootModule().OutputValues {
			if v.Value != cty.NilVal {
				outputVal, _ := ctyjson.Marshal(v.Value, v.Value.Type())
				ret[k] = output{
					Value:     outputVal,
					Sensitive: v.Sensitive,
				}
			}
		}
	}

	if changes.Outputs == nil {
		// No changes - we're done here!
		return ret, nil
	}

	// overwrite the current state's outputs with any changes
	// this will also add any outputs not in the state
	for _, oc := range changes.Outputs {
		if oc.ChangeSrc.Action == plans.Delete {
			delete(ret, oc.Addr.String())
		}

		var after []byte
		changeV, err := oc.Decode()
		if err != nil {
			return ret, err
		}

		if changeV.After != cty.NilVal {
			if changeV.After.IsWhollyKnown() {
				after, err = ctyjson.Marshal(changeV.After, changeV.After.Type())
				if err != nil {
					return ret, err
				}
			}
		}

		ret[oc.Addr.OutputValue.Name] = output{
			Value:     json.RawMessage(after),
			Sensitive: oc.Sensitive,
		}
	}

	return ret, nil

}

func (sv *stateValues) MarshalState(s *states.State, schemas *terraform.Schemas) error {
	if s.Empty() {
		return nil
	}

	// var rootModule module
	// var rs []resource

	// start with the root module
	// rs, err := marshalResources(s.RootModule().Resources, schemas)

	return nil
}

// func marshalResources(resources map[string]*states.Resource, schemas *terraform.Schemas) ([]resource, error) {
// 	var rs []resource
// 	for _, v := range resources {

// 	}
// 	return rs, nil
// }
