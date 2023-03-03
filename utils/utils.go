package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func TildeToAbsolutePath(path string) string {
	var s string
	s = path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error could not get home directory")
	}

	if strings.Contains(path, "~") {
		s = strings.Replace(path, "~", homeDir, 1)
	}

	return s
}

func FindPath(program string) (string, error) {
	path, err := exec.LookPath(program)
	if err != nil {
		return "", errors.New("could not find " + program)
	}
	return path, nil
}

func ExecCommand(path string, command []string, stdout, stderr bool) ([]string, error) {
	var subs []string
	var out bytes.Buffer

	cmd := exec.Command(path, command...)
	if stdout && stderr {
		cmd.Stderr = os.Stderr
		cmd.Stdout = io.MultiWriter(os.Stdout, &out)
	} else if stdout && !stderr {
		cmd.Stdout = io.MultiWriter(os.Stdout, &out)
	} else if !stdout && stderr {
		cmd.Stderr = os.Stderr
		cmd.Stdout = &out
	} else {
		cmd.Stdout = &out
	}

	err := cmd.Run()
	if err != nil {
		return nil, errors.New("error running command " + path)
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		subs = append(subs, scanner.Text())
	}
	return subs, nil
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

type lines []string

func ReadLines(path string) lines {
	var lines []string

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open %s, err: %v", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines
}

func (l lines) Print() {
	for _, v := range l {
		fmt.Printf("%s\n", v)
	}
}
