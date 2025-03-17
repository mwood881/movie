package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "modernc.org/sqlite" // Import SQLite
)

func main() {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite", "movies.db")
	if err != nil {
		log.Fatal(err) //exit if fails
	}
	defer db.Close()

	// Create necessary tables if they don't exist
	createTables(db)

	// Load movie data from CSV into the database
	loadMovies(db, "IMDB-movies.csv")

	// Load genre data from CSV into the database
	loadGenres(db, "IMDB-movies_genres.csv")

	// Execute and display the top-rated movie genres
	executeSampleQuery(db)
}

// createTables creates the movies and genres tables in the SQLite database
func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY,
			title TEXT,
			year INTEGER,
			rating FLOAT
		);`,
		`CREATE TABLE IF NOT EXISTS genres (
			movie_id INTEGER,
			genre TEXT,
			FOREIGN KEY(movie_id) REFERENCES movies(id)
		);`,
	}
	// Execute each query to create tables if they don't exist
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err) // Log and exit if table creation fails
		}
	}
}

// loadMovies reads a CSV file and inserts movie data into the database
func loadMovies(db *sql.DB, filePath string) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err) // Log and exit if file opening fails
	}
	defer file.Close() // Ensure file is closed

	// Create a CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	//  inserting movie records
	stmt, err := db.Prepare("INSERT INTO movies (id, title, year, rating) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Ensure statement is closed

	// Iterate  skipping the header row
	for i, record := range records[1:] {
		if len(record) < 4 {
			log.Printf("Skipping row %d: not enough columns", i+1)
			continue
		}

		// Convert year from string to integer
		year, err := strconv.Atoi(record[2])
		if err != nil {
			log.Printf("Skipping row %d: invalid year format (%v)", i+1, err)
			continue
		}

		// Convert rating from string to float, handling NULL values
		rating := 0.0
		if record[3] != "NULL" && record[3] != "" {
			rating, err = strconv.ParseFloat(record[3], 64)
			if err != nil {
				log.Printf("Skipping row %d: invalid rating format (%v)", i+1, err)
				continue
			}
		}

		// Insert movie record into the database
		_, err = stmt.Exec(record[0], record[1], year, rating)
		if err != nil {
			log.Printf("Skipping row %d: %v", i+1, err)
		}
	}
}

// loadGenres reads a CSV file and inserts genre data into the database
func loadGenres(db *sql.DB, filePath string) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err) // Log and exit if file opening fails
	}
	defer file.Close() // Ensure file is closed when function exits

	// Create a CSV reader with flexible field counts and quote handling
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading genres CSV:", err)
	}

	// Prepare SQL statement for inserting genre records
	stmt, err := db.Prepare("INSERT INTO genres (movie_id, genre) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Ensure statement is closed when function exits

	// Iterate over CSV records, skipping the header row
	for i, record := range records[1:] {
		if len(record) < 2 {
			log.Printf("Skipping row %d: not enough columns (%v)", i+1, record)
			continue
		}

		// Insert genre record into the database
		_, err := stmt.Exec(record[0], record[1])
		if err != nil {
			log.Printf("Skipping row %d due to error: %v", i+1, err)
		}
	}
}

// retrieves and displays the top 10 highest-rated movie genres
func executeSampleQuery(db *sql.DB) {
	query := `SELECT g.genre, AVG(m.rating) AS avg_rating 
			FROM movies m
			JOIN genres g ON m.id = g.movie_id
			GROUP BY g.genre
			ORDER BY avg_rating DESC
			LIMIT 10;`

	// Execute the SQL query
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print the query results
	fmt.Println("Top 10 Highest Rated Genres:")
	for rows.Next() {
		var genre string
		var avgRating float64
		if err := rows.Scan(&genre, &avgRating); err != nil {
			log.Fatal(err) // Log and exit if data retrieval fails
		}
		fmt.Printf("%s: %.2f\n", genre, avgRating)
	}
}
