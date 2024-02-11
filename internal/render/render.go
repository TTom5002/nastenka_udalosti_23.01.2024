package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/helpers"
	"nastenka_udalosti/internal/models"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	// "iterate":    Iterate,
}

var app *config.AppConfig
var pathToTemplates = "./templates"

// TODO: Pak možná odstraň
// Iterate returns a slice of ints, starting at 1, going to count
// func Iterate(count int) []int {
// 	var i int
// 	var items []int
// 	for i = 0; i < count; i++ {
// 		items = append(items, i)
// 	}
// 	return items
// }

// NewRenderer nastavá config pro template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate vrací datum v YYYY-MM-DD formatu
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate formátuje datum dle požadoveného stringu
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// AddDefaultData přidá data do templatu
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "userInfo") {
		td.IsAuthenticated = 1
		td.User, _ = helpers.GetUserInfo(r)
		// td.User = app.Session.Get(r.Context(), "userInfo").(models.User)
	}
	return td
}

// Template renderuje templaty pomocí html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache {
		// získá template cache z app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// dotaz na template z cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	//renderuje template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache vytvoří template cache jako map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	// dostane všechny soubory s názvem *.page.tmpl z ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// projeť všechny soubory končící na *.pages.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
