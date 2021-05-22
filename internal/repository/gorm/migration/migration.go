package migration

import "github.com/go-gormigrate/gormigrate/v2"

// Migrations is a collection of storage migration patterns
var Migrations = []*gormigrate.Migration{
	v202105221650,
}
