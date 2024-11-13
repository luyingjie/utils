package cmd

import (
	"os/exec"
	"strings"
)

func Cmd(cm string) error {
	if cm == "" {
		return nil
	}
	// cmd := exec.Command("bash", "-c", cm)
	// if err := cmd.Run(); err != nil {
	// 	return err
	// }
	// return nil
	cmd := exec.Command("bash", "-c", cm)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func CmdToStr(cm string) (string, error) {
	if cm == "" {
		return "", nil
	}
	out, err := exec.Command("bash", "-c", cm).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func CmdToArr(cm string) ([]string, error) {
	if cm == "" {
		return nil, nil
	}
	out, err := exec.Command("bash", "-c", cm).Output()
	if err != nil {
		return nil, err
	}

	outs := strings.Split(string(out), "\n")
	return outs[0 : len(outs)-2], nil
}
