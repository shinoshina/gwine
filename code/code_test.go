package code

import (
	"fmt"
	"testing"
)

func TestByteEncode(t *testing.T) {
	fmt.Printf("%b\n",Make(OpConstant,1000))
	fmt.Printf("%s\n",Instructions(Make(OpConstant,1000)).String())
}