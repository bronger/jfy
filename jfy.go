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

	"github.com/bronger/jfy/ls"
)

var logger = log.New(os.Stderr, "", 0)

type dispatcher func(stdout, stderr []byte, args ...string) (any, any, error)

var dispatchers = map[string]dispatcher{
	"true": ls.Handle,
}

type settingsType struct {
	ExitCode int
}

var settings settingsType

func init() {
	settings = settingsType{ExitCode: 221}
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
	if err := cmd.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ProcessState.ExitCode())
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
	if data, dataErr, err := handler(stdout, stderr, os.Args[2:]...); err != nil {
		logger.Println(err)
		os.Exit(settings.ExitCode)
	} else {
		if serializedJSON, err := json.Marshal(data); err != nil {
			logger.Println(err)
			os.Exit(settings.ExitCode)
		} else {
			fmt.Printf("%s\n", serializedJSON)
		}
		if serializedJSON, err := json.Marshal(dataErr); err != nil {
			logger.Println(err)
			os.Exit(settings.ExitCode)
		} else {
			logger.Printf("%s", serializedJSON)
		}
	}
}
