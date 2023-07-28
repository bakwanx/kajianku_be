package main

import "kajianku_be/routes"

func main() {
	e := routes.New()
	e.Logger.Fatal(e.Start(":8080"))
}
