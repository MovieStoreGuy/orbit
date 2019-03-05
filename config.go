package main

import "github.com/MovieStoreGuy/termi"

type Config struct {
	Project string `env:"ORBIT_GOOGLE_PROJECT" validate:"required"`
	Filters string `env:"ORBIT_INSTANCE_FILTERS" description:"Allows to filter instances based on GCP properties"`
}

func registerFlags(fs termi.FlagSet, conf *Config) {
	fs.Register(termi.Must(termi.NewFlag(&conf.Project)).
		SetDescription("defines the gcp project to examine").
		SetName("project").
		SetName("p"))
	fs.Register(termi.Must(termi.NewFlag(&conf.Filters)).
		SetDescription("filter instances based off gcp properties").
		SetName("filter"))
}
