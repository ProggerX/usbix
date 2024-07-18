package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Input []struct {
	Device     string `json:"device"`
	Type       string `json:"type"`
	Partitions []struct {
		Label string `json:"label"`
		Size  int    `json:"size"`
		Units string `json:"units"`
	} `json:"partitions"`
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("I expected only one argument - path to directory containing flake.nix, you given me", len(os.Args))
	}
	eval_cmd := exec.Command("nix", "eval", "--json", os.Args[1]+"#usbix")
	out, err := eval_cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Error while evaluating flake,", "output:", string(out))
	}
	var decoded Input
	fmt.Println("unmarshling", strings.TrimSpace(strings.Split(string(out), "\n")[len(strings.Split(string(out), "\n"))-1]))
	err = json.Unmarshal([]byte(strings.TrimSpace(strings.Split(string(out), "\n")[len(strings.Split(string(out), "\n"))-1])), &decoded)
	if err != nil {
		log.Fatalln("Error while running unmarshal,", "err:", err)
	}
	for _, device := range decoded {
		dev := device.Device
		typ := device.Type
		if !exists(dev) {
			fmt.Println("File", dev, "does not exist, skipping...")
			continue
		}
		fmt.Println("Erasing device", dev, "[y/N]")
		var ans string
		fmt.Scan(&ans)
		if ans != "y" {
			fmt.Println("You said no, skipping...")
			continue
		}
		if typ != "fat32-partitions" {
			fmt.Println("Type", typ, "is not supported, skipping...")
			continue
		}

		exec.Command("parted", dev, "-s", "--", "mklabel", "gpt").Run()

		prev := 0
		num := 1
		for _, part := range device.Partitions {
			lab := part.Label
			siz := part.Size
			uni := part.Units
			output, err := exec.Command("parted", dev, "-s", "--", "mkpart", lab, "fat32", strconv.Itoa(prev)+uni, strconv.Itoa(prev+siz)+uni).CombinedOutput()
			if err != nil {
				log.Fatalln("Error while running mkpart,", "output:", string(output))
			}
			prev += siz
			output, err = exec.Command("mkfs.fat", "-F32", "-n", lab, dev+strconv.Itoa(num)).CombinedOutput()
			if err != nil {
				log.Fatalln("Error while running mkfs,", "output:", string(output))
			}
			num++
		}

		fmt.Println("Successfully partitioned", dev)
	}
}
