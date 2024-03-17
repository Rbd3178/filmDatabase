package entities

// ActorWithFilm
type ActorWithFilm struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Gender string `db:"gender"`
	BirthDate string `db:"birth_date"`
	FilmId int `db:"film_id"`
	Title string `db:"title"`
}

// Actor
type Actor struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Gender string `db:"gender"`
	BirthDate string `db:"birth_date"`
}