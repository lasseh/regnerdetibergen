package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("index.tmpl"))

type status struct {
	Message string `json:"message"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := indexTemplate.Execute(&buf, nil)
	if err != nil {
		log.Println("Index template rendering failed: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(buf.Bytes())
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	rain := GetRain()
	s := status{Message: rain}
	json.NewEncoder(w).Encode(s)
}
