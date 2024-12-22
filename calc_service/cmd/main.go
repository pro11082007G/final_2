package main

import (
	"github.com/pro11082007G/calc_service/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
