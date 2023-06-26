package test

import (
	"bytes"
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
