package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int
	Title  string
	Artist string
	Price  float64
}

/* Identificador de la ddbb */
var db *sql.DB

func main() {
	/* Capture connection properties */
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  os.Getenv("tcp"),
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	/* Get a database handle */
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("¡Connected!")

	/* Find records that match the passed parameter as a string */
	albums, err := AlbumsByArtist("Anyone")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	/* Hard-code ID 2 here to test the query */
	alb, err := AlbumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found album: %v\n", alb)
}

func AlbumsByArtist(name string) ([]Album, error) {
	/* An abums slice to hold data from returned rows */
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	/* Loop through rows, using Scan to assign column data to struct fields. */
	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	/* Check for an error from the overall query, using rows.Err */
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

/* Querys for the album with the specified ID */
func AlbumByID(id int64) (Album, error) {
	/* An album to hold data from the returned row */
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("AlbumByID: %d: not such album", id)
		}
		return alb, fmt.Errorf("AlbumByID: %d: %v", id, err)
	}
	return alb, nil
}
