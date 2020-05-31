package settings

import (
	"os"
)

func getEnv(envName, envDefault string) string {
	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}

	return envDefault
}

func GetLogLevel() string {
	return getEnv("LOG_LEVEL", "info")
}

func GetDBName() string {
	return getEnv("DB_NAME", "currency_api_db")
}

func GetDBHost() string {
	return getEnv("DB_HOST", "localhost")
}

func GetDBPort() string {
	return getEnv("DB_PORT", "5432")
}

func GetDBUser() string {
	return getEnv("DB_USER", "secretdbuser")
}

func GetDBPass() string {
	return getEnv("DB_PASS", "secretdbpass")
}

func GetXMLDataURLPath() string {
	return getEnv("XML_URL_PATH", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml1")
}

func GetXMLDataFilePath() string {
	return getEnv("XML_FILE_PATH", "./xmlfile/eurofxref-hist-90d.xml")
}

func GetServerHost() string {
	return getEnv("SERVER_HOST", "localhost")
}

func GetServerPort() string {
	return getEnv("SERVER_PORT", "9988")
}

func GetServerPublicKey() string {
	return getEnv("SERVER_PUBLIC_KEY", "./certs/cert.pem")
}

func GetServerPrivateKey() string {
	return getEnv("SERVER_PRIVATE_KEY", "./certs/key.pem")
}

func GetTokenSecret() string {
	return getEnv("TOKEN_SECRET", "notSoSecret")
}
