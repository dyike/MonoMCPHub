package tools

import "os/exec"

func ExecuteAdbCommand(args []string) (string, error) {
	cmd := exec.Command("adb", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
