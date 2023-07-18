package test

import (
	"offstorage/data"
	"testing"
)

// executer
func TestDataAddGet(t *testing.T) {
	t.Log("begin")
	dsh, _, err := data.NewDataShell(
		data.DataWithHost("127.0.0.1"),
		data.DataWithPort(5001),
		data.DataWithKeyName("data"),
		data.DataWithIpnsName("k51qzi5uqu5dki16z7czuv0wjv47dpl112v5pqhiifqzkcje9dt42lvdk9paif"),
		data.DataWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/test_get_data"),
		data.DataWithRole("executer"),
	)

	if err != nil {
		t.Error(err)
	}

	err = dsh.AddData("test", "test_value")
	if err != nil {
		t.Error(err)
	}

	data, err := dsh.GetData("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("data: %s", data)

	err = dsh.CloseImage()
	if err != nil {
		t.Error(err)
	}
}

// verifier
func TestVDataAddGet(t *testing.T) {
	t.Log("begin")
	dsh, _, err := data.NewDataShell(
		data.DataWithHost("127.0.0.1"),
		data.DataWithPort(5001),
		data.DataWithKeyName("data"),
		data.DataWithIpnsName("k51qzi5uqu5dki16z7czuv0wjv47dpl112v5pqhiifqzkcje9dt42lvdk9paif"),
		data.DataWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/test_get_data"),
		data.DataWithRole("verifier"),
		data.DataWithTaskID("4c5f23da-67aa-4b5e-8b02-7087d99e831d"),
	)

	if err != nil {
		t.Error(err)
	}

	err = dsh.AddData("test", "test_value")
	if err != nil {
		t.Error(err)
	}

	data, err := dsh.GetData("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("data: %s", data)

	err = dsh.CloseImage()
	if err != nil {
		t.Error(err)
	}
}
