package model

type Test struct {
	Id        int `db:"id" json:"id"`
	Title     string  `db:"title" json:"title"`
	Type      int `db:"f_type" json:"type"`
}
