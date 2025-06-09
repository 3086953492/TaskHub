package configs

import "net/http"

type SessionConfig struct {
	HttpOnly bool          `mapstructure:"http_only" json:"http_only"`
	Secure   bool          `mapstructure:"secure" json:"secure"`
	SameSite http.SameSite `mapstructure:"same_site" json:"same_site"`
	MaxAge   int           `mapstructure:"max_age" json:"max_age"`
	Path     string        `mapstructure:"path" json:"path"`
	Domain   string        `mapstructure:"domain" json:"domain"`
}
