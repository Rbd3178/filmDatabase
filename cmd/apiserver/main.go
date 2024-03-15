package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Rbd3178/filmDatabase/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.json", "path to configuration file")
}

func main() {
	flag.Parse()
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := apiserver.NewConfig()
	if err := json.Unmarshal(jsonData, &config); err != nil {
		log.Fatal(err)
	}
	
	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
