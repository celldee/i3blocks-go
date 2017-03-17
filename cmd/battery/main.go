package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"io/ioutil"
)

func main() {

	// Set display texts to defaults.
	var icon string
	var output string
	var fullText string = "error"
	var shortText string = "error"

	// Read charging status information from kernel
	// pseudo-file-system mounted at /sys.
	chargingRaw, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {

		// Write an error to STDERR, fallback display values
		// to STDOUT and exit with failure code.
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read charging file: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(1)
	}

	// Read full capacity information from /sys.
	chargeFullRaw, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/charge_full")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read full capacity file: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(1)
	}

	// Read current capacity information from /sys.
	chargeNowRaw, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/charge_now")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Failed to read current capacity file: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(1)
	}

	// Trim whitespace.
	charging := strings.TrimSpace(string(chargingRaw))
	chargeFullString := strings.TrimSpace(string(chargeFullRaw))
	chargeNowString := strings.TrimSpace(string(chargeNowRaw))

	// Convert full capacity string to float32.
	chargeFull, err := strconv.ParseFloat(chargeFullString, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Could not convert full capacity value: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(1)
	}

	// Convert current capacity string to float32.
	chargeNow, err := strconv.ParseFloat(chargeNowString, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-go] Could not convert current capacity value: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(1)
	}

	// Calculate current battery charge percentage.
	chargePerc := int((chargeNow / chargeFull) * 100.0)

	// Avoid unreasonable rounding or reporting errors.
	if chargePerc > 100 {
		chargePerc = 100
	}

	if (chargePerc < 8) && (charging != "Charging") {

		// If charge percentage is very low and battery
		// is currently not being charged, print values
		// and exit with urgent return code.
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(33)
	}

	// Depending on current charge percentage,
	// set appropriate battery icon.
	if (chargePerc >= 0) && (chargePerc <= 20) {
		icon = ""
	} else if (chargePerc >= 21) && (chargePerc <= 40) {
		icon = ""
	} else if (chargePerc >= 41) && (chargePerc <= 60) {
		icon = ""
	} else if (chargePerc >= 61) && (chargePerc <= 80) {
		icon = ""
	} else {
		icon = ""
	}

	// Construct and color final output string based
	// on charging status and percentage.
	if charging == "Charging" {
		output = fmt.Sprintf("<span foreground=\"#378c1a\">%s</span>%4d%%", icon, chargePerc)
	} else {

		if (chargePerc >= 0) && (chargePerc <= 20) {
			output = fmt.Sprintf("<span foreground=\"#ff0000\">%s</span>%4d%%", icon, chargePerc)
		} else if (chargePerc >= 21) && (chargePerc <= 30) {
			output = fmt.Sprintf("<span foreground=\"#ffae00\">%s</span>%4d%%", icon, chargePerc)
		} else if (chargePerc >= 31) && (chargePerc <= 40) {
			output = fmt.Sprintf("<span foreground=\"#fff600\">%s</span>%4d%%", icon, chargePerc)
		} else {
			output = fmt.Sprintf("%s%4d%%", icon, chargePerc)
		}
	}

	fullText = output
	shortText = output

	// Write out gathered information to STDOUT.
	fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
	os.Exit(0)
}
