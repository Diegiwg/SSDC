package main

import (
	"strconv"
	"strings"
)

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
