// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mnptu

import (
	"testing"
)

func TestNullGraphWalker_impl(t *testing.T) {
	var _ GraphWalker = NullGraphWalker{}
}
