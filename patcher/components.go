package patcher

type DashboarComponents struct {
	Sliders      []Slider
	TextFields   []TextField
	NumberFields []NumberField
}

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

type NumberField struct {
	Id           string
	InitialValue int
	Name         string
}
