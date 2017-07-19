package main

import (
    "log"
    "html/template"
    "net/http"
)

type Option struct {
    Description string
    Link string
}

type Scene struct {
    Description string
    Options []Option
}

func main() {
    gameTemplate := template.Must(template.ParseFiles("template.html"))
    notFoundTemplate := template.Must(template.ParseFiles("404.html"))

    game := make(map[string]Scene)
    game[""] = Scene{Description: "You need to buy some milk at the grocery store. What transportation do you use?",
    Options: []Option{
        Option{Description: "Take the car", Link: "car"},
        Option{Description: "It's nice weather, take the bicycle this time", Link: "bicycle"},
        Option{Description: "Milk? Just stay at home instead.", Link: "lazy"},
    }}

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/html; charset=utf-8")
        w.Header().Add("Cache-Control", "no-cache")
        title := r.URL.Path[1:]
        scene, exists := game[title]
        if exists {
            gameTemplate.Execute(w, scene)
        } else {
            w.WriteHeader(404)
            notFoundTemplate.Execute(w, title)
        }
    })

    log.Println("gocyoa starting to serve...")
    http.ListenAndServe(":8080", nil)
}
