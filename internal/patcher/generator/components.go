package generator

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

type DashboardGroup struct {
	Sliders      []Slider
	TextFields   []TextField
	NumberFields []NumberField
	Group        string
}

type dashboardComponents map[string]DashboardGroup

func (components dashboardComponents) toSlice() []DashboardGroup {
	var componentsSlice []DashboardGroup

	for groupName, dashboardGroup := range components {
		dashboardGroup.Group = groupName
		componentsSlice = append(componentsSlice, dashboardGroup)
	}

	return componentsSlice
}
