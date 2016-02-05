package main

import (
	"math/rand"
)

var quotes = map[string][]Quote{
	"science": {
		{Quote: "Only two things are infinite, the universe and human stupidity, and I'm not sure about the former.", Author: "Albert Einstein"},
		{Quote: "To raise new questions, new possibilities, to regard old problems from a new angle, requires creative imagination and marks real advance in science.", Author: "Albert Einstein"},
		{Quote: "Everything is theoretically impossible, until it is done.", Author: "Robert A. Heinlein"},
		{Quote: "We live in a society exquisitely dependent on science and technology, in which hardly anyone knows anything about science and technology.", Author: "Carl Sagan"},
		{Quote: "The saddest aspect of life right now is that science gathers knowledge faster than society gathers wisdom.", Author: "Isaac Asimov"},
		{Quote: "Equipped with his five senses, man explores the universe around him and calls the adventure Science.", Author: "Edwin Powell Hubble"},
		{Quote: "It is inexcusable for scientists to torture animals; let them make their experiments on journalists and politicians.", Author: "Henrik Ibsen"},
	},
	"drinking": {
		{Quote: "Alcohol may be man's worst enemy, but the bible says love your enemy.", Author: "Frank Sinatra"},
		{Quote: "Alcohol gives you infinite patience for stupidity.", Author: "Sammy Davis Jr."},
		{Quote: "My rule of life prescribed as an absolutely sacred rite smoking cigars and also the drinking of alcohol before, after and if need be during all meals and in the intervals between them.", Author: "Winston Churchill"},
		{Quote: `Death: "THERE ARE BETTER THINGS IN THE WORLD THAN ALCOHOL, ALBERT."

Albert: "Oh, yes, sir. But alcohol sort of compensates for not getting them."`, Author: "Terry Pratchett"},
		{Quote: "I cook with wine, sometimes I even add it to the food.", Author: "W.C. Fields"},
		{Quote: "There comes a time in every woman's life when the only thing that helps is a glass of champagne.", Author: "Bette Davis"},
		{Quote: "I'd rather have a bottle in front of me than a frontal lobotomy.", Author: "Dorothy Parker"},
	},
	"computers": {
		{Quote: "Doing research on the Web is like using a library assembled piecemeal by pack rats and vandalized nightly.", Author: "Roger Ebert"},
		{Quote: "Idiots emit bogons, causing machinery to malfunction in their presence. System administrators absorb bogons, letting machinery work again.", Author: "Charles Stross"},
		{Quote: "To be fair, faking GPG usage is almost as hard as actually using GPG.", Author: "Dr. Diogo Monica"},
		{Quote: "The use of COBOL cripples the mind; its teaching should, therefore, be regarded as a criminal offense", Author: "Edsger W. Dijkstra"},
		{Quote: "Any sufficiently advanced technology is indistinguishable from magic.", Author: "Arthur C. Clarke"},
		{Quote: "We are stuck with technology when what we really want is just stuff that works.", Author: "Douglas Adams"},
		{Quote: "There will come a time when it isn't 'They're spying on me through my phone' anymore. Eventually, it will be 'My phone is spying on me'.", Author: "Philip K. Dick"},
	},
}

type MockQuoter struct{}

func (m *MockQuoter) Quote(topic string) (Quote, error) {
	t := quotes[topic]
	return t[rand.Intn(len(t))], nil
}
