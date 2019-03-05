package main

const UsageTemplate = `{{.name}} -- usage
Go Runtime version ({{.GoVersion}})

{{.name}} is an application that allows you to view your VMs within your cloud provider.
This allows you to generate an inventory that can be used to with server validation tools
to ensure everything is working correctly.

Environments :
{{ range $name, $description := .environments }}{{ $name }}	{{ $description }}
{{ end }}
Flags:
{{ range $flag := .flags }}{{ $flag }}
{{ end }}
`
