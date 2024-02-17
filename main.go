package main

import (
	"errors"
	"strconv"

	"github.com/Diegiwg/cli"
)

var G_SAVE_FILES = false

func RunCMD(ctx *cli.Context, callback func(int) error) error {
	if len(ctx.Args) < 1 {
		return errors.New("the distance to be covered is not provided")
	}

	distance, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return errors.New("the distance to be covered is not a integer number")
	}

	if _, ok := ctx.Flags["save"]; ok {
		G_SAVE_FILES = true
	}

	callback(distance)

	return nil
}

func LocalCMD(ctx *cli.Context) error {
	return RunCMD(ctx, func(i int) error {
		// TODO: IMPL the local data handler
		return nil
	})
}

func main() {
	app := cli.NewApp()

	app.AddCommand(&cli.Command{
		Name:  "local",
		Desc:  "Test the application using the data saved in the API folder",
		Help:  "Test the application using the data saved in the API folder",
		Usage: "<distance: int as MGLT> [--save: optional]",
		Exec:  LocalCMD,
	})

	err := app.Run()
	if err != nil {
		println(err.Error())
	}
}
