package server

import "github.com/metafates/avito-task/db"

// Options for the server
type Options struct {
	// Connections to the databases
	Connections db.Connections
}
