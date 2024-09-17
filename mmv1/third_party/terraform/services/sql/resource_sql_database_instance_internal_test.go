package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func TestMaintenanceVersionDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New       string
		ShouldSuppress bool
	}{
		"older configuration maintenance version than current version should suppress diff": {
			Old:            "MYSQL_8_0_26.R20220508.01_09",
			New:            "MYSQL_5_7_37.R20210508.01_03",
			ShouldSuppress: true,
		},
		"older configuration maintenance version than current version should suppress diff with lexicographically smaller database version": {
			Old:            "MYSQL_5_8_10.R20220508.01_09",
			New:            "MYSQL_5_8_7.R20210508.01_03",
			ShouldSuppress: true,
		},
		"newer configuration maintenance version than current version should not suppress diff": {
			Old:            "MYSQL_5_7_37.R20210508.01_03",
			New:            "MYSQL_8_0_26.R20220508.01_09",
			ShouldSuppress: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if maintenanceVersionDiffSuppress("version", tc.Old, tc.New, nil) != tc.ShouldSuppress {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", tc.Old, tc.New, tc.ShouldSuppress)
			}
		})
	}
}

type updateData struct {
	wantInstance *sqladmin.DatabaseInstance
	op           *sqladmin.Operation
}

var (
	mockConfig = &transport_tpg.Config{UserAgent: "unittest-user-agent"}

	switchoverOpData = map[string]*sqladmin.Operation{}
	updateOpData     = map[string]updateData{}

	instanceData = map[string]*sqladmin.DatabaseInstance{}
	opWaitData   = map[*sqladmin.Operation]*sqladmin.DatabaseInstance{}
)

type savedFuncs struct {
	genUserAgent   func(d *schema.ResourceData, userAgent string) (string, error)
	getProject     func(d *schema.ResourceData, config *transport_tpg.Config) (string, error)
	instancesGet   func(config *transport_tpg.Config, userAgent, project, instanceName string) (*sqladmin.DatabaseInstance, error)
	update         func(config *transport_tpg.Config, userAgent, project, instanceName string, instance *sqladmin.DatabaseInstance) (*sqladmin.Operation, error)
	patch          func(config *transport_tpg.Config, userAgent, project, instanceName string, instance *sqladmin.DatabaseInstance) (*sqladmin.Operation, error)
	promoteReplica func(config *transport_tpg.Config, userAgent, project, instanceName string) (*sqladmin.Operation, error)
	switchover     func(config *transport_tpg.Config, userAgent, project, instanceName string) (*sqladmin.Operation, error)
	restoreBackup  func(config *transport_tpg.Config, userAgent, project, instanceName string, backupRequest *sqladmin.InstancesRestoreBackupRequest) (*sqladmin.Operation, error)
	usersUpdate    func(config *transport_tpg.Config, userAgent, project, instanceName string, user *sqladmin.User, host, name string) (*sqladmin.Operation, error)
	opWait         func(config *transport_tpg.Config, opqOp interface{}, project, activity, userAgent string, timeout time.Duration) error
}

func saveOrigFuncs() savedFuncs {
	saveData := savedFuncs{}

	saveData.genUserAgent = generateUserAgentString
	saveData.getProject = getProject
	saveData.instancesGet = instancesGet
	saveData.update = instancesUpdate
	saveData.patch = instancesPatch
	saveData.promoteReplica = instancesPromoteReplica
	saveData.switchover = instancesSwitchover
	saveData.restoreBackup = instancesRestoreBackup
	saveData.usersUpdate = usersUpdate
	saveData.opWait = sqlAdminOperationWaitTime

	return saveData
}

func installMockFuncs() {
	generateUserAgentString = mockGenUserAgent
	getProject = mockGetProject
	instancesGet = mockInstancesGet
	instancesUpdate = mockUpdate
	instancesPatch = mockPatch
	instancesPromoteReplica = mockPromoteReplica
	instancesSwitchover = mockSwitchover
	instancesRestoreBackup = mockRestoreBackup
	usersUpdate = mockUsersUpdate
	sqlAdminOperationWaitTime = mockOpWait
}

func restoreFuncsFrom(saveData savedFuncs) {
	generateUserAgentString = saveData.genUserAgent
	getProject = saveData.getProject
	instancesGet = saveData.instancesGet
	instancesUpdate = saveData.update
	instancesPatch = saveData.patch
	instancesPromoteReplica = saveData.promoteReplica
	instancesSwitchover = saveData.switchover
	instancesRestoreBackup = saveData.restoreBackup
	usersUpdate = saveData.usersUpdate
	sqlAdminOperationWaitTime = saveData.opWait
}

