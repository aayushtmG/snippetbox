package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/aayushtmG/snippetbox/internal/models"
	"github.com/aayushtmG/snippetbox/ui"
)


type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string
}



func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}


var functions = template.FuncMap{
	"humanDate": humanDate,
}


func newTemplateCache() (map[string]*template.Template,error){

	//just a regular map acting as cache
	cache := map[string]*template.Template{}
	
	pages,err := fs.Glob(ui.Files,"html/pages/*.tmpl")	
	if err != nil {
		return nil,err
	}


	for _,page := range pages {

		//extracts the file name like (home.tmpl)
		name := filepath.Base(page)	


		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files,patterns...)
		if err != nil {
			return nil, err
		}

		// ts,err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		// if err!= nil {
		// 	return nil, err
		// }
		//
		//
		// ts, err = ts.ParseFiles(page)
		// if err!= nil {
		// 	return nil, err
		// }

		cache[name] = ts
	}




	return cache,nil
}
