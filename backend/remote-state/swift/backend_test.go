package swift

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/state/remote"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/states/statefile"
)

// verify that we are doing ACC tests or the Swift tests specifically
func testACC(t *testing.T) {
	skip := os.Getenv("TF_ACC") == "" && os.Getenv("TF_SWIFT_TEST") == ""
	if skip {
		t.Log("swift backend tests require setting TF_ACC or TF_SWIFT_TEST")
		t.Skip()
	}
	t.Log("swift backend acceptance tests enabled")
}

func TestBackend_impl(t *testing.T) {
	var _ backend.Backend = new(Backend)
}

func testAccPreCheck(t *testing.T) {
	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}
}

func TestBackendConfig(t *testing.T) {
	testACC(t)

	// Build config
	config := map[string]interface{}{
		"archive_container": "test-tfstate-archive",
		"container":         "test-tfstate",
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)

	if b.container != "test-tfstate" {
		t.Fatal("Incorrect path was provided.")
	}
	if b.archiveContainer != "test-tfstate-archive" {
		t.Fatal("Incorrect archivepath was provided.")
	}
}

func TestBackend(t *testing.T) {
	testACC(t)

	// test with a slash at the end
	container := fmt.Sprintf("terraform-state-swift-test-%x/test/", time.Now().Unix())

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"container": container,
	})).(*Backend)

	createSwiftContainer(t, b.client, container)
	defer deleteSwiftContainer(t, b.client, container)

	backend.TestBackendStates(t, b)
}

func TestBackendPath(t *testing.T) {
	testACC(t)

	// test without a slash at the end
	path := fmt.Sprintf("terraform-state-swift-test-%x/test", time.Now().Unix())
	t.Logf("[DEBUG] Generating backend config")
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"container": path,
	})).(*Backend)
	t.Logf("[DEBUG] Backend configured")

	defer deleteSwiftContainer(t, b.client, path)

	t.Logf("[DEBUG] Testing Backend")

	// Generate some state
	state1 := states.NewState()

	// RemoteClient to test with
	client := &RemoteClient{
		name:             DEFAULT_NAME,
		client:           b.client,
		archive:          b.archive,
		archiveContainer: b.archiveContainer,
		container:        b.container,
	}

	stateMgr := &remote.State{Client: client}
	stateMgr.WriteState(state1)
	if err := stateMgr.PersistState(); err != nil {
		t.Fatal(err)
	}

	if err := stateMgr.RefreshState(); err != nil {
		t.Fatal(err)
	}

	// Add some state
	mod := state1.EnsureModule(addrs.RootModuleInstance)
	mod.SetOutputValue("bar", cty.StringVal("baz"), false)
	stateMgr.WriteState(state1)
	if err := stateMgr.PersistState(); err != nil {
		t.Fatal(err)
	}

}

func TestBackendArchive(t *testing.T) {
	testACC(t)

	// without a slash at all
	container := fmt.Sprintf("terraform-state-swift-test-%x", time.Now().Unix())
	archiveContainer := fmt.Sprintf("%s_archive", container)

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"archive_container": archiveContainer,
		"container":         container,
	})).(*Backend)

	defer deleteSwiftContainer(t, b.client, container)
	defer deleteSwiftContainer(t, b.client, archiveContainer)

	// Generate some state
	state1 := states.NewState()

	// RemoteClient to test with
	client := &RemoteClient{
		name:             DEFAULT_NAME,
		client:           b.client,
		archive:          b.archive,
		archiveContainer: b.archiveContainer,
		container:        b.container,
	}

	stateMgr := &remote.State{Client: client}
	stateMgr.WriteState(state1)
	if err := stateMgr.PersistState(); err != nil {
		t.Fatal(err)
	}

	if err := stateMgr.RefreshState(); err != nil {
		t.Fatal(err)
	}

	// Add some state
	mod := state1.EnsureModule(addrs.RootModuleInstance)
	mod.SetOutputValue("bar", cty.StringVal("baz"), false)
	stateMgr.WriteState(state1)
	if err := stateMgr.PersistState(); err != nil {
		t.Fatal(err)
	}

	cnt, prefix := getContainerAndPrefix(archiveContainer)
	archiveObjects := getSwiftObjectNames(t, b.client, cnt, prefix)
	t.Logf("archiveObjects len = %d. Contents = %+v", len(archiveObjects), archiveObjects)
	if len(archiveObjects) != 1 {
		t.Fatalf("Invalid number of archive objects. Expected 1, got %d", len(archiveObjects))
	}

	// Download archive state to validate
	archiveData := downloadSwiftObject(t, b.client, archiveContainer, archiveObjects[0])
	t.Logf("Archive data downloaded... Looks like: %+v", archiveData)
	archiveStateFile, err := statefile.Read(archiveData)
	if err != nil {
		t.Fatalf("Error Reading State: %s", err)
	}

	t.Logf("Archive state lineage = %s, serial = %d", archiveStateFile.Lineage, archiveStateFile.Serial)
	if stateMgr.StateSnapshotMeta().Lineage != archiveStateFile.Lineage {
		t.Fatal("Got a different lineage")
	}

}

