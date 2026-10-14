package models

type Category struct {
	ID          int    `bson:"id" json:"id"`
	Description string `bson:"description" json:"description"`
}

type Genre struct {
	ID          int    `bson:"id" json:"id"`
	Description string `bson:"description" json:"description"`
}

type Game struct {
	ID              int        `bson:"_id" json:"id"`
	Title           string     `bson:"title" json:"title"`
	RequiredAge     int        `bson:"required_age" json:"required_age"`
	IsFree          bool       `bson:"is_free" json:"is_free"`
	Description     string     `bson:"description" json:"description"`
	Image           string     `bson:"image" json:"image"`
	Website         string     `bson:"website" json:"website"`
	Developers      []string   `bson:"developers" json:"developers"`
	Publishers      []string   `bson:"publishers" json:"publishers"`
	Platforms       []string   `bson:"platforms" json:"platforms"`
	MetacriticScore int        `bson:"metacritic_score" json:"metacritic_score"`
	Categories      []Category `bson:"categories" json:"categories"`
	Genres          []Genre    `bson:"genres" json:"genres"`
	Recommendations int        `bson:"recommendations" json:"recommendations"`
	ReleaseDate     string     `bson:"release_date" json:"release_date"`
	ReleaseYear     int        `bson:"release_year" json:"release_year"`
}