// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package statefile

import (
	"io"

	tfversion "github.com/hashicorp/mnptu/version"
)

// Write writes the given state to the given writer in the current state
// serialization format.
func Write(s *File, w io.Writer) error {
	// Always record the current mnptu version in the state.
	s.mnptuVersion = tfversion.SemVer

	diags := writeStateV4(s, w)
	return diags.Err()
}

// WriteForTest writes the given state to the given writer in the current state
// serialization format without recording the current mnptu version. This is
// intended for use in tests that need to override the current mnptu
// version.
func WriteForTest(s *File, w io.Writer) error {
	diags := writeStateV4(s, w)
	return diags.Err()
}
