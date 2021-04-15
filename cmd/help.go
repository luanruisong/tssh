package cmd

import (
	"fmt"

	"tssh/constant"
)

func Logo() {
	fmt.Print(constant.LogoStr)
}

func Help() {
	fmt.Print(constant.HelpStr)
}
