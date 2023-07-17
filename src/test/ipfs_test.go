package test

import (
	"bytes"
	"fmt"
	"offstorage/ipfs"
	"strings"
	"testing"
)

func TestIpfsInit(t *testing.T) {
	api, err := ipfs.NewShell(
		ipfs.ShellWithHost("127.0.0.1"),
		ipfs.ShellWithPort(5001),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cid, err := api.Sh.Add(strings.NewReader("12345"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(cid)

	content, err := api.Sh.Cat(cid)
	if err != nil {
		t.Fatalf(err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	finalStr := buf.String()
	t.Logf(finalStr)
}

func TestIpfsPublish(t *testing.T) {
	api, err := ipfs.NewShell(
		ipfs.ShellWithHost("127.0.0.1"),
		ipfs.ShellWithPort(5001),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cid, err := api.AddString("{}")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(cid)

	ipnsname, err := api.PublishFile(cid, "test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(ipnsname)

	content, err := api.Sh.Cat(
		fmt.Sprintf("/ipns/%s", ipnsname),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	finalStr := buf.String()
	t.Logf(finalStr)
}

func TestIpfsPublishFolder(t *testing.T) {
	api, err := ipfs.NewShell(
		ipfs.ShellWithHost("127.0.0.1"),
		ipfs.ShellWithPort(5001),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cid, err := api.AddFolder("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/empty_data")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(cid)

	ipnsname, err := api.PublishFile(cid, "data")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(ipnsname)

	content, err := api.Sh.Cat(
		fmt.Sprintf("/ipns/%s/db", ipnsname),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	finalStr := buf.String()
	t.Logf(finalStr)
}
