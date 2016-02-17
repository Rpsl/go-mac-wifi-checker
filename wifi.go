package main

// import "fmt"
import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go Wifi checker"
	app.Usage = ""
	app.HideVersion = true

	var timeout string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "timeout, t",
			Value:       "2s",
			Usage:       "timeout between checks",
			Destination: &timeout,
		},
	}

	app.Action = func(c *cli.Context) {
		loop(timeout)
	}

	app.Run(os.Args)
}

func loop(timeout string) {

	sleep := checkParams(timeout)
	eth := getWifiInterface()

	log.Printf("Start checking interface %s with timeout %s", eth, sleep.String())

	for {

		enabled, err := exec.Command("networksetup", "getairportpower", eth).Output()

		if strings.Contains(strings.Trim(string(enabled), "\n "), "On") {

			err = exec.Command("ping", "-c", "2", "8.8.8.8").Run()

			if err != nil {
				restartWifi(string(eth))
			}

			time.Sleep(sleep)
		}
	}

}

func checkParams(timeout string) time.Duration {
	sleep, err := time.ParseDuration(timeout)

	if err != nil {
		log.Fatalln("Incorrect usage of --timeout param; Pls use like --timeout=5s")
	}

	if sleep < time.Duration(time.Second) {
		log.Fatal("Timeout can't be less than 1s")
	}

	return sleep
}

func getWifiInterface() string {
	cmd := "networksetup -listallhardwareports | fgrep Wi-Fi -A1 | awk 'NF==2{print $2}'"
	eth, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		log.Fatal("Can't detect wifi interface")
	}

	return strings.Trim(string(eth), " \n")
}

func restartWifi(eth string) {
	log.Println("Wifi restarting")
	exec.Command("networksetup", "-setairportpower", eth, "off").Run()
	time.Sleep(2 * time.Second)
	exec.Command("networksetup", "-setairportpower", eth, "on").Run()
}
