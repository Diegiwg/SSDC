package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/Diegiwg/cli"
)

var G_SAVE_FILES = false

type StarshipInfo struct {
	UID string `json:"uid"`
}

type GeneralResponse struct {
	TotalRecords int            `json:"total_records"`
	Results      []StarshipInfo `json:"results"`
}

func local_general_data() *GeneralResponse {
	file, err := os.Open("api/general.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := &GeneralResponse{}
	json.NewDecoder(file).Decode(data)

	return data
}

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
		data := local_general_data()

		println("Total records: " + strconv.Itoa(data.TotalRecords))

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