// Helper function to download an object in a Swift container
func downloadSwiftObject(t *testing.T, osClient *gophercloud.ServiceClient, container, object string) (data io.Reader) {
	container, prefix := getContainerAndPrefix(container)

	t.Logf("Attempting to download object %s from container %s", prefix+object, container)
	res := objects.Download(osClient, container, prefix+object, nil)
	if res.Err != nil {
		t.Fatalf("Error downloading object: %s", res.Err)
	}
	data = res.Body
	return
}

// Helper function to get a list of objects in a Swift container
func getSwiftObjectNames(t *testing.T, osClient *gophercloud.ServiceClient, container, prefix string) (objectNames []string) {
	listOpts := &objects.ListOpts{
		Prefix:    prefix,
		Delimiter: "/",
		Full:      false,
	}

	_ = objects.List(osClient, container, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		// Get a slice of object names
		names, err := objects.ExtractNames(page)
		if err != nil {
			t.Fatalf("Error extracting object names from page: %s", err)
		}
		for _, object := range names {
			if strings.HasSuffix(object, "/") {
				recursiveObjects := getSwiftObjectNames(t, osClient, container, prefix+object)
				objectNames = append(objectNames, recursiveObjects...)
				continue
			}
			objectNames = append(objectNames, object)
		}

		return true, nil
	})
	return
}

// Helper function to create Swift container
func createSwiftContainer(t *testing.T, osClient *gophercloud.ServiceClient, container string) {
	container, _ = getContainerAndPrefix(container)

	// Create the container
	createResult := containers.Create(osClient, container, &containers.CreateOpts{})
	if createResult.Err != nil {
		t.Fatalf("Error creating %s container: %s", container, createResult.Err)
	}
}

// Helper function to delete Swift container
func deleteSwiftContainer(t *testing.T, osClient *gophercloud.ServiceClient, container string) {
	container, prefix := getContainerAndPrefix(container)

	warning := "WARNING: Failed to delete the test Swift container. It may have been left in your Openstack account and may incur storage charges. (error was %s)"

	// Remove any objects
	deleteSwiftObjects(t, osClient, container, prefix)

	// Delete the container
	deleteResult := containers.Delete(osClient, container)
	if deleteResult.Err != nil {
		if _, ok := deleteResult.Err.(gophercloud.ErrDefault404); !ok {
			t.Fatalf(warning, deleteResult.Err)
		}
	}
}

// Helper function to delete Swift objects within a container
func deleteSwiftObjects(t *testing.T, osClient *gophercloud.ServiceClient, container, prefix string) {
	objectNames := getSwiftObjectNames(t, osClient, container, prefix)

	// Get a slice of object names
	for _, object := range objectNames {
		result := objects.Delete(osClient, container, object, nil)
		if result.Err != nil {
			if _, ok := result.Err.(gophercloud.ErrDefault404); !ok {
				t.Fatalf("Error deleting object %s from container %s: %s", object, container, result.Err)
			}
		}
	}
}
