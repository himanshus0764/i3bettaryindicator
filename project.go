package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func zenity(message string) {
	cmd := exec.Command("zenity", "--warning", fmt.Sprintf("--text=%s", message))
	err := cmd.Run()
	if err != nil {
		fmt.Println("zenity error:", err)
	}
}

func getChargingStatus() int {
	status := exec.Command("cat", "/sys/class/power_supply/ACAD/online")
	output, err := status.Output()
	if err != nil {
		fmt.Println("Charging status error:", err)
		return -1
	}
	result := strings.TrimSpace(string(output))
	val, err := strconv.Atoi(result)
	if err != nil {
		fmt.Println("Conversion error (charging status):", err)
		return -1
	}
	return val
}

func getBatteryPercentage() int {
	battery := exec.Command("bash", "-c", "upower -i $(upower -e | grep BAT) | grep -E 'percentage'")
	output, err := battery.Output()
	if err != nil {
		fmt.Println("Battery percentage error:", err)
		return -1
	}
	line := string(output)
	percent := strings.TrimSpace(strings.Split(line, ":")[1])
	percent = strings.TrimSuffix(percent, "%")
	val, err := strconv.Atoi(strings.TrimSpace(percent))
	if err != nil {
		fmt.Println("Conversion error (battery percentage):", err)
		return -1
	}
	return val
}

func main() {
	for {
		chargingStatus := getChargingStatus()
		batteryPercentage := getBatteryPercentage()

		if chargingStatus == -1 || batteryPercentage == -1 {
			zenity("error to execute commant cause to may be wrong command string check and fix it according to you machine")
		} else {
			if chargingStatus == 0 {
				switch {
				case batteryPercentage <= 5:
					zenity("Battery in critical stage, charge it.")
				case batteryPercentage == 10:
					zenity("Put battery on charge. Only 10% remaining.")
				case batteryPercentage <= 20:
					zenity("Battery low, please charge.")
				}
			} else if chargingStatus == 1 && batteryPercentage >= 95 && batteryPercentage <= 100 {
				zenity("Battery full, remove charger.")
			}
		}
	}
}
