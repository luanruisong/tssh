// +build darwin freebsd linux netbsd openbsd solaris

package constant

import "github.com/manifoldco/promptui"

const (
	HOME = "HOME"
)

var (
	ListTpl = &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336  {{ .FmtName | cyan }} ({{ .User | yellow }}@{{ .Ip | red }})",
		Inactive: "   {{ .FmtName | cyan }} ({{ .User | yellow }}@{{ .Ip | red }})",
		Selected: "start connect {{ .Name | cyan }}({{ .User | yellow }}@{{ .Ip | red }})...",
		Details:  Detail,
	}
)
