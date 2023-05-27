package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"bookings/pkg/config"
	"bookings/pkg/models"
)

var app *config.AppConfig

// RenderTemplate renders templates using html template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from app config
		tc = app.TemplateCache
	} else {
		// recreate template cache
		tc, _ = CreateTemplatesFromAllFiles()
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	// render the template

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}


func CreateTemplatesFromAllFiles() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	const baseTmpl = "./templates/*.layout.tmpl"

	// get all the files in template dir with name *.page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	
	// get all the files in template dir with name *.layout.tmpl
	layouts, err := filepath.Glob(baseTmpl)
	if err != nil {
		return myCache, err
	}
	
	// loop through the pages
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println("Could not parse template", page, err)
			return myCache, err
		}
		
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(baseTmpl)
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	
	return myCache, nil
}

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}
