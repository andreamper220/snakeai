package get_handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func UserAuth(w http.ResponseWriter, r *http.Request) {
	dir, _ := filepath.Split(os.Args[0])
	filePath := filepath.Join(dir, "internal/game/templates/auth.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
