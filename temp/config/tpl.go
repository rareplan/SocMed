package config

import (
	"html/template"
	"log"
)

var (
	TPL *template.Template
	fm  = template.FuncMap{}
)

func init() {

	var err error
	TPL, err = template.New("").Funcs(fm).ParseGlob("template/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

}
