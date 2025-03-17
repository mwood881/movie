# Movie Database

## Description:
This repository serves the purpose of creating a movie database using SQLite. CSV files are used to then load data and create tables for genres and movie titles. Then the top rated genres are shown as the top 10. 

## Features:
- Movie data from IMDB-movies.csv
- Genre data from IMDB-movies_genres.csv
- Top 10 highest rated genres

## Installing the Package for Use
(1) Clone repository: 


   ```bash
   git clone https://github.com/mwood881/movie_database.git
cd movie_database
   ```


(2) Install package

  ```bash
  go get modernc.org/sqlite
   ```

(3) Run the program

  ```bash
go run main.go
   ```

## Results
Top 10 Highest Rated Genres:
Film-Noir: 6.47
Thriller: 2.73
Horror: 2.63
War: 2.52
Sci-Fi: 2.47
Family: 2.44
Romance: 2.34
Mystery: 2.31
Adventure: 2.23
Fantasy: 2.19

## Future Improvements
You could use more tables to track movie reviews and watch history. 
## Resources Used
I used github copilot to help clean up errors in my code. mostly for learning how to upload this repository to github. 
