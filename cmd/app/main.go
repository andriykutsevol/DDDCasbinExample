package main

import (
	"github.com/andriykutsevol/DDDCasbinExample/internal/app"
)

const configsDir = "../../configs/"

//go:generate go env -w GO111MODULE=on
//go:generate go mod tidy
//go:generate go mod download

func main() {
	app.Run(configsDir)
}
