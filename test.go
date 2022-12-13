package main

import (
	app "androidtoolbox/backend"
	"fmt"
)

func main1() {
	app := app.NewApp()
	a := app.GetCurrentActivity()
	fmt.Printf("%+v\n", a)
}
