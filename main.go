package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Diegiwg/cli"
)

var HTTPClient = &http.Client{Timeout: 10 * time.Second}
var G_SAVE_FILES = false

type Starship struct {
	Name        string `json:"name"`
	MGLT        string `json:"MGLT"`
	Consumables string `json:"consumables"`
}

func (s *Starship) Consume() int {
	split := strings.Split(s.Consumables, " ")

	value, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}

	units := strings.TrimSuffix(split[1], "s")

	switch units {
	case "year":
		return 8760 * value

	case "month":
		return 730 * value

	case "week":
		return 168 * value

	case "day":
		return 24 * value

	default:
		return 1
	}
}

func (s *Starship) Stops(distance_MGLT int) int {
	// f:
	// distance / (consumables%toHours * MGLT)

	local_MGLT, err := strconv.Atoi(s.MGLT)
	if err != nil {
		panic(err)
	}

	return distance_MGLT / (s.Consume() * local_MGLT)
}

type StarshipData struct {
	Properties Starship `json:"properties"`
}

type StarshipResponse struct {
	Result StarshipData `json:"result"`
}

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

func local_starship_data(uid string) *StarshipResponse {
	file, err := os.Open("api/starships/" + uid + ".json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := &StarshipResponse{}
	json.NewDecoder(file).Decode(data)

	return data
}

func local_process_data(data *GeneralResponse) []Starship {
	result := []Starship{}

	for index := range data.Results {
		starship := local_starship_data(data.Results[index].UID)
		result = append(result, starship.Result.Properties)
	}

	return result
}

func remote_general_data(page int, limit int) *GeneralResponse {
	URL := fmt.Sprintf("https://www.swapi.tech/api/starships/?page=%d&limit=%d", page, limit)
	req, err := HTTPClient.Get(URL)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		panic(fmt.Sprintf("status code error: %d %s", req.StatusCode, req.Status))
	}

	data := &GeneralResponse{}
	json.NewDecoder(req.Body).Decode(data)

	if G_SAVE_FILES {
		os.Mkdir("api", 0755)

		file, err := os.Create("api/general.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		asJson, _ := json.MarshalIndent(data, "", "  ")
		file.Write(asJson)
	}

	return data
}

func remote_starship_data(uid string) *StarshipResponse {
	URL := fmt.Sprintf("https://www.swapi.tech/api/starships/%s", uid)

	req, err := HTTPClient.Get(URL)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		panic(fmt.Sprintf("status code error: %d %s", req.StatusCode, req.Status))
	}

	data := &StarshipResponse{}
	json.NewDecoder(req.Body).Decode(data)

	if G_SAVE_FILES {
		os.Mkdir("api/starships", 0755)

		file, err := os.Create("api/starships/" + uid + ".json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		asJson, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
		file.Write(asJson)
	}

	return data
}

func remote_process_data(data *GeneralResponse) []Starship {
	result := []Starship{}

	for index := range data.Results {
		starship := remote_starship_data(data.Results[index].UID)
		result = append(result, starship.Result.Properties)
	}

	return result
}

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

func LocalCMD(ctx *cli.Context) error {
	return RunCMD(ctx, func(distance int) {
		G_SAVE_FILES = false

		data := local_general_data()
		starships := local_process_data(data)

		for index := range starships {
			starship := starships[index]

			if starship.MGLT == "unknown" {
				continue
			}

			println(starship.Name+":", starship.Stops(distance))
		}
	})
}

func RemoteCMD(ctx *cli.Context) error {
	return RunCMD(ctx, func(distance int) {
		data := remote_general_data(1, 100)
		starships := remote_process_data(data)
		println("Total:", len(starships))
	})
}

func main() {
	app := cli.NewApp()

	app.AddCommand(&cli.Command{
		Name:  "local",
		Desc:  "Test the application using the data saved in the API folder",
		Help:  "Test the application using the data saved in the API folder",
		Usage: "<distance: int as MGLT>",
		Exec:  LocalCMD,
	})

	app.AddCommand(&cli.Command{
		Name:  "remote",
		Desc:  "Test the application using the live data",
		Help:  "Test the application using the live data",
		Usage: "<distance: int as MGLT> [--save]",
		Exec:  RemoteCMD,
	})

	err := app.Run()
	if err != nil {
		println(err.Error())
	}
}
