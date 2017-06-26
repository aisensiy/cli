// +build !windows

package cmd

import (
	"github.com/kr/pty"
	"os"
	"os/exec"
)

func CmdStart(cmd *exec.Cmd) (*os.File, error) {
	return pty.Start(cmd)
}
