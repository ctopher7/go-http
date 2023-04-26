package main

import (
	"fmt"
	"net/http"
	"os"

	"example.com/m/v2/config"
	db "example.com/m/v2/database"
	"example.com/m/v2/dependency"
	"example.com/m/v2/resource"
	"example.com/m/v2/route"
)

func main() {
	//init config
	cfg, err := config.ReadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	//init resources
	res := resource.Init(&cfg)
	defer res.PostgresDb.Close()

	switch os.Args[2] {
	case "server":
		dep := dependency.Init(&cfg, res)

		route.Init(dep)

		fmt.Printf("running server on %s \n", cfg.ServerAddress)
		http.ListenAndServe(cfg.ServerAddress, nil)
	case "migrate":
		db.Migrate(res)
	case "seed":
		db.Seed(res)
	}
}