func mockGenUserAgent(d *schema.ResourceData, userAgent string) (string, error) {
	return userAgent, nil
}

func mockGetProject(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	return "unittest-project", nil
}

func mockSwitchover(_ *transport_tpg.Config, _, _, instanceName string) (*sqladmin.Operation, error) {
	op, ok := switchoverOpData[instanceName]
	if !ok {
		return nil, fmt.Errorf("unexpected Switchover called for instance %q", instanceName)
	}
	delete(switchoverOpData, instanceName)
	return op, nil
}

func mockUpdate(_ *transport_tpg.Config, _, _, instanceName string, instance *sqladmin.DatabaseInstance) (*sqladmin.Operation, error) {
	opData, ok := updateOpData[instanceName]
	if !ok {
		var inToPrint interface{} = instance
		if instance != nil {
			jsonData, err := json.Marshal(*instance)
			if err == nil {
				inToPrint = string(jsonData)
			}
		}

		return nil, fmt.Errorf("unexpcted update called for instance %q, new instance: %+v", instanceName, inToPrint)
	}
	if opData.wantInstance != nil {
		if diff := cmp.Diff(opData.wantInstance, instance); diff != "" {
			return nil, fmt.Errorf("unexpected instance diff (-want, got):\n%s", diff)
		}
	}
	delete(updateOpData, instanceName)
	return opData.op, nil
}

func mockInstancesGet(_ *transport_tpg.Config, _, _, instanceName string) (*sqladmin.DatabaseInstance, error) {
	instance, ok := instanceData[instanceName]
	if !ok {
		return nil, fmt.Errorf("unexpected Get called for instance %q", instanceName)
	}
	return instance, nil
}

func mockOpWait(_ *transport_tpg.Config, opqOp interface{}, _, _, _ string, _ time.Duration) error {
	op := opqOp.(*sqladmin.Operation)
	newInstance, ok := opWaitData[op]
	if !ok {
		return fmt.Errorf("waiting for unexpected operation %+v", op)
	}
	if newInstance != nil {
		instanceData[op.TargetId] = newInstance
	}
	delete(opWaitData, op)
	return nil
}

func mockPromoteReplica(_ *transport_tpg.Config, _, _, _ string) (*sqladmin.Operation, error) {
	return nil, fmt.Errorf("mock for promote replica unimplemented")
}

func mockPatch(_ *transport_tpg.Config, _, _, _ string, _ *sqladmin.DatabaseInstance) (*sqladmin.Operation, error) {
	return nil, fmt.Errorf("mock for patch unimplemented")
}

func mockRestoreBackup(_ *transport_tpg.Config, _, _, _ string, _ *sqladmin.InstancesRestoreBackupRequest) (*sqladmin.Operation, error) {
	return nil, fmt.Errorf("mock restore backup unimpelmented")
}

func mockUsersUpdate(_ *transport_tpg.Config, _, _, _ string, _ *sqladmin.User, _, _ string) (*sqladmin.Operation, error) {
	return nil, fmt.Errorf("mock users update unimplemented")
}

// func TestSwitchover(t *testing.T) {
// 	saveData := saveOrigFuncs()
// 	installMockFuncs()
// 	defer restoreFuncsFrom(saveData)

// 	ctx := context.Background()

// 	state := map[string]interface{}{
// 		"region":               "us-central1",
// 		"project":              "unittest-project",
// 		"database_version":     "MYSQL_8_0",
// 		"name":                 "original-replica",
// 		"instance_type":        "READ_REPLICA_INSTANCE",
// 		"master_instance_name": "original-primary",
// 		"settings": []interface{}{
// 			map[string]interface{}{
// 				"version": 1,
// 				"tier":    "db-n1-standard-1",
// 			},
// 		},
// 		"replica_configuration": []interface{}{
// 			map[string]interface{}{
// 				"failover_target": false,
// 			},
// 		},
// 		"replica_names": []interface{}{
// 			"leaf-replica-1", "leaf-replica-2",
// 		},
// 	}

// 	diff := &terraform.InstanceDiff{
// 		Attributes: map[string]*terraform.ResourceAttrDiff{
// 			"master_instance_name": &terraform.ResourceAttrDiff{
// 				NewRemoved: true,
// 			},
// 			"instance_type": &terraform.ResourceAttrDiff{
// 				Old: "READ_REPLICA_INSTANCE",
// 				New: "CLOUD_SQL_INSTANCE",
// 			},
// 			"replica_names.#": &terraform.ResourceAttrDiff{
// 				Old: "2",
// 				New: "3",
// 			},
// 			"replica_names.2": &terraform.ResourceAttrDiff{
// 				New: "original-primary",
// 			},
// 		},
// 	}

