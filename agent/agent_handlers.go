package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

func ListAgents() {
	fmt.Println("Get List of Agents in the file")
	files, err := os.ReadDir("./user/agents")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".toml" {
			filePath := filepath.Join("./user/agents", f.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", f.Name(), err)
				continue
			}
			fmt.Printf("----- %s -----\n", f.Name())
			fmt.Println(string(data))
			fmt.Println("---------------------")
		}
	}
}
