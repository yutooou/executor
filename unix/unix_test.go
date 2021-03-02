// created by yutooou
package unix

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"
)

func TestLookPath(t *testing.T) {
	path, err := exec.LookPath("gcc")
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
}

func TestUnixShell(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	result, err := UnixShell(&ShellOptions{
		Name:    "go",
		Args:    []string{"version"},
		Context: ctx,
	})
	if err != nil {
		log.Panic(err)
	} else {
		fmt.Println(result.Stdout)
	}
}
