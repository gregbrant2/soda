package dtos

type Server struct {
	Id        int64
	Name      string
	Type      string
	Databases int
	IpAddress string
	Port      string
	Status    string
}
