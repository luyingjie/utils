package cmd

import (
	"fmt"
	"testing"
)

func TestCMD(t *testing.T) {
	// cmd := exec.Command("ls", "-l", "-a")
	// if err := cmd.Run(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(cmd)
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// cmd := exec.Command("bash", "-c", "ifconfig")
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	// if err := cmd.Run(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(stdout.String())
	// fmt.Println("------------")
	// fmt.Println(stderr.String())
	// out, err1 := exec.Command("bash", "-c", "ls").CombinedOutput() //.Output()
	// if err1 != nil {
	// 	fmt.Println(err1)
	// 	return
	// }
	// fmt.Println(string(out))
	pu, _ := CmdToArr("ls")
	fmt.Println(len(pu))
	fmt.Println(pu)
}
