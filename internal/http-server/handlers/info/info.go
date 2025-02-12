package info

import "github.com/justcgh9/merch_store/internal/models/inventory"

type Informator interface {
	Informate(username string) inventory.Info
}

type InfoRequest struct {
}
