package main

// import "fmt"
import "time"
import "os/exec"
import "log"
import "flag"

import "strings"

func main() {
	// TODO move arguments parsing in fuction
	var timeout = flag.String("timeout", "2s", "timeout for checks, 2 seconds by default")

	flag.Parse()

	sleep, err := time.ParseDuration(*timeout)

	if err != nil {
		log.Fatalln("Incorrect usage of --timeout param; Pls use like --timeout=5s")
	}

	if sleep < time.Duration(time.Second) {
		log.Fatal("Timeout can't be less than 1s")
	}

	cmd := "networksetup -listallhardwareports | fgrep Wi-Fi -A1 | awk 'NF==2{print $2}'"
	eth, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		log.Fatal("Can't detect wifi interface")
	}

	log.Printf("Starting checking interface %s with timeout %s", strings.Trim(string(eth), "\n "), sleep.String())

	for {

		enabled, err := exec.Command("networksetup", "getairportpower", strings.Trim(string(eth), "\n ")).Output()

		if strings.Contains(strings.Trim(string(enabled), "\n "), "On") {

			err = exec.Command("ping", "-c", "2", "8.8.8.8").Run()

			if err != nil {
				restartWifi(string(eth))
			}

			time.Sleep(sleep)
		}
	}
}

func restartWifi(eth string) {
	log.Println("Wifi restarting")
	exec.Command("networksetup", "-setairportpower", eth, "off").Run()
	time.Sleep(2 * time.Second)
	exec.Command("networksetup", "-setairportpower", eth, "on").Run()
}
