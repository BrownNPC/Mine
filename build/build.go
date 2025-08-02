package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var usage string = `usage:
	go run ./build/build.go windows
`

func main() {
	flag.Parse()
	GOOS := flag.Arg(0)
	switch GOOS {
	case "windows":
		BuildWindows()
	default:
		fmt.Println(usage)
	}
}
func BuildWindows() {
	os.Setenv("CGO_ENABLED", "1")
	os.Setenv("CC", "x86_64-w64-mingw32-gcc")
	os.Setenv("GOOS", "windows")

	cmd := exec.Command("rm", "./Mine.exe")
	cmd.Run()

	cmd = exec.Command("go", "build", "-v", "-o", "Mine.exe", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("./Mine.exe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
