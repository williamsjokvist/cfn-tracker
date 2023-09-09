//go:build darwin

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
)

func cleanUpProcess() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()

	fmt.Println("Current Process ID: ", os.Getpid())

	processes, _ := process.Processes()
	for _, process := range processes {
		name, _ := process.Name()

		if !strings.Contains(name, `CFN Tracker`) || process.Pid == int32(os.Getpid()) {
			continue
		}

		fmt.Println("Killing remnant process by ID: ", process.Pid, " Name: ", name)
		err := process.Kill()
		if err != nil {
			err = process.Terminate()
			if err != nil {
				fmt.Println("Error killing remnant process by ID: ", process.Pid, " msg: ", err)
			}
			continue
		}

		time.Sleep(1 * time.Second)
	}
}
