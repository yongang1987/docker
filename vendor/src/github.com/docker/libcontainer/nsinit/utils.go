package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/docker/libcontainer"
	"github.com/docker/libcontainer/configs"
)

func loadConfig(context *cli.Context) (*configs.Config, error) {
	if context.Bool("create") {
		config := getTemplate()
		modify(config, context)
		return config, nil
	}
	f, err := os.Open(context.String("config"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config *configs.Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadFactory(context *cli.Context) (libcontainer.Factory, error) {
	return libcontainer.New(context.GlobalString("root"), libcontainer.Cgroupfs)
}

func getContainer(context *cli.Context) (libcontainer.Container, error) {
	factory, err := loadFactory(context)
	if err != nil {
		return nil, err
	}
	container, err := factory.Load(context.String("id"))
	if err != nil {
		return nil, err
	}
	return container, nil
}

func fatal(err error) {
	if lerr, ok := err.(libcontainer.Error); ok {
		lerr.Detail(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func fatalf(t string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, t, v...)
	os.Exit(1)
}
