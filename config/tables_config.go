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
			title TEXT,
			actors TEXT ,
			rating INTEGER,
			details TEXT,
			genre TEXT,
			saison INTEGER,
            episode INTEGER
		);
	`
const CREATE_TABLE_FAVORITE = `
	CREATE TABLE IF NOT EXISTS favorite (
		userid INTEGER PRIMARY KEY ,
		moviesid TEXT
	);
`
const CREATE_TABLE_RATING = `
CREATE TABLE IF NOT EXISTS rating (
	moviesid INTEGER PRIMARY KEY ,
	stars INTEGER,
	userid INTEGER
);
`
