package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func run(dir, bin string, args ...string) {

	path, err := exec.LookPath(bin);
	if err != nil {
		log.Fatalf("cant locate %s. Please install!", bin)
	}

	cmd := exec.Command(path, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("→ (%s) %s %v\n", dir, path, args)
	if err := cmd.Run(); err != nil {
		log.Fatalf("%s failed: %v", path, err)
	}
}

func main() {

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	webDir := filepath.Join(root, "src/web")

	run(webDir, "npm", "install")
	run(webDir, "npm", "run", "build")

	run(root, "go", "build", "-o", "bin/scanbridge", "./src")
}
