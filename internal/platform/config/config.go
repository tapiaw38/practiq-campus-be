package config

type (
	ServerConfig struct {
		AppName       string
		Port          string
		GinMode       string
		JWTSecret     string
		FrontendURL   string
		AuthAPIURL    string
		PractiqAPIURL string
	}

	DatabaseConfig struct {
		DatabaseURL string
	}

	// S3Config is optional: when empty, storage degrades to a no-op and
	// materials can still be created as external links.
	S3Config struct {
		AWSRegion          string
		AWSAccessKeyID     string
		AWSSecretAccessKey string
		AWSSessionToken    string
		AWSBucket          string
		AWSEndpoint        string
	}

	Config struct {
		ServerConfig   ServerConfig
		DatabaseConfig DatabaseConfig
		S3Config       S3Config
	}
)
