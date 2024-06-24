package get_handlers

import (
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"snake_ai/pkg/logger"
)

func PlayerMatch(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	dir, _ := filepath.Split(os.Args[0])
	filePath := filepath.Join(dir, "internal/templates/play.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
