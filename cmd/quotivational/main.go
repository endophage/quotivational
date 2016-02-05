package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	w, err := setupWindow("Quotivational", 400, 175)
	if err != nil {
		log.Fatal(err)
	}
	err = setupWidgets(w)
	if err != nil {
		log.Fatal(err)
	}

	// Recursively show all widgets contained in this window.
	w.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
