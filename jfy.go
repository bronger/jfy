package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

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
	//	stderr := stderrBuf.Bytes()
	fmt.Printf("%s", stdout)
}
