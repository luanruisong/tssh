package cmd

import (
	"fmt"

	"github.com/luanruisong/tssh/constant"
)

func Logo() {
	fmt.Print(constant.LogoStr)
}

func Help() {
	fmt.Print(constant.HelpStr)
}
