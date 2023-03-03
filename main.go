package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dvvedz/gore/utils"
)

func main() {
	fileFlag := flag.String("f", "", "Takes a list of sub-domains or urls")
	wildFlag := flag.Bool("wild", false, "Add wildcards found to stdout (last)")
	wildOnlyFlag := flag.Bool("wildonly", false, "only print wildcards to stdout")
	resolversFlag := flag.String("r", "~/Hacking/wordlists/resolvers.txt", "use custom resolvers")
	cleanFlag := flag.Bool("clean", false, "clean up all files generated from script, after script run")

	flag.Parse()

	// TODO read input from stdin

	// Required flags
	if *fileFlag == "" {
		fmt.Println("-file flag is required")
		fmt.Println("")
		flag.Usage()
		os.Exit(1)
	}
	if *wildFlag && *wildOnlyFlag {
		fmt.Println("can't use both -wild and -wildonly flags at the same time")
		fmt.Println("")
		flag.Usage()
		os.Exit(1)
	}

	// Test inputs, Should later be supplied as flags adn arguments
	rfile := *resolversFlag
	wfile := *fileFlag

	timestamp := time.Now().Unix()
	outWildcards := fmt.Sprintf("/tmp/wildcards-%d.txt", timestamp)

	rfile = utils.TildeToAbsolutePath(rfile)
	wfile = utils.TildeToAbsolutePath(wfile)

	// check if provided file exits
	if !utils.FileExists(rfile) {
		fmt.Printf("error: file: %s does not exist\n", rfile)
		os.Exit(1)
	}

	if !utils.FileExists(wfile) {
		fmt.Printf("error: file: %s does not exist\n", wfile)
		os.Exit(1)
	}

	p, perr := utils.FindPath("puredns")
	if perr != nil {
		log.Fatalf("puredns not found, err: %v", perr)
	}

	cmd := []string{"resolve", wfile, "-r", rfile, "--write-wildcards", outWildcards}

	var printStdout bool

	if !*wildFlag && *wildOnlyFlag {
		printStdout = false
	} else {
		printStdout = true
	}

	_, cerr := utils.ExecCommand(p, cmd, printStdout, false)
	if cerr != nil {
		log.Fatalf("%v", cerr)
	}

	if *wildFlag && !*wildOnlyFlag {
		printWildcards(outWildcards, "*")
	} else if !*wildFlag && *wildOnlyFlag {
		printWildcards(outWildcards, "")
	}

	// cleanup
	if *cleanFlag {
		fileCleanup()
	}
}

func printWildcards(path, prefix string) {
	for _, v := range utils.ReadLines(path) {
		if len(prefix) > 0 {
			fmt.Printf("%s.%s\n", prefix, v)
		} else {
			fmt.Printf("%s\n", v)
		}
	}
}

func fileCleanup() {
	files, err := filepath.Glob("/tmp/wildcards-*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
