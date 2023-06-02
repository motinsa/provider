package esxi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func getESXiClient(hostname string, username string, password string) (*govmomi.Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url, err := url.Parse(fmt.Sprintf("https://%s/sdk", hostname))
	if err != nil {
		return nil, err
	}
	url.User = url.UserPassword(username, password)

	client, err := govmomi.NewClient(ctx, url, true)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createVM(esxiHostname string, esxiUsername string, esxiPassword string, guestName string, memSize int64, numVCPUs int32, guestOS string) error {
	client, err := getESXiClient(esxiHostname, esxiUsername, esxiPassword)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(client.Client, true)

	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return err
	}

	finder.SetDatacenter(dc)

	cr, err := finder.DefaultResourcePool(ctx)
	if err != nil {
		return err
	}

	var mds mo.Datastore
	ds, err := finder.DefaultDatastore(ctx)
	if err != nil {
		return err
	}

	err = ds.Properties(ctx, ds.Reference(), []string{"info"}, &mds)
	if err != nil {
		return err
	}

	spec := types.VirtualMachineConfigSpec{
		Name:     guestName,
		GuestId:  guestOS,
		Files:    &types.VirtualMachineFileInfo{VmPathName: fmt.Sprintf("[%s]", mds.Info.GetDatastoreInfo().Name)},
		NumCPUs:  numVCPUs,
		MemoryMB: memSize,
	}

	folders, err := dc.Folders(ctx)
	if err != nil {
		return err
	}

	task, err := folders.VmFolder.CreateVM(ctx, spec, cr, nil)
	if err != nil {
		return err
	}

	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return err
	}

	if info.State != types.TaskInfoStateSuccess {
		return fmt.Errorf("createVM task did not complete successfully")
	}

	return nil
}
