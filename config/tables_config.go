package config

const CREATE_TABLE_USERS = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user TEXT,
			pswd TEXT,
			account TEXT,
			name TEXT,
			firstname TEXT,
			mail TEXT UNIQUE, 
			cell INTEGER UNIQUE,
			adress TEXT
		);
	`
