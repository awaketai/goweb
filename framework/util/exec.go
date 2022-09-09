package util

import (
	"os"
	"syscall"
)

// GetExecDir return absolute directory by invoker
func GetExecDir() (string, error) {
	file, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return file, nil
}

func CheckProcessExists(pid int) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, err
	}
	return true, nil
}
