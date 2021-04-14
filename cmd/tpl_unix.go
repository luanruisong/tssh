// +build darwin freebsd linux netbsd openbsd solaris

package cmd

import "github.com/manifoldco/promptui"

var listTpl = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F336  {{ .FmtName | cyan }} ({{ .User | yellow }}@{{ .Ip | red }})",
	Inactive: "   {{ .FmtName | cyan }} ({{ .User | yellow }}@{{ .Ip | red }})",
	Selected: "start connect {{ .Name | cyan }}({{ .User | yellow }}@{{ .Ip | red }})...",
	Details: `
{{ "Name:" | faint }}	{{ .Name }}
{{ "Ip:" | faint }}	{{ .Ip }}
{{ "User:" | faint }}	{{ .User }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "ConnMode:" | faint }}	{{ .ConnMode }}
{{ "SaveAt:" | faint }}	{{ .SaveAt }}`,
}
