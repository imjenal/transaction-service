package config

type (
	Environment string

	//Server has all the server related config
	Server struct {
		Port        int         `validate:"required"`
		Address     string      `validate:"required"`
		Environment Environment `validate:"required,oneof=PROD STAGING DEV TEST"`
	}

	DB struct {
		Host     string `validate:"required,hostname_rfc1123"`
		Port     string `validate:"required,number"`
		User     string `validate:"required"`
		Password string `validate:"required"`
		Name     string `validate:"required"`
	}
)
