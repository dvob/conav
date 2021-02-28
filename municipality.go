package conav

//go:generate go run municipality_list_gen.go

type Municipality struct {
	ID         int
	Name       string
	Population int
	Canton     string
}
