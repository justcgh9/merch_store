package main

import "github.com/justcgh9/merch_store/internal/config"

func main() {
	cfg := config.MustLoad()

	_ = cfg
}
