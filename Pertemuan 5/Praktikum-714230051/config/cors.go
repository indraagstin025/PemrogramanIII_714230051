package config

var allowedOrigins = []string{
	"https://indrariksa.github.io/",
}

func GetAllowedOrigins() []string {
	return allowedOrigins
}
