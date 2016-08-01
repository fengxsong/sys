package sys

import (
	"bytes"
	"os/exec"
)

func ExecCmdDirBytes(dir, cmdName string, args ...string) ([]byte, []byte, error) {
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	cmd.Stdout = bufOut
	cmd.Stderr = bufErr

	err := cmd.Run()
	return bufOut.Bytes(), bufErr.Bytes(), err
}

func ExecCmdBytes(cmdName string, args ...string) ([]byte, []byte, error) {
	return ExecCmdDirBytes("", cmdName, args...)
}

func ExecCmdDir(dir, cmdName string, args ...string) (string, string, error) {
	bufOut, bufErr, err := ExecCmdDirBytes(dir, cmdName, args...)
	return string(bufOut), string(bufErr), err
}

func ExecCmd(cmdName string, args ...string) (string, string, error) {
	return ExecCmdDir("", cmdName, args...)
}
