package patcher

type Slider struct {
	Id           string
	InitialValue int
	Min          int
	Max          int
	Step         int
	Name         string
}

type TextField struct {
	Id           string
	InitialValue string
	Name         string
}
