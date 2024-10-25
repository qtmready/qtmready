package main

type (
	Component struct {
		Name    string
		Service Service
	}
)

func NewComponent(name string, service Service) *Component {
	return &Component{
		Name:    name,
		Service: service,
	}
}
