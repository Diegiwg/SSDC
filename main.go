package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	cli "github.com/Diegiwg/cli"
)

var HTTPClient = &http.Client{Timeout: 10 * time.Second}
var G_SAVE_FILES = false

func RunCMD(ctx *cli.Context, callback func(distance int)) error {
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

func main() {
	app := cli.NewApp()

	app.AddCommand(&cli.Command{
		Name:  "web",
		Desc:  "Test the application on a web page",
		Help:  "Test the application on a web page",
		Usage: "[--local]",
		Exec:  WebCMD,
	})

	app.AddCommand(&cli.Command{
		Name:  "remote",
		Desc:  "Test the application using the live data",
		Help:  "Test the application using the live data",
		Usage: "<distance> [--save]",
		Exec:  RemoteCMD,
	})

	app.AddCommand(&cli.Command{
		Name:  "local",
		Desc:  "Test the application using the data saved in the API folder",
		Help:  "Test the application using the data saved in the API folder",
		Usage: "<distance>",
		Exec:  LocalCMD,
	})

	err := app.Run()
	if err != nil {
		println(err.Error())
	}
}
