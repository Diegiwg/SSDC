package main

import (
	"encoding/json"
	"fmt"
	"os"

	cli "github.com/Diegiwg/cli"
)

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

func RemoteCMD(ctx *cli.Context) error {
	return RunCMD(ctx, func(distance int) {
		data := remote_general_data(1, 100)
		starships := remote_process_data(data)

		for index := range starships {
			starship := starships[index]

			if starship.MGLT == "unknown" {
				continue
			}

			println(starship.Name+":", starship.Stops(distance))
		}
	})
}
