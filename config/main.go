// Package config defines the env configuration variables
package config

type Config struct {
	Port              string
	CookieSecret      string
	Database          string
	DatabaseHost      string
	DatabasePort      string
	DatabaseName      string
	DatabaseUsername  string
	DatabasePassword  string
	BaseURL           string
	SMTPUsername      string
	SMTPPassword      string
	SMTPHost          string
	SMTPPort          string
	SMTPSender        string
	RequestsPerMinute int
	CacheParameter    string
}
