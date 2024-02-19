package main

import (
	"os"
	"strconv"
	"strings"

	cli "github.com/Diegiwg/cli"
	web "github.com/Diegiwg/dwork-web/dw"
	log "github.com/Diegiwg/dwork-web/dw/logger"
)

func load_file(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}

var CACHED_STARSHIPS []Starship
var BASE_CSS = load_file("web/file.css")
var BASE_JS = load_file("web/file.js")
var BASE_HTML = load_file("web/file.html")

func render(template string, data map[string]string) string {
	data["CSS"] = "<style>" + BASE_CSS + "</style>"
	data["JS"] = "<script>" + BASE_JS + "</script>"

	for key, value := range data {
		template = strings.ReplaceAll(template, "{{"+key+"}}", value)
	}

	return template
}

func render_results(starships []Starship, distance int) string {
	if len(starships) == 0 {
		return `<div id="results"></div>`
	}

	result := `<ul id="results">`
	for index := range starships {
		starship := starships[index]

		if starship.MGLT == "unknown" {
			continue
		}

		result += "<li>" + starship.Name + ": " + strconv.Itoa(starship.Stops(distance)) + "</li>"
	}
	result += "</ul>"

	return result
}

func render_home() string {
	data := map[string]string{
		"Distance": "",
		"Results":  render_results([]Starship{}, 0),
	}

	return render(BASE_HTML, data)
}

func render_calculation_results(distance int) string {
	data := map[string]string{
		"Distance": strconv.Itoa(distance),
		"Results":  render_results(CACHED_STARSHIPS, distance),
	}

	return render(BASE_HTML, data)
}

func render_error(msg string) string {
	return render(BASE_HTML, map[string]string{
		"Distance": "",
		"Results":  "<h2 class=\"error\">" + msg + "</h2>",
	})
}

func WebCMD(cliCtx *cli.Context) error {
	if _, ok := cliCtx.Flags["local"]; ok {
		log.Info("Using local starship database...")

		starships_data := local_general_data()
		CACHED_STARSHIPS = local_process_data(starships_data)

	} else {
		log.Info("Checking and caching the starship database...")

		starships_data := remote_general_data(1, 100)
		CACHED_STARSHIPS = remote_process_data(starships_data)

		if len(CACHED_STARSHIPS) == 0 {
			log.Fatal("Starship database is not available")
		}
	}

	webApp := web.MakeApp()

	webApp.GET("/", func(webCtx web.Context) {
		webCtx.Response.Html(render_home())
	})

	webApp.GET("/<int:distance>", func(webCtx web.Context) {
		distance, err := webCtx.Request.Params.Int("distance")
		if err != nil || distance < 1 {
			webCtx.Response.Html(render_error("Invalid distance"))
			return
		}

		webCtx.Response.Html(render_calculation_results(distance))
	})

	log.Success("Starting web server...")

	webApp.Serve(":8081")

	return nil
}
