package main

import (
	"fmt"

	"github.com/CustomCloudStorage/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	fmt.Println(cfg)
}
