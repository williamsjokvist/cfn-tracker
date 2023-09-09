//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

func cleanUpProcess() {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return
	}

	p := windows.ProcessEntry32{Size: 568}
	didKillProcess := false
	for {
		err = windows.Process32Next(handle, &p)
		if err != nil {
			break
		}

		if !strings.Contains(windows.UTF16ToString(p.ExeFile[:]), "CFN Tracker") || p.ProcessID == uint32(os.Getpid()) {
			continue
		}

		process, err := os.FindProcess(int(p.ProcessID))
		if err != nil {
			fmt.Println("Error killing remnant process by ID: ", p.ProcessID, " msg: ", err)
			continue
		}

		err = process.Kill()
		if err != nil {
			process.Signal(syscall.SIGQUIT)
			process.Signal(syscall.SIGTERM)
			fmt.Println("Error killing remnant process by ID: ", p.ProcessID, " msg: ", err)
		}

		didKillProcess = true
	}

	if didKillProcess {
		time.Sleep(1 * time.Second)
	}
}
