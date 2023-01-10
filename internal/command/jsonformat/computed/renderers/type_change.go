package renderers

import (
	"fmt"

	"github.com/hashicorp/terraform/internal/command/jsonformat/computed"
)

func TypeChange(before, after computed.Diff) computed.DiffRenderer {
	return &typeChangeRenderer{
		before: before,
		after:  after,
	}
}

type typeChangeRenderer struct {
	NoWarningsRenderer

	before computed.Diff
	after  computed.Diff
}

func (renderer typeChangeRenderer) RenderHuman(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {
	opts.OverrideNullSuffix = true // Never render null suffix for children of type changes.
	return fmt.Sprintf("%s [yellow]->[reset] %s", renderer.before.RenderHuman(indent, opts), renderer.after.RenderHuman(indent, opts))
}
