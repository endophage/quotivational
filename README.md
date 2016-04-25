# Quotivational

This is the best Linux app for getting motivational quotes.

Instructions to build:

- Install requirements:
	- MySQL (Maria DB version >=10.1.10) up and running
	- Redis (version >= 3.0) up and running
	- Golang >= 1.5
	- GTK development packages for your OS
		- Linux: libgtk+-2.0
		- Windows: see http://www.gtk.org/download/windows.php
		- Mac Os: gtk+3

- `git clone https://github.com/endophage/quotivational.git ~/Go/src/github.com/endophage/quotivational`
- `cd ~/Go/src/github.com/endophage/quotivational`
- `go build -o bin/quotivational ./cmd/quotivational`
- `go build -o bin/quotivational-server ./cmd/server`
- `go build -o bin/quotivational-auth ./cmd/auth`
- `mysql < mysqlsetup/initial.sql`
- In one terminal:  `bin/quotivational-auth gooduser <user1> <user2>`
- In another terminal: `bin/quotivational-server -db server:password@tcp(localhost:3306)/quotes?parseTime=true -auth http://localhost:8081 -redis localhost:6379`
