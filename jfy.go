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

type Dispatcher func(stdout, stderr []byte, args ...string) (any, any, error)

var Dispatchers = map[string]Dispatcher{
	"true": ls.Handle,
}

func main() {
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
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ProcessState.ExitCode())
		} else {
			panic(err)
		}
	}
	stdout := stdoutBuf.Bytes()
	stderr := stderrBuf.Bytes()
	handler := Dispatchers[os.Args[1]]
	if handler == nil {
		panic("No handler found")
	}
	if data, dataErr, err := handler(stdout, stderr, os.Args[2:]...); err != nil {
		panic(err)
	} else {
		if serializedJSON, err := json.Marshal(data); err != nil {
			panic(err)
		} else {
			fmt.Printf("%s\n", serializedJSON)
		}
		if serializedJSON, err := json.Marshal(dataErr); err != nil {
			panic(err)
		} else {
			logger.Printf("%s", serializedJSON)
		}
	}
}
