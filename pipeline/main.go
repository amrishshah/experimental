package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func runPipeline() {

	//cmd1Name, args1 := splitCommandAndArguments(left)
	//cmd2Name, args2 := splitCommandAndArguments(right)
	cmd1Name := "tail"
	cmd2Name := "wc"
	cmd1 := exec.Command(cmd1Name, "-f", "/tmp/foo/bar.md")
	cmd2 := exec.Command(cmd2Name)

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd1.Stderr = os.Stderr
	cmd2.Stdin = r
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	if err := cmd1.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting %s: %v\n", cmd1Name, err)
		w.Close()
		r.Close()
		return
	}
	if err := cmd2.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting %s: %v\n", cmd2Name, err)
		w.Close()
		r.Close()
		return
	}

	go func() {
		cmd1.Wait()
		w.Close()
	}()

	cmd2.Wait()
}

func main() {
	runPipeline()
}
