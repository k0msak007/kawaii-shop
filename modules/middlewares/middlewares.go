package middlewares

type Role struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}
