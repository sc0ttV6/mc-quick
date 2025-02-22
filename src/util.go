package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/computerdane/gears"
)

func printStep(step string) {
	fmt.Printf("\n\033[0;32m=> %s...\033[0m\n", step)
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func checkSha1(file string, sha1sum string) bool {
	if !fileExists(file) {
		return false
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return false
	}

	hash := sha1.New()
	hash.Write(data)
	sum := hash.Sum(nil)
	actual := hex.EncodeToString(sum)
	if actual == sha1sum {
		fmt.Printf("OK %s\n", file)
		return true
	} else {
		fmt.Printf("Warning: Checksum did not match for %s. Expected %s but got %s\n", file, sha1sum, actual)
		return false
	}
}

func download(file string, loc string, sha1sum string) {
	if gears.BoolValue("overwrite") || sha1sum == "" || !checkSha1(file, sha1sum) {
		fmt.Printf("Downloading %s...\n", file)

		res, err := http.Get(loc)
		if err != nil {
			log.Fatalf("Failed to download %s: %s", loc, err)
		}
		defer res.Body.Close()

		dir := path.Dir(file)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to make directory %s: %s", dir, err)
		}

		out, err := os.Create(file)
		if err != nil {
			log.Fatalf("Failed to create file %s: %s", file, err)
		}

		if _, err := io.Copy(out, res.Body); err != nil {
			log.Fatalf("Failed to write file %s: %s", file, err)
		}
	}
}

func run(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("%s command failed: %s", name, err)
	}
}
