package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"log"
)

const version = "0.1"

var flagVersion = false
var flagHelp = false

var printVersion = flag.CommandLine.Name() + " version: " + version
var printUsage = "usage: " + flag.CommandLine.Name() + " keyword ..."

func manpath() [][]byte {
	// Should return a list of elements
	out, err := exec.Command("man", "--path").Output()
	if err != nil {
		log.Fatalf("Oops man command failed: %v", err)
	}
	return bytes.Split(out, []byte{0x3a})
}

func getOccurences(search string) {
	mp := manpath()
	var wg sync.WaitGroup
	wg.Add(len(mp))
	go func() {
		for d := range mp {
			if _, err := os.Stat(string(mp[d])); err == nil {
				stat, err := os.Lstat(string(mp[d]))
				if err != nil {
					log.Fatal(err)
				}
				// Ensure if the statement check if first condition is a file o else it's dir
				if _, err := os.Stat(string(mp[d]) + "/whatis"); err == nil {
					// Write the correct implementation
					c, err := exec.Command("grep", "-i", search, string(mp[d])+"/whatis").CombinedOutput()
					if err != nil {
						// fmt.Println("ERRRRRRRR0: ", err)
					} else {
						fmt.Printf("%s", c)
					}
				} else if stat.IsDir() || stat.Mode()&os.ModeSymlink != 0 {
					c, err := exec.Command("/usr/libexec/makewhatis.local", "-o /dev/fd/1", string(mp[d])).Output()
					if err != nil {
						// fmt.Println("ERRRRRRRR1: ", err)
					} else if strings.ContainsAny(string(c), search) {
						re := regexp.MustCompile("(?i).*" + search + ".*")
						matches := re.FindAllStringSubmatch(string(c), 100)
						for m := range matches {
							fmt.Println(strings.Join(matches[m], ""))
						}
					}
				} else {
					fmt.Printf("%s: nothing appropriate\n", search)
				}
			}
			defer wg.Done()
		}
	}()
	wg.Wait()
}
func init() {
	flag.BoolVar(&flagVersion, "version", false, "display program version")
	flag.BoolVar(&flagHelp, "help", false, printUsage)
}

func main() {

	flag.Parse()
	switch {
	case flagVersion:
		fmt.Println(printVersion)
	case flagHelp:
		flag.Usage()
	case len(os.Args) <= 1:
		flag.Usage()
		os.Exit(1)
	case len(os.Args) > 1:
		getOccurences(os.Args[1])
	}

}
