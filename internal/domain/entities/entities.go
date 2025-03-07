package entities

type Server struct {
	Id        int64
	Name      string
	Type      string
	Databases int
	IpAddress string
	Port      string `fake:"{number:1024,49151}"`
	Status    string
	Username  string
	Password  string
}

type Database struct {
	Id     int64
	Name   string
	Server string
}
