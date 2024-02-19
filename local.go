package main

import (
	"encoding/json"
	"os"

	cli "github.com/Diegiwg/cli"
)

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
