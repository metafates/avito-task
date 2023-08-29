package server

import "github.com/metafates/avito-task/db"

// Options for the server
type Options struct {
	// Pools of databases
	Pools db.Pools
}
