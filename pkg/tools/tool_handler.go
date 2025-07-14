package tools

import (
	"fmt"
	"io/fs"
  "path/filepath"
	"plugin"
)


func GetAvailableToolPacks() []string {
    rootDir := "./configs/user/tools"
    var soFiles []string

    filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            fmt.Printf("error accessing %q: %v\n", path, err)
            return nil
        }
        if !d.IsDir() && filepath.Ext(path) == ".so" {
            soFiles = append(soFiles, path)
        }
        return nil
    })

	return soFiles
}

func CheckOutToolPack() {
	rootDir := "./configs/user/tools"
	toolPath := filepath.Join(rootDir, "reaper_tools.so")
	plug, err:= plugin.Open(toolPath)

	if err != nil {
		fmt.Println(err)
			return
	}

	sym, err:= plug.Lookup("PluginSpecs")
	if err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println(sym)

  specsFunc, ok := sym.(func() []Tool)
	if !ok{
		fmt.Println("yikes")
		return
	}

	for _, spec:= range specsFunc(){
		fmt.Println(spec)
		fmt.Println(spec.Name)
	}
}



