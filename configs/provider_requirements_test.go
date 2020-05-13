package configs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcltest"
	"github.com/hashicorp/terraform/addrs"
	"github.com/zclconf/go-cty/cty"
)

var (
	ignoreUnexported = cmpopts.IgnoreUnexported(version.Constraint{})
	comparer         = cmp.Comparer(func(x, y RequiredProvider) bool {
		if x.Name != y.Name {
			return false
		}
		if x.Type != y.Type {
			return false
		}
		if x.Source != y.Source {
			return false
		}
		if x.Requirement.Required.String() != y.Requirement.Required.String() {
			return false
		}
		if x.DeclRange != y.DeclRange {
			return false
		}
		return true
	})
	blockRange = hcl.Range{
		Filename: "mock.tf",
		Start:    hcl.Pos{Line: 3, Column: 12, Byte: 27},
		End:      hcl.Pos{Line: 3, Column: 19, Byte: 34},
	}
	mockRange = hcl.Range{
		Filename: "MockExprLiteral",
	}
)

func TestDecodeRequiredProvidersBlock(t *testing.T) {
	tests := map[string]struct {
		Block *hcl.Block
		Want  *RequiredProviders
		Error string
	}{
		"legacy": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"default": {
							Name: "default",
							Expr: hcltest.MockExprLiteral(cty.StringVal("1.0.0")),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"default": {
						Name:        "default",
						Type:        addrs.NewDefaultProvider("default"),
						Requirement: testVC("1.0.0"),
						DeclRange:   mockRange,
					},
				},
				DeclRange: blockRange,
			},
		},
		"provider source": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"my-test": {
							Name: "my-test",
							Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
								"source":  cty.StringVal("mycloud/test"),
								"version": cty.StringVal("2.0.0"),
							})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"my-test": {
						Name:        "my-test",
						Source:      "mycloud/test",
						Type:        addrs.NewProvider(addrs.DefaultRegistryHost, "mycloud", "test"),
						Requirement: testVC("2.0.0"),
						DeclRange:   mockRange,
					},
				},
				DeclRange: blockRange,
			},
		},
		"mixed": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"legacy": {
							Name: "legacy",
							Expr: hcltest.MockExprLiteral(cty.StringVal("1.0.0")),
						},
						"my-test": {
							Name: "my-test",
							Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
								"source":  cty.StringVal("mycloud/test"),
								"version": cty.StringVal("2.0.0"),
							})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"legacy": {
						Name:        "legacy",
						Type:        addrs.NewDefaultProvider("legacy"),
						Requirement: testVC("1.0.0"),
						DeclRange:   mockRange,
					},
					"my-test": {
						Name:        "my-test",
						Source:      "mycloud/test",
						Type:        addrs.NewProvider(addrs.DefaultRegistryHost, "mycloud", "test"),
						Requirement: testVC("2.0.0"),
						DeclRange:   mockRange,
					},
				},
				DeclRange: blockRange,
			},
		},
		"version-only block": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"test": {
							Name: "test",
							Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
								"version": cty.StringVal("~>2.0.0"),
							})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"test": {
						Name:        "test",
						Type:        addrs.NewDefaultProvider("test"),
						Requirement: testVC("~>2.0.0"),
						DeclRange:   mockRange,
					},
				},
				DeclRange: blockRange,
			},
		},
		"invalid source": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"my-test": {
							Name: "my-test",
							Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
								"source":  cty.StringVal("some/invalid/provider/source/test"),
								"version": cty.StringVal("~>2.0.0"),
							})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"my-test": {
						Name:   "my-test",
						Source: "some/invalid/provider/source/test",
						// decodeRequiredProviders falls back to determining the type from the local name
						Type:        addrs.NewDefaultProvider("my-test"),
						Requirement: testVC("~>2.0.0"),
						DeclRange:   mockRange,
					},
				},
				DeclRange: blockRange,
			},
			Error: "Invalid provider source string",
		},
		"invalid localname": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"my_test": {
							Name: "my_test",
							Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
								"version": cty.StringVal("~>2.0.0"),
							})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"my_test": {
						Name:        "my_test",
						Type:        addrs.Provider{},
						Requirement: testVC("~>2.0.0"),
						DeclRange:   mockRange,
					},
				},
				DeclRange: blockRange,
			},
			Error: "Invalid provider local name",
		},
		"version constraint error": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"my-test": {
							Name: "my-test",
							Expr: hcltest.MockExprLiteral(cty.ObjectVal(map[string]cty.Value{
								"source":  cty.StringVal("mycloud/test"),
								"version": cty.StringVal("invalid"),
							})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"my-test": {
						Name:      "my-test",
						Source:    "mycloud/test",
						Type:      addrs.NewProvider(addrs.DefaultRegistryHost, "mycloud", "test"),
						DeclRange: mockRange,
					},
				},
				DeclRange: blockRange,
			},
			Error: "Invalid version constraint",
		},
		"invalid required_providers attribute value": {
			Block: &hcl.Block{
				Type: "required_providers",
				Body: hcltest.MockBody(&hcl.BodyContent{
					Attributes: hcl.Attributes{
						"test": {
							Name: "test",
							Expr: hcltest.MockExprLiteral(cty.ListVal([]cty.Value{cty.StringVal("2.0.0")})),
						},
					},
				}),
				DefRange: blockRange,
			},
			Want: &RequiredProviders{
				RequiredProviders: map[string]*RequiredProvider{
					"test": {
						Name:      "test",
						Type:      addrs.NewDefaultProvider("test"),
						DeclRange: mockRange,
					},
				},
				DeclRange: blockRange,
			},
			Error: "Invalid required_providers syntax",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, diags := decodeRequiredProvidersBlock(test.Block)
			if diags.HasErrors() {
				if test.Error == "" {
					t.Fatalf("unexpected error")
				}
				if gotErr := diags[0].Summary; gotErr != test.Error {
					t.Errorf("wrong error, got %q, want %q", gotErr, test.Error)
				}
			} else if test.Error != "" {
				t.Fatalf("expected error")
			}

			if !cmp.Equal(got, test.Want, ignoreUnexported, comparer) {
				t.Fatalf("wrong result:\n %s", cmp.Diff(got, test.Want, ignoreUnexported, comparer))
			}
		})
	}
}

func testVC(ver string) VersionConstraint {
	constraint, _ := version.NewConstraint(ver)
	return VersionConstraint{
		Required:  constraint,
		DeclRange: hcl.Range{},
	}
}
