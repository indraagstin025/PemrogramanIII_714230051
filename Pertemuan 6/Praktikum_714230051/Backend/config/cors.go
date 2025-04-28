package config

var allowedOrigins = []string{
	"https://indrariksa.github.io/",
	"http://localhost:5173/",
}

func GetAllowedOrigins() []string {
	return allowedOrigins
}
