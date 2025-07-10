package device

import (
	"os/exec"
	"fmt"
	"bufio"
	"bytes"
	"strings"
)
type MyNetwork struct{
	BSSID string
	SSID  string
}


func GetMyCurrentNetwork() (*MyNetwork, error) {
	cmd := exec.Command("ipconfig", "getsummary", "en0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run ipconfig: %w", err)
	}

	var ssid, bssid string

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SSID") {
			// Example: "SSID                 : MyNetwork"
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				ssid = strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "BSSID") {
			// Example: "BSSID                : 12:34:56:78:90:ab"
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				bssid = strings.TrimSpace(parts[1])
			}
		}
	}

	return &MyNetwork{BSSID: bssid, SSID: ssid}, nil
}
