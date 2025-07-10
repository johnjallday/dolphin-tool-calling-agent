package tools

import (
	"fmt"
	"os/exec"
)

func BuildPlugin(){
	cmd :=exec.Command("go build -buildmode=plugin -o ./user/tools/reaper_rtemplate_loader.so ../dolphin-reaper-track-manager/main.go/")
//go build -buildmode=plugin -o ./reaper_rtemplate_loader.so ./main.go
	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output
	fmt.Println(string(output))

}

