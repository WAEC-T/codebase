// internal/config/config.go
package config

import (
	"text/template"

	"gorm.io/gorm"
)

var (
	// Exported global variables
	DB  *gorm.DB
	Tpl *template.Template
)
