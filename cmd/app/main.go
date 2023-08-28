package main

import "github.com/ABDURAZZAKK/avito_experiment/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
