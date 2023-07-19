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
		data.DataWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/test_data"),
		data.DataWithRole("executer"),
	)

	if err != nil {
		t.Error(err)
	}

	// err = dsh.AddData("test", "test_value3")
	// if err != nil {
	// 	t.Error(err)
	// }

	data, err := dsh.GetData("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("data: %s", data)

	err = dsh.CloseData()
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
		data.DataWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/test_data"),
		data.DataWithRole("verifier"),
		data.DataWithTaskID("771dceb6-a069-46ae-90e6-c5818bab91ba"),
	)

	if err != nil {
		t.Error(err)
	}

	// err = dsh.AddData("test", "test_value1")
	// if err != nil {
	// 	t.Error(err)
	// }

	data, err := dsh.GetData("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("data: %s", data)

	err = dsh.CloseData()
	if err != nil {
		t.Error(err)
	}
}

func TestDBPersistant(t *testing.T) {
	t.Log("begin")
	dsh, _, err := data.NewDataShell(
		data.DataWithHost("127.0.0.1"),
		data.DataWithPort(5001),
		data.DataWithKeyName("data"),
		data.DataWithIpnsName("k51qzi5uqu5dki16z7czuv0wjv47dpl112v5pqhiifqzkcje9dt42lvdk9paif"),
		data.DataWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/test_data"),
		data.DataWithRole("judge"),
		data.DataWithTaskID("b5427047-d2c2-499f-ad64-f92425057a53"),
	)
	if err != nil {
		t.Error(err)
	}

	// text, err := dsh.TraverseTable("8645d145-3dec-44c4-8f0b-1fbfaf3334ba", "")
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(text)

	err = dsh.DataPersistance()
	if err != nil {
		t.Error(err)
	}

	err = dsh.CloseJudgeData()
	if err != nil {
		t.Error(err)
	}
}
