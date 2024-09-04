package request

type DemoParams struct {
	RequestString string `json:"requeststring"`
	RequestNumber int    `json:"requestnumber"`
}

type DemoPUTParams struct {
	Property1 string `json:"property1"`
	Property2 int    `json:"property2"`
}
