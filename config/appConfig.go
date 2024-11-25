package config

type AppConfig struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	FromEmail         string
	FromEmailSMTP     string
	FromEmailPassword string
	SMTPAddrress      string
}
