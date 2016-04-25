CREATE DATABASE IF NOT EXISTS `quotes`;

CREATE USER "server"@"%" IDENTIFIED BY "password";

GRANT
	ALL PRIVILEGES ON `quotes`.*
	TO "server"@"%";

USE `quotes`;

CREATE TABLE `quote` (
	`id`      INT(11) NOT NULL AUTO_INCREMENT,
	`topic`   VARCHAR(30) NOT NULL,
	`author`  VARCHAR(255) NOT NULL,
	`text`    VARCHAR(767) NOT NULL,
	`created` DATETIME DEFAULT CURRENT_TIMESTAMP,
	`updated` DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	UNIQUE KEY `quote`  (`author`, `text`)
);

INSERT INTO `quote` (`topic`, `author`, `text`) VALUES
	("life", "Mahatma Gandhi", "Live as if you were to die tomorrow; learn as if you were to live forever."),
	("life", "Mahatma Gandhi", "The weak can never forgive. Forgiveness is the attribute of the strong."),
	("life", "Sun Tzu", "The supreme art of war is to subdue the enemy without fighting."),
	("science", "Albert Einstein", "Only two things are infinite, the universe and human stupidity, and I'm not sure about the former."),
	("science", "Albert Einstein", "To raise new questions, new possibilities, to regard old problems from a new angle, requires creative imagination and marks real advance in science."),
	("science", "Robert A. Heinlein", "Everything is theoretically impossible, until it is done."),
	("science", "Carl Sagan", "We live in a society exquisitely dependent on science and technology, in which hardly anyone knows anything about science and technology."),
	("science", "Isaac Asimov", "The saddest aspect of life right now is that science gathers knowledge faster than society gathers wisdom."),
	("science", "Edwin Powell Hubble", "Equipped with his five senses, man explores the universe around him and calls the adventure Science."),
	("science", "Henrik Ibsen", "It is inexcusable for scientists to torture animals; let them make their experiments on journalists and politicians."),
	("drinking", "Frank Sinatra", "Alcohol may be man's worst enemy, but the bible says love your enemy."),
	("drinking", "Sammy Davis Jr.", "Alcohol gives you infinite patience for stupidity."),
	("drinking", "Winston Churchill", "My rule of life prescribed as an absolutely sacred rite smoking cigars and also the drinking of alcohol before, after and if need be during all meals and in the intervals between them."),
	("drinking", "Terry Pratchett", "Death: \"THERE ARE BETTER THINGS IN THE WORLD THAN ALCOHOL, ALBERT.\"

Albert: \"Oh, yes, sir. But alcohol sort of compensates for not getting them.\""),
	("drinking", "W.C. Fields", "I cook with wine, sometimes I even add it to the food."),
	("drinking", "Bette Davis", "There comes a time in every woman's life when the only thing that helps is a glass of champagne."),
	("drinking", "Dorothy Parker", "I'd rather have a bottle in front of me than a frontal lobotomy."),
	("computers", "Roger Ebert", "Doing research on the Web is like using a library assembled piecemeal by pack rats and vandalized nightly."),
	("computers", "Charles Stross", "Idiots emit bogons, causing machinery to malfunction in their presence. System administrators absorb bogons, letting machinery work again."),
	("computers", "Dr. Diogo Monica", "To be fair, faking GPG usage is almost as hard as actually using GPG."),
	("computers", "Edsger W. Dijkstra", "The use of COBOL cripples the mind; its teaching should, therefore, be regarded as a criminal offense"),
	("computers", "Arthur C. Clarke", "Any sufficiently advanced technology is indistinguishable from magic."),
	("computers", "Douglas Adams", "We are stuck with technology when what we really want is just stuff that works."),
	("computers", "Philip K. Dick", "There will come a time when it isn't 'They're spying on me through my phone' anymore. Eventually, it will be 'My phone is spying on me'.");
