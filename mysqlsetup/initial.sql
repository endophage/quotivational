CREATE DATABASE IF NOT EXISTS `quotes`;

CREATE USER "server"@"%" IDENTIFIED BY "password";

GRANT
	ALL PRIVILEGES ON `quotes`.*
	TO "server"@"%";

USE `quotes`;

CREATE TABLE `quote` (
	`id`      INT(11) NOT NULL AUTO_INCREMENT,
	`author`  VARCHAR(255) NOT NULL,
	`text`    VARCHAR(767) NOT NULL,
	`created` DATETIME NOT NULL,
	`updated` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY `quote`  (`author`, `text`)
);

INSERT INTO `quote` (`author`, `text`, `created`, `updated`) VALUES (
	"Mahatma Gandhi",
	"Live as if you were to die tomorrow; learn as if you were to live forever.",
	CURDATE(),
	CURDATE()
);

INSERT INTO `quote` (`author`, `text`, `created`, `updated`) VALUES (
	"Mahatma Gandhi",
	"The weak can never forgive. Forgiveness is the attribute of the strong.",
	CURDATE(),
	CURDATE()
);

INSERT INTO `quote` (`author`, `text`, `created`, `updated`) VALUES (
	"Sun Tzu",
	"The supreme art of war is to subdue the enemy without fighting.",
	CURDATE(),
	CURDATE()
);
