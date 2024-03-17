package entities

// FilmWithActor
type FilmWithActor struct {
	ID          int     `db:"id"`
	Title       string  `db:"title"`
	Description string  `db:"description"`
	ReleaseDate string  `db:"release_date"`
	Rating      float64 `db:"rating"`
	ActorID     *int    `db:"actor_id"`
	Name        *string `db:"name"`
}

// Actor
/*type Film struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Gender string `db:"gender"`
	BirthDate string `db:"birth_date"`
}*/
