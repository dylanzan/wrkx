package ostools

import (
	"fmt"
	"sshx/internal/cores/sshx"
	"testing"
)

func TestGetOsersion(t *testing.T) {

	sshx, err := sshx.NewSSHX("", "", "")
	if err != nil {
		t.Error(err)
	}

	osversion, err := GetOSVersionName(sshx)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(osversion)
}
