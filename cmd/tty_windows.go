// +build windows
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func CmdStart(cmd *exec.Cmd) (*os.File, error) {
	readerCloser, _ := cmd.StdoutPipe()
	all, err := ioutil.ReadAll(readerCloser)
	file, err := ioutil.TempFile(os.TempDir(), "output")
	fmt.Println(file.Name())
	ioutil.WriteFile(file.Name(), all, 0644)
	return file, err
}
