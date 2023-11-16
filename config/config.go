package config

type API struct {
	Name       string
	Port       string
	TagVersion string
	Env        string
	Host       string
}
type Database struct {
	Host string
	User string
	Pass string
	DB   string
	Port int64
}
