package bbrpc

import (
	"strings"
	"testing"
)

func TestClient_Version(t *testing.T) {
	args := DefaultDebugBBArgs()
	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      args,
	})
	tShouldNil(t, err)
	defer killBigBangServer()

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err)

	ver, err := client.Version()
	tShouldNil(t, err)
	tShouldTrue(t, strings.Contains(ver, "."))
	tShouldTrue(t, strings.Contains(ver, "v"))
}
