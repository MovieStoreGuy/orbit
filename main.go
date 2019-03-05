package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	"github.com/MovieStoreGuy/termi"

	"github.com/MovieStoreGuy/orbit/instances/gcp"
)

var (
	helpRequired = false
	fs           = termi.NewFlagSet()
	ctx          = context.Background()
)

var (
	settings = &Config{}
)

func init() {
	fs.SetEnvironment(settings).
		SetDescription(UsageTemplate)
	fs.Register(termi.Must(termi.NewFlag(&helpRequired)).
		SetDescription("displays this help message").
		SetName("help").
		SetName("h"))
	registerFlags(fs, settings)
}

func main() {
	if err := fs.ParseEnvironment(); err != nil {
		panic(err)
	}
	if _, err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	if helpRequired {
		if err := fs.PrintDescription(os.Stderr); err != nil {
			panic(err)
		}
		return
	}
	if err := validator.New().Struct(settings); err != nil {
		panic(err)
	}
	client, err := gcp.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	instances, err := client.GetInstances(ctx, settings.Project, strings.Split(settings.Filters, ",")...)
	if err != nil {
		panic(err)
	}
	for _, instance := range instances {
		fmt.Println("Instance", instance)
	}
}
