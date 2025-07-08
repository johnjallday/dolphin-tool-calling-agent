package main


import(
	"github.com/johnjallday/dolphin-tool-calling-agent/tools"
	"fmt"
)


func main(){
	soFiles := tools.GetAvailableToolPacks()
	fmt.Println(soFiles)

	tools.CheckOutToolPack()

}
