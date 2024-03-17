package store

// Store
type Store interface {
	User() UserRepository
	Film() FilmRepository
	Actor() ActorRepository
}
