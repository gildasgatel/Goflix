package config

const CREATE_TABLE_USERS = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user TEXT UNIQUE,
			pswd TEXT,
			account TEXT,
			name TEXT,
			firstname TEXT,
			mail TEXT UNIQUE, 
			cell INTEGER UNIQUE,
			adress TEXT
		);
	`

const CREATE_TABLE_MOVIES = `
		CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT UNIQUE,
			actors TEXT ,
			rating INTEGER,
			details TEXT,
			genre TEXT,
			saison INTEGER,
            episode INTEGER
		);
	`
