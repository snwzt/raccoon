package models

type Status struct {
	Path      string
	FinalSize int64
	Parts     []int64
	Done      bool
	Err       error
}

type ReadData struct {
	N   int
	Err error
}
