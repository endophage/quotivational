package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

var topic = "science"

func setupWindow(title string, width, height int) (*gtk.Window, error) {
	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	w.SetTitle(title)
	w.Connect("destroy", func() {
		gtk.MainQuit()
	})
	w.SetDefaultSize(width, height)
	return w, nil
}

func setupMenuBar(g *gtk.Grid) error {
	bar, err := gtk.MenuBarNew()
	if err != nil {
		return err
	}
	g.Add(bar)
	topics, err := gtk.MenuItemNewWithLabel("Topics")
	if err != nil {
		return err
	}

	topicList, err := gtk.MenuNew()
	for t, _ := range quotes {
		s, err := gtk.MenuItemNewWithLabel(t)
		if err != nil {
			return err
		}
		fixedT := t
		s.Connect("activate", func() {
			topic = fixedT
		})
		topicList.Append(s)
	}

	topics.SetSubmenu(topicList)
	bar.Append(topics)
	return nil
}

func setupWidgets(w *gtk.Window) error {
	grid, err := gtk.GridNew()
	if err != nil {
		return err
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	err = setupMenuBar(grid)
	if err != nil {
		return err
	}

	quote, err := gtk.LabelNew("Click the button!")
	if err != nil {
		return err
	}
	quote.SetSizeRequest(400, 100)
	quote.SetLineWrapMode(pango.WRAP_WORD)
	quote.SetLineWrap(true)
	quote.SetMarginTop(20)
	quote.SetMarginStart(20)
	quote.SetMarginEnd(20)

	author, err := gtk.LabelNew("")
	if err != nil {
		return err
	}
	author.SetSizeRequest(400, 50)

	b, err := gtk.ButtonNewWithLabel("Quotivate Me!")
	if err != nil {
		return err
	}
	b.SetSizeRequest(100, 50)
	b.SetMarginStart(50)
	b.SetMarginEnd(50)
	b.SetMarginBottom(30)
	q := MockQuoter{}
	b.Connect("clicked", func() {
		s, err := q.Quote(topic)
		if err != nil {
			quote.SetLabel("couldn't get a quote")
			author.SetLabel("")
			return
		}
		quote.SetLabel(s.Quote)
		author.SetLabel("- " + s.Author)
	})

	grid.Add(quote)
	grid.Add(author)
	grid.Add(b)
	w.Add(grid)
	return nil
}
