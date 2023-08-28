package entity

type User struct {
	Id   int    `db:"id"`
	Slug string `db:"slug"`
}
