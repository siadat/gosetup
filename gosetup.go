// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Setup is a simple ... TODO(sina)
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	stateConfirmGopath = iota
	stateEditGopath
	statePersistGoPath
	stateAddSampleProg
	stateExit
)

const defaultGOPATH = "$HOME/go"
const helloProgram = `package main

import "fmt"

func main() {
    fmt.Printf("hello, world\n")
}`

func main() {
	log.SetFlags(0)
	log.SetPrefix(">>> ")

	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		fmt.Printf("GOPATH=%q\n", gopath)
		// fmt.Fprintf(os.Stderr, "GOPATH is not set.\n")
		if err := validateGoDirs(gopath); err != nil {
			fmt.Println(err)
		}
		return
	}

	profile := shellProfile()
	if profile == "" {
		fmt.Println("shell profile not found")
		os.Exit(1)
	}

	gopath = defaultGOPATH
	state := stateConfirmGopath
	for {
		switch state {
		case stateConfirmGopath:
			s := fmt.Sprintf("export GOPATH=%s\nexport PATH=$PATH:$GOPATH/bin", gopath)

			log.Printf("Adding the following lines to %s:", profile)
			log.Printf("\n")
			fmt.Println(prefixLines(s, ">>>\t"))
			log.Printf("\n")

			switch prompt("Continue [Y,n,e,?]? ") {
			case "", "y", "Y":
				state = statePersistGoPath
			case "e":
				state = stateEditGopath
			case "n":
				state = stateExit
			default:
				log.Println()
				log.Println("Y - append to file (default)")
				log.Println("n - quit")
				log.Println("e - set a different GOPATH")
				log.Println("? - show command meanings")
				log.Println()
			}
		case statePersistGoPath:
			s := fmt.Sprintf("export GOPATH=%s\nexport PATH=$PATH:$GOPATH/bin", gopath)
			if err := appendToFile(profile, s+"\n"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := os.MkdirAll(filepath.Join(expand(gopath), "src"), 0755); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			log.Println("Done. Changes will be reflected next time a terminal is started.")
			fmt.Println("")
			state = stateAddSampleProg
		case stateAddSampleProg:
			if err := promptSampleProg(gopath); err != nil {
				fmt.Println(err)
				break
			}
			state = stateExit
		case stateEditGopath:
			var err error
			if gopath, err = promptGOPATH(); err != nil {
				fmt.Println(err)
			} else {
				state = stateConfirmGopath
			}
		case stateExit:
			os.Exit(0)
		}
	}
}

func prompt(p string) string {
	fmt.Print(p)
	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')
	if err == io.EOF {
		// CTRL-d
		fmt.Println("exit")
		os.Exit(0)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res = strings.TrimSpace(res)
	if res == "exit" || res == "quit" {
		os.Exit(0)
	}
	return res
}

func promptSampleProg(gopath string) error {
	user := "$USERNAME"
	log.Printf("Creating a hello world program with the following file:")
	log.Printf("\n")
	path := filepath.Join(expand(gopath), "src", "github.com", user, "hello", "hello.go")
	fmt.Println(prefixLines(path, ">>>\t"))
	log.Printf("\n")

	if user = prompt("Enter your GitHub username or exit: "); user == "" {
		return fmt.Errorf("error: empty user")
	}

	dir := filepath.Join(expand(gopath), "src", "github.com", user, "hello")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(dir, "hello.go"), []byte(helloProgram), 0666); err != nil {
		return err
	}

	log.Printf("Done.")
	fmt.Println()
	log.Printf("Run this program with the following commands:\n")
	log.Printf("\n")
	log.Printf("\tcd %s\n", dir)
	log.Printf("\tgo run hello.go\n")
	log.Printf("\tgo install\n")
	log.Printf("\thello\n")
	log.Printf("\n")

	return nil
}

func expand(p string) string {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch true {
	case strings.HasPrefix(p, "$HOME"):
		return strings.Replace(p, "$HOME", user.HomeDir, 1)
	case strings.HasPrefix(p, "${HOME}"):
		return strings.Replace(p, "${HOME}", user.HomeDir, 1)
	case strings.HasPrefix(p, "~"):
		return strings.Replace(p, "~", user.HomeDir, 1)
	}
	return p
}

func prefixLines(s, prefix string) string {
	return prefix + strings.Replace(s, "\n", "\n"+prefix, -1)
}

func promptGOPATH() (string, error) {
	gopath := prompt(fmt.Sprintf("Enter new GOPATH: "))

	switch true {
	case strings.HasPrefix(gopath, "$HOME"):
	case strings.HasPrefix(gopath, "${HOME}"):
	case strings.HasPrefix(gopath, "~"):
	default:
		if !filepath.IsAbs(gopath) {
			return "", fmt.Errorf("error: GOPATH must be absolute path")
		}
	}
	return gopath, nil
}

func validateGoDirs(gopath string) error {
	if err := validateDir(gopath); err != nil {
		return err
	}
	if err := validateDir(filepath.Join(gopath, "src")); err != nil {
		return err
	}
	if err := validateDir(filepath.Join(gopath, "bin")); err != nil {
		return err
	}
	if err := validateDir(filepath.Join(gopath, "pkg")); err != nil {
		return err
	}
	return nil
}

func validateDir(dir string) error {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("error: %q does not exist", dir)
	}
	if os.IsPermission(err) {
		return fmt.Errorf("error: %q permission denied", dir)
	}
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("error: %q not a directory", dir)
	}

	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("error: GOPATH must be absolute path")
		}
	}
	return nil
}

func appendToFile(file, s string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(s); err != nil {
		return err
	}
	return nil
}

func shellProfile() (p string) {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch filepath.Base(os.Getenv("SHELL")) {
	case "bash":
		for _, filename := range []string{".bashrc", ".bash_profile"} {
			p = filepath.Join(user.HomeDir, filename)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	case "zsh":
		p = filepath.Join(user.HomeDir, ".zshrc")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	p = filepath.Join(user.HomeDir, ".profile")
	if _, err := os.Stat(p); err == nil {
		return p
	}

	return ""
}
