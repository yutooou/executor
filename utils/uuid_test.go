// created by yutooou
package utils

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	str1 := UUID(12)
	str2 := UUID(12)
	fmt.Println(str1, str2)
}
