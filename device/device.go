package device

import ( 
	"fmt"
	"log"
	"github.com/gordonklaus/portaudio"
)

func GetAudioDevice() { 
	err := portaudio.Initialize() 
	if err != nil { 
		log.Fatal(err) 
	}

	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Printf("Name: %s, Max Input Channels: %d, Max Output Channels: %d\n", device.Name, device.MaxInputChannels, device.MaxOutputChannels)
	}

}

func GetCurrentAudioDevice(){
	err := portaudio.Initialize() 
	if err != nil { 
		log.Fatal(err) 
	}

	inputDev, err := portaudio.DefaultInputDevice()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Default input: %s\n", inputDev.Name)
	// Get default output device
	outputDev, err := portaudio.DefaultOutputDevice()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Default output: %s\n", outputDev.Name)
}
