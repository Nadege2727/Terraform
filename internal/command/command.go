package command

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hashicorp/terraform/internal/terraform"
	"github.com/hashicorp/terraform/internal/tfdiags"
)

// Set to true when we're testing
var test bool = false

// DefaultDataDir is the default directory for storing local data.
const DefaultDataDir = ".terraform"

// PluginPathFile is the name of the file in the data dir which stores the list
// of directories supplied by the user with the `-plugin-dir` flag during init.
const PluginPathFile = "plugin_path"

// pluginMachineName is the directory name used in new plugin paths.
const pluginMachineName = runtime.GOOS + "_" + runtime.GOARCH

// DefaultPluginVendorDir is the location in the config directory to look for
// user-added plugin binaries. Terraform only reads from this path if it
// exists, it is never created by terraform.
const DefaultPluginVendorDir = "terraform.d/plugins/" + pluginMachineName

// DefaultStateFilename is the default filename used for the state file.
const DefaultStateFilename = "terraform.tfstate"

// DefaultVarsFilename is the default filename used for vars
const DefaultVarsFilename = "terraform.tfvars"

// DefaultBackupExtension is added to the state file to form the path
const DefaultBackupExtension = ".backup"

// DefaultParallelism is the limit Terraform places on total parallel
// operations as it walks the dependency graph.
const DefaultParallelism = 10

// ErrUnsupportedLocalOp is the common error message shown for operations
// that require a backend.Local.
const ErrUnsupportedLocalOp = `The configured backend doesn't support this operation.

The "backend" in Terraform defines how Terraform operates. The default
backend performs all operations locally on your machine. Your configuration
is configured to use a non-local backend. This backend doesn't support this
operation.
`

// ModulePath returns the path to the root module and validates CLI arguments.
//
// This centralizes the logic for any commands that previously accepted
// a module path via CLI arguments. This will error if any extraneous arguments
// are given and suggest using the -chdir flag instead.
//
// If your command accepts more than one arg, then change the slice bounds
// to pass validation.
func ModulePath(args []string) (string, error) {
	// TODO: test

	if len(args) > 0 {
		return "", fmt.Errorf("Too many command line arguments. Did you mean to use -chdir?")
	}

	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting pwd: %s", err)
	}

	return path, nil
}

// CheckRequiredVersion loads the config and check if the
// core version requirements are satisfied.
func CheckRequiredVersion(c *Meta) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	loader, err := c.initConfigLoader()
	if err != nil {
		diags = diags.Append(err)
		return diags
	}

	pwd, err := os.Getwd()
	fmt.Println(pwd)
	if err != nil {
		diags = diags.Append(fmt.Errorf("Error getting pwd: %s", err))
		return diags
	}

	config, configDiags := loader.LoadConfig(pwd)
	if configDiags.HasErrors() {
		diags = diags.Append(configDiags)
		return diags
	}

	versionDiags := terraform.CheckCoreVersionRequirements(config)
	if versionDiags.HasErrors() {
		diags = diags.Append(versionDiags)
		return diags
	}

	return nil
}
