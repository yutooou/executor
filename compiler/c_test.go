// created by yutooou
package compiler

import (
	"fmt"
	"log"
	"testing"
)

var sourceCodeC = `#include <stdio.h>

int main() {
	int a, b;
	while (~scanf("%d%d", &a, &b)) {
	    printf("%d\n", a+b);
	}
}`

func TestC_Compile(t *testing.T) {
	compiler, err := NewCompiler("c")
	if err != nil {
		log.Panic(err)
	}
	err = compiler.Init(sourceCodeC, "/tmp/yu")
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("source file create success")
	}
	err = compiler.Compile()
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("executable file create success")
	}

	runArgs := compiler.RunArgs()
	fmt.Println(runArgs)
}
