package device

import ( 
	"fmt"
	"log"
	"github.com/gordonklaus/portaudio"
)

func GetAudioDevices() ([]*portaudio.DeviceInfo, error) {
    if err := portaudio.Initialize(); err != nil {
        return nil, err
    }
    defer portaudio.Terminate()

    devices, err := portaudio.Devices()
    if err != nil {
        return nil, err
    }

    return devices, nil
}

func GetAudioOutputDevices() ([]*portaudio.DeviceInfo, error) {
    devices, err := GetAudioDevices()
    if err != nil {
        return nil, err
    }

		fmt.Println("outputDevices")
    var outputs []*portaudio.DeviceInfo
    for _, d := range devices {
        if d.MaxOutputChannels > 0 {
            outputs = append(outputs, d)
        }
    }
    return outputs, nil
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



func PrintDevices(devices []*portaudio.DeviceInfo) {
    for i, d := range devices {
        api := "unknown"
        if d.HostApi != nil {
            api = d.HostApi.Name
        }
        fmt.Printf(
            "[%d] Name: %s, API: %s, In: %d, Out: %d\n",
            i, d.Name, api, d.MaxInputChannels, d.MaxOutputChannels,
        )
    }
}
