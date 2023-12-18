package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	GetMoviesById(id int) (*models.Movies, error)
	GetMovies() ([]*models.Movies, error)
	DeleteMovieByID(id int) error
	GetSeries() ([]*models.Movies, error)
	AddMovie(movie *models.Movies) error
	SaveRating(rating *models.Rating) error
	GetFavoriteByUser(favorite *models.Favorite) error
	SaveFavorite(favorite *models.Favorite) error
	GetRatingByUser(id int) ([]*models.Rating, error)
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

	err = db.InitTables()
	if err != nil {
		return err
	}
	fmt.Println("Database connected!")

	return nil
}

func (db *DbSqlite) InitTables() error {
	_, err := db.sqlite.Exec(config.CREATE_TABLE_USERS)
	if err != nil {
		return err
	}
	fmt.Println("user created!")
	_, err = db.sqlite.Exec(config.CREATE_TABLE_MOVIES)
	if err != nil {
		return err
	}
	fmt.Println("movies created!")
	_, err = db.sqlite.Exec(config.CREATE_TABLE_FAVORITE)
	if err != nil {
		return err
	}
	fmt.Println("favorite created!")

	_, err = db.sqlite.Exec(config.CREATE_TABLE_RATING)
	if err != nil {
		return err
	}
	fmt.Println("rating created!")

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
	hashPswd, err := utils.HashPasswd([]byte(user.Pswd))
	if err != nil {
		return err
	}
	updateSQL := "UPDATE users SET user = ? , pswd = ? , account = ? , name = ? , firstname = ? , mail = ? , cell = ? , adress = ? WHERE id = ?"
	res, err := db.sqlite.Exec(updateSQL,
		&user.User,
		&hashPswd,
		&user.Account,
		&user.Info.Name,
		&user.Info.Firstname,
		&user.Info.Mail,
		&user.Info.Cell,
		&user.Info.Adress,
		&user.Id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); n < 1 || err != nil {
		return fmt.Errorf("update failed with id: %d", user.Id)
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

// * * *

func (db *DbSqlite) GetMoviesById(id int) (*models.Movies, error) {
	rows, err := db.sqlite.Query("SELECT * FROM movies  WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movie models.Movies
	if rows.Next() {
		err = rows.Scan(&movie.Id,
			&movie.Title,
			&movie.Actors,
			&movie.Rating,
			&movie.Details,
			&movie.Genre,
			&movie.Saison,
			&movie.Episode)
		if err != nil {
			return nil, err
		}
	}
	if movie.Id == 0 {
		return nil, errors.New("movie not found")
	}
	return &movie, nil
}
func (db *DbSqlite) GetMovies() ([]*models.Movies, error) {
	rows, err := db.sqlite.Query("SELECT * FROM movies WHERE saison = 0 ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []*models.Movies
	for rows.Next() {
		movie := models.Movies{}
		err = rows.Scan(&movie.Id,
			&movie.Title,
			&movie.Actors,
			&movie.Rating,
			&movie.Details,
			&movie.Genre,
			&movie.Saison,
			&movie.Episode)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if len(movies) == 0 {
		return nil, errors.New("movies not found")
	}
	return movies, nil
}
func (db *DbSqlite) DeleteMovieByID(id int) error {
	deleteSQL := "DELETE FROM movies WHERE id = ?"
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
func (db *DbSqlite) GetSeries() ([]*models.Movies, error) {
	rows, err := db.sqlite.Query("SELECT * FROM movies WHERE saison > 0 ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var series []*models.Movies
	for rows.Next() {
		serie := models.Movies{}
		err = rows.Scan(&serie.Id,
			&serie.Title,
			&serie.Actors,
			&serie.Rating,
			&serie.Details,
			&serie.Genre,
			&serie.Saison,
			&serie.Episode)
		if err != nil {
			return nil, err
		}
		series = append(series, &serie)
	}
	if len(series) == 0 {
		return nil, errors.New("series not found")
	}
	return series, nil
}
func (db *DbSqlite) AddMovie(movie *models.Movies) error {
	insertSQL := "INSERT INTO movies (title, actors, rating, details, genre, saison, episode) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.sqlite.Exec(insertSQL,
		movie.Title, movie.Actors, movie.Rating, movie.Details, movie.Genre, movie.Saison, movie.Episode)
	if err != nil {
		return err
	}

	return nil
}

func (db *DbSqlite) SaveFavorite(favorite *models.Favorite) error {

	updateSQL := "UPDATE favorite SET moviesid = ? WHERE userid = ?"
	res, err := db.sqlite.Exec(updateSQL, &favorite.MoviesID, &favorite.UserId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		log.Println("favorite updated")
		return nil
	}
	insertSQL := "INSERT INTO favorite (userid, moviesid) VALUES (?, ?)"
	_, err = db.sqlite.Exec(insertSQL, &favorite.UserId, &favorite.MoviesID)
	if err != nil {
		return err
	}
	log.Println("favorite add")
	return nil
}
func (db *DbSqlite) GetFavoriteByUser(favorite *models.Favorite) error {
	rows, err := db.sqlite.Query("SELECT * FROM favorite WHERE userid = ?", favorite.UserId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&favorite.UserId, &favorite.MoviesID)
		if err != nil {
			return err
		}
	}
	return nil
}
func (db *DbSqlite) SaveRating(rating *models.Rating) error {
	updateSQL := "UPDATE rating SET stars = ? WHERE movieid = ? AND userid = ?"
	res, err := db.sqlite.Exec(updateSQL, &rating.Stars, &rating.MovieId, &rating.UserId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		log.Println("rating updated")
		return nil
	}
	insertSQL := "INSERT INTO rating (movieid, stars, userid) VALUES (?, ?, ?)"
	_, err = db.sqlite.Exec(insertSQL, &rating.MovieId, &rating.Stars, &rating.UserId)
	if err != nil {
		return err
	}
	log.Println("rating add")
	return nil
}
func (db *DbSqlite) GetRatingByUser(id int) ([]*models.Rating, error) {
	rows, err := db.sqlite.Query("SELECT * FROM rating WHERE userid = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ranting []*models.Rating
	for rows.Next() {
		newRating := models.Rating{}
		err = rows.Scan(&newRating.MovieId, &newRating.Stars, &newRating.UserId)
		if err != nil {
			return nil, err
		}
		ranting = append(ranting, &newRating)
	}
	return ranting, nil
}
