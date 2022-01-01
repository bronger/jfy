package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/bronger/jfy/lib"
	"github.com/bronger/jfy/uptime"
)

var logger = log.New(os.Stderr, "", 0)

var dispatchers = map[string]lib.Dispatcher{
	"uptime": uptime.Handle,
}

var settings lib.SettingsType

func init() {
	settings = lib.SettingsType{ExitCode: 221, Version: 99999999}
	data := os.Getenv("JFY_SETTINGS")
	if data != "" {
		if err := json.Unmarshal([]byte(data), &settings); err != nil {
			logger.Println(err)
			os.Exit(221)
		}
		if settings.ExitCode < 1 || settings.ExitCode > 255 {
			logger.Println("Invalid exit code in settings")
			os.Exit(221)
		}
		if settings.Version < 0 || settings.Version > 99999999 {
			logger.Println("Invalid version number in settings")
			os.Exit(settings.ExitCode)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		logger.Println("Too few arguments")
		os.Exit(settings.ExitCode)
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	go func() {
		signals := make(chan os.Signal, 32)
		signal.Notify(signals)
		for cmd.ProcessState == nil {
			sig := <-signals
			for cmd.Process == nil {
				time.Sleep(10 * time.Millisecond)
			}
			cmd.Process.Signal(sig)
		}
		signal.Reset()
	}()
	cmd.Stdin = os.Stdin
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if err := cmd.Start(); err != nil {
		logger.Println(err)
		os.Exit(settings.ExitCode)
	}
	var exitCode int
	if err := cmd.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ProcessState.ExitCode()
		} else {
			logger.Println(err)
			os.Exit(settings.ExitCode)
		}
	}
	stdout := stdoutBuf.Bytes()
	stderr := stderrBuf.Bytes()
	handler := dispatchers[os.Args[1]]
	if handler == nil {
		panic("No handler found")
	}
	if data, dataErr, err := handler(settings, stdout, stderr, os.Args[1:]...); err != nil {
		logger.Println(err)
		os.Exit(settings.ExitCode)
	} else {
		if data != nil {
			if serializedJSON, err := json.Marshal(data); err != nil {
				logger.Println(err)
				os.Exit(settings.ExitCode)
			} else {
				fmt.Printf("%s\n", serializedJSON)
			}
		}
		if dataErr != nil {
			if serializedJSON, err := json.Marshal(dataErr); err != nil {
				logger.Println(err)
				os.Exit(settings.ExitCode)
			} else {
				logger.Printf("%s", serializedJSON)
			}
		}
	}
	os.Exit(exitCode)
}
