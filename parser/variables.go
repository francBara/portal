package parser

type NumberVariable struct {
	Name  string
	Value int
	Max   int
	Min   int
	Step  int
}

type StringVariable struct {
	Name  string
	Value string
}

type PortalVariables struct {
	Number []NumberVariable
	String []StringVariable
}

func (pv PortalVariables) Concat(newPv PortalVariables) PortalVariables {
	var concatenatedPv PortalVariables
	concatenatedPv.Number = append(pv.Number, newPv.Number...)
	concatenatedPv.String = append(pv.String, newPv.String...)
	return concatenatedPv
}
