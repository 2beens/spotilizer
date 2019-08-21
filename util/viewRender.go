package util

import (
	"fmt"
	"github.com/2beens/spotilizer/models"
	"github.com/prometheus/common/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

var templatesMap map[string]*template.Template

func SetupTemplates() {
	templatesMap = make(map[string]*template.Template)
	layoutFiles := []string{
		"public/views/layouts/layout.html",
		"public/views/layouts/footer.html",
		"public/views/layouts/navbar.html",
	}

	viewFileNames, err := getViewFileNames()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range viewFileNames {
		viewPath := "public/views/" + v
		t, err := template.New("layout").ParseFiles(append(layoutFiles, viewPath)...)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Infof(" > read template view file: " + viewPath)
		templatesMap[v] = t
	}
}

func getViewFileNames() ([]string, error) {
	var viewFileNames []string
	viewFiles, err := ioutil.ReadDir("./public/views")
	if err != nil {
		return viewFileNames, err
	}
	for _, f := range viewFiles {
		fileName := f.Name()
		if strings.HasSuffix(fileName, ".html") {
			viewFileNames = append(viewFileNames, fileName)
		}
	}
	return viewFileNames, nil
}

// templates cheatsheet
// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet
func RenderView(w http.ResponseWriter, page string, viewData interface{}) {
	t, ok := templatesMap[page+".html"]
	if !ok {
		log.Error(" >>> error rendering view, cannot find view template: " + page + ".html")
		http.Error(w, "internal server error (error rendering view)", http.StatusInternalServerError)
	}

	err := t.ExecuteTemplate(w, "layout", viewData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderSpAPIErrorView(w http.ResponseWriter, username string, title string, apiErr *models.SpAPIError) {
	RenderView(w, "error", models.ErrorViewData{Title: title, Error: fmt.Sprintf("Status: [%d]: %s", apiErr.Error.Status, apiErr.Error.Message), Username: username})
}

func RenderErrorView(w http.ResponseWriter, username string, title string, status int, message string) {
	RenderView(w, "error", models.ErrorViewData{Title: title, Error: fmt.Sprintf("Status: [%d]: %s", status, message), Username: username})
}
