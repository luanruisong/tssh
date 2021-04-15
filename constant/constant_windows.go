// +build windows

package constant

import "github.com/manifoldco/promptui"

const (
	HOME = "HOMEPATH"
)

var (
	ListTpl = &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "->{{ .FmtName | cyan }} ({{ .User | yellow }}@{{ .Ip | red }})",
		Inactive: "  {{ .FmtName | cyan }} ({{ .User | yellow }}@{{ .Ip | red }})",
		Selected: "start connect {{ .Name | cyan }}({{ .User | yellow }}@{{ .Ip | red }})...",
		Details:  Detail,
	}
)