// 	rd, err := testResourceData(ctx, t, state, diff)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rd.Set("master_instance_name", nil)
// 	rd.Set("replica_configuration", nil)

// 	t.Logf("LYCH DEBUG replica config:\n%+v", rd.Get("replica_configuration"))
// 	t.Logf("LYCH DEBUG mastername:\n%+v", rd.Get("master_instance_name"))
// 	om, nm := rd.GetChange("master_instance_name")
// 	t.Logf("LYCH DEBUG mastername:\n%+v %+v", om, nm)
// 	t.Logf("LYCH DEBUG it:\n%+v", rd.Get("instance_type"))
// 	oit, nit := rd.GetChange("instance_type")
// 	t.Logf("LYCH DEBUG it:\n%+v %+v", oit, nit)

// 	rd.Set("master_instance_name", nil)
// 	rd.Set("replica_configuration", nil)

// 	t.Logf("LYCH DEBUG replica config:\n%+v", rd.Get("replica_configuration"))
// 	t.Logf("LYCH DEBUG mastername:\n%+v", rd.Get("master_instance_name"))
// 	om, nm = rd.GetChange("master_instance_name")
// 	t.Logf("LYCH DEBUG mastername:\n%+v %+v", om, nm)
// 	t.Logf("LYCH DEBUG it:\n%+v", rd.Get("instance_type"))
// 	oit, nit = rd.GetChange("instance_type")
// 	t.Logf("LYCH DEBUG it:\n%+v %+v", oit, nit)

// 	fakeSwitchoverOp := &sqladmin.Operation{}
// 	switchoverOpData["original-replica"] = fakeSwitchoverOp
// 	opWaitData[fakeSwitchoverOp] = nil
// 	fakeUpdateOp := &sqladmin.Operation{}
// 	updateOpData["original-replica"] = updateData{op: fakeUpdateOp}
// 	opWaitData[fakeUpdateOp] = nil

// 	instanceData["original-replica"] = &sqladmin.DatabaseInstance{
// 		Name: "original-replica",
// 		Settings: &sqladmin.Settings{
// 			SettingsVersion: 2,
// 		},
// 	}

// 	if err := resourceSqlDatabaseInstanceUpdate(rd, mockConfig); err != nil {
// 		t.Fatal(err)
// 	}
// 	if _, ok := switchoverOpData["original-replica"]; ok {
// 		t.Fatal("switchover didn't called")
// 	}
// }

var schemaMap = schema.InternalMap(ResourceSqlDatabaseInstance().SchemaMap())

func testResourceData(ctx context.Context, t *testing.T, state map[string]interface{}, diffState *terraform.InstanceDiff) (*schema.ResourceData, error) {
	rd1 := schema.TestResourceDataRaw(t, schemaMap, state)
	rd1.SetId(fmt.Sprintf("unit-test-id-%d", rand.Uint32()))
	state1 := rd1.State()

	// diffState, err := schemaMap.Diff(ctx, state1, diffCfg, nil, nil, false)
	// if err != nil {
	// 	return nil, err
	// }
	ret, err := schemaMap.Data(state1, diffState)
	if err != nil {
		return nil, err
	}
	ret.SetId(fmt.Sprintf("unit-test-id-%d", rand.Uint32()))
	return ret, nil

}

func copyState(t *testing.T, src map[string]interface{}) map[string]interface{} {
	if src == nil {
		t.Fatalf("source state is nil")
	}

	dst := map[string]interface{}{}
	for k, v := range src {
		dst[k] = copyVal(t, v)
	}
	return dst
}

func copyVal(t *testing.T, val interface{}) interface{} {
	if val == nil {
		t.Fatalf("value is nil")
	}
	switch val.(type) {
	case bool, int, string:
		return val
	case []interface{}:
		newVal := []interface{}{}
		for _, listElt := range val.([]interface{}) {
			newVal = append(newVal, copyVal(t, listElt))
		}
		return newVal
	case map[string]interface{}:
		newVal := map[string]interface{}{}
		for k, innerVal := range val.(map[string]interface{}) {
			newVal[k] = copyVal(t, innerVal)
		}
		return newVal
	default:
		t.Fatalf("value with unknown type %T: %+v", val, val)
		return nil
	}
}
