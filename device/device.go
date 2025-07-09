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

func GetCurrentAudioDevice() (string, string) {
    err := portaudio.Initialize()
    if err != nil {
        log.Fatal(err)
    }
    defer portaudio.Terminate()

    inputDev, err := portaudio.DefaultInputDevice()
    if err != nil {
        log.Fatal(err)
    }
    outputDev, err := portaudio.DefaultOutputDevice()
    if err != nil {
        log.Fatal(err)
    }
    return inputDev.Name, outputDev.Name
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
