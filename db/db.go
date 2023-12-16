package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"goflix/config"
	"goflix/models"
	"goflix/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Setup() error
	Close()
	GetUser(int) (*models.User, error)
	SaveUser(*models.User) error
	DeleteUser(int) error
	UpdateUser(*models.User) error
	GetID(user *models.User) error
}

type DbSqlite struct {
	sqlite *sql.DB
}

func New() Storage {
	return &DbSqlite{}
}

func (db *DbSqlite) Setup() error {
	var err error
	db.sqlite, err = sql.Open(config.DRIVE_NAME, config.DATA_SOURCE_NAME)
	if err != nil {
		return err
	}
	err = db.sqlite.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Database connected!")
	_, err = db.sqlite.Exec(config.CREATE_TABLE_USERS)
	if err != nil {
		return err
	}

	return nil
}

func (db *DbSqlite) Close() {
	db.sqlite.Close()
}

func (db *DbSqlite) GetUser(id int) (*models.User, error) {
	rows, err := db.sqlite.Query("SELECT * FROM users  WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user models.User
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.User, &user.Pswd, &user.Account,
			&user.Info.Name,
			&user.Info.Firstname,
			&user.Info.Mail,
			&user.Info.Cell,
			&user.Info.Adress)
		if err != nil {
			return nil, err
		}
	}
	if user.Id == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (db *DbSqlite) SaveUser(user *models.User) error {
	hashPswd, err := utils.HashPasswd([]byte(user.Pswd))
	if err != nil {
		return err
	}
	insertSQL := "INSERT INTO users (user,pswd,account,name,firstname,mail,cell,adress) VALUES (?,?,?,?,?,?,?,?)"
	_, err = db.sqlite.Exec(insertSQL,
		user.User,
		hashPswd,
		user.Account,
		user.Info.Name,
		user.Info.Firstname,
		user.Info.Mail,
		user.Info.Cell,
		user.Info.Adress)
	if err != nil {
		return err
	}

	return nil
}
func (db *DbSqlite) GetID(user *models.User) error {
	pswd := user.Pswd
	rows, err := db.sqlite.Query("SELECT * FROM users  WHERE user=? ", user.User)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.User, &user.Pswd, &user.Account,
			&user.Info.Name,
			&user.Info.Firstname,
			&user.Info.Mail,
			&user.Info.Cell,
			&user.Info.Adress)
		if err != nil {
			return err
		}
	}
	if user.Id == 0 {
		return errors.New("user not found")
	}
	err = utils.CompareHashAndPassword([]byte(pswd), []byte(user.Pswd))
	if err != nil {
		return errors.New("error password")
	}
	return nil
}

func (db *DbSqlite) UpdateUser(user *models.User) error {
	updateSQL := "UPDATE users SET user =? AND pswd=? AND account=? AND name=? AND firstname=? AND mail=? AND cell=? AND adress=? WHERE id = ?"
	_, err := db.sqlite.Exec(updateSQL,
		user.User,
		strings.TrimSpace(user.Pswd),
		user.Account,
		user.Info.Name,
		user.Info.Firstname,
		user.Info.Mail,
		user.Info.Cell,
		user.Info.Adress,
		user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DbSqlite) DeleteUser(id int) error {
	deleteSQL := "DELETE FROM users WHERE id = ?"
	result, err := db.sqlite.Exec(deleteSQL, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("errors want delete 1 reccord got: %d", rowsAffected)
	}

	return nil
}
