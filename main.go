package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Option struct {
	Description string
	Link        string
}

type Scene struct {
	Description string
	Options     []Option
}

func main() {
	game := make(map[string]Scene)
	game[""] = Scene{"You need to buy some milk at the grocery store. What transportation do you use?",
		[]Option{
			Option{"Take the car.", "car"},
			Option{"It's nice weather, take the bicycle this time.", "bicycle"},
			Option{"Milk? Just stay at home instead.", "lazy"},
		}}
	game["car"] = Scene{"You safely arrive at the store. Your car has always been a reliable one.",
		[]Option{
			Option{"Ok", "store"},
		}}
	game["store"] = Scene{"You enter the store. Do you look around a bit or go directly to the milk?",
		[]Option{
			Option{"Look at what else is in the store.", "around"},
			Option{"Go straight for the milk.", "milk"},
		}}
	game["around"] = Scene{"You see an advertisement on the wall in the store. It seems a local crazy scientist is looking for a tester for his \"time modification device\". \"Looking someone with a weak sense of self-preservation who wants to make some easy money in groundbreaking experiments!\"",
		[]Option{
			Option{"Ignore it and go for the milk.", "milk"},
			Option{"Sounds good. Apply as tester, get the milk later.", "time"},
		}}
	game["milk"] = Scene{"You get some cheap milk, your friend John told you the expensive stuff is just repackaged from the same factory anyway. That was only after he checked his wallet though. John is always struck by wisdom after checking his wallet. You wish you had a gift like that. Stepping out of your thoughts, you find the line at the counter is incredibly long. What do you do?",
		[]Option{
			Option{"Wait in line like a good customer.", "line"},
			Option{"Try to scare away some of the other customers in line.", "scare"},
			Option{"Just walk outside without paying.", "theft"},
		}}
	game["line"] = Scene{"You're still waiting in line. What do you do?",
		[]Option{
			Option{"Wait in line more like a good customer.", "line2"},
			Option{"Try to scare away some of the other customers in line.", "scare"},
			Option{"Just walk outside without paying.", "theft"},
		}}
	game["line2"] = Scene{"You're STILL waiting in line. Is that grandma trying to pay her groceries in pennies?",
		[]Option{
			Option{"Wait in line even more like a good sheep.", "line3"},
			Option{"Try to scare away some of the other customers in line.", "scare"},
			Option{"Just walk outside without paying.", "theft"},
		}}
	game["line3"] = Scene{"You're almost at the end of the line. But you've been waiting long enough. And on top of that it seems that the fat guy in front of you is moving extra slowly. You're asking yourself whether he's doing it on purpose. Or maybe God is doing it on purpose to teach you something. You wonder what it is God would want to teach you. After all you had felt you learned a lot already and there couldn't be that much left to be teached. You decide it must be the fat guy that's doing it on purpose.",
		[]Option{
			Option{"Just keep waiting.", "line4"},
			Option{"Try to scare this guy away.", "scare"},
			Option{"Just walk outside without paying.", "theft"},
		}}
	game["line4"] = Scene{"FINALLY the last customer goes away and you're at the counter. Indian shopkeeper Amesh greets you with an extremely strong accent. Then you find out you forgot your wallet at home. You have no cash and no card. What now?",
		[]Option{
			Option{"Try to haggle for the milk.", "haggle"},
			Option{"Give back the milk and go back home empty-handed.", "home"},
			Option{"Just walk outside without paying.", "theft"},
		}}
	game["theft"] = Scene{"You try to walk outside, but the Indian shopkeeper Amesh catches you. He runs away from the counter towards the exit with an M4 assault rifle. He says \"Drop the milk!\" with a ridiculously overdone accent.",
		[]Option{
			Option{"Give back the milk and go back home empty-handed.", "home"},
			Option{"Run for it!", "run"},
			Option{"Argue that his weapon is just as illegal as your theft.", "argue"},
		}}
	game["run"] = Scene{"You find out the hard way that you can't outrun bullets, especially the M4 in three-round burst fire mode. the 5.56x45mm NATO, air-cooled, direct impinged gas-operated, magazine-fed carbine fires multiple bursts into your body. As you fall and gasp for air, you feel the life slipping out of you. The last words you hear are \"please come again\". You are dead.",
		[]Option{
			Option{"Try again", ""},
		}}
	game["argue"] = Scene{"You remember the book you recently read about debating and influencing people and decide what works for a used car salesman to sell even more in 5 easy steps must also work on convincing angry, gun-wielding shopkeepers not to sell you bullets. You try to argue using the techniques of your book but you feel your arguments become increasingly philosphical. In fact you find that Amesh has a good rebuttal to your arguments and cannot find any contradiction in his rationality. You feel increasingly convinced that you should be shot.",
		[]Option{
			Option{"Give back the milk and go back home empty-handed.", "home"},
			Option{"Run for it!", "run"},
			Option{"Try to counter his arguments.", "counter"},
			Option{"Agree with him that you deserve to be shot.", "agree"},
		}}
	game["agree"] = Scene{"It turns out Amesh is also a master philosopher and convincing speaker. You accept that by all logic you deserve to be shot. It's all he can do now, really. You feel a bit bad that Amesh needs to go through this. Not only is Amesh good at convincing, he's also good at shooting, having been in the military for 20 years prior to retiring to his civil life as a simple shopkeep. He fires a single shot. The last thing that goes through your mind is the Ren & Stimpy cartoon you saw this morning. How unfortunate. You are dead.",
		[]Option{
			Option{"Try again", ""},
		}}

	gameTemplate := template.Must(template.ParseFiles("template.html"))
	notFoundTemplate := template.Must(template.ParseFiles("404.html"))
	homePage, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/play/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Add("Cache-Control", "no-cache")
		title := r.URL.Path[len("/play/"):]
		scene, exists := game[title]
		if exists {
			gameTemplate.Execute(w, scene)
		} else {
			w.WriteHeader(404)
			notFoundTemplate.Execute(w, title)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(homePage)
	})

	log.Println("gocyoa starting to serve...")
	http.ListenAndServe(":8080", nil)
}
