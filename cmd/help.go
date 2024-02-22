package cmd

import (
	"fmt"
	"strings"

	"github.com/luanruisong/tssh/constant"
)

func Logo(version string) {
	if len(version) > 0 {
		s := strings.ReplaceAll(constant.LogoStr, "unknown", version)
		fmt.Print(s)
		return
	}
	fmt.Print(constant.LogoStr)
}

func Help() {
	fmt.Print(constant.HelpStr)
}
