package util

import "os"

// GetExecDir return absolute directory by invoker
func GetExecDir() (string, error) {
	file, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return file, nil
}
