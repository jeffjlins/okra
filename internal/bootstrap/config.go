package bootstrap

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Firestore FirestoreConfig
}

type ServerConfig struct {
	Port string
}

type FirestoreConfig struct {
	ProjectID       string
	DatabaseID      string // Optional: defaults to "(default)" if not specified
	CredentialsFile string // Optional: path to service account JSON file. If empty, uses GOOGLE_APPLICATION_CREDENTIALS env var or default credentials
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.okra")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("firestore.project_id", "")
	viper.SetDefault("firestore.database_id", "(default)")
	viper.SetDefault("firestore.credentials_file", "")

	// Environment variables
	viper.SetEnvPrefix("OKRA")
	viper.AutomaticEnv()
	viper.BindEnv("server.port", "OKRA_SERVER_PORT")
	viper.BindEnv("firestore.project_id", "OKRA_FIRESTORE_PROJECT_ID")
	viper.BindEnv("firestore.database_id", "OKRA_FIRESTORE_DATABASE_ID")
	viper.BindEnv("firestore.credentials_file", "OKRA_FIRESTORE_CREDENTIALS_FILE")

	// Read config file (optional - will use defaults if not found)
	if err := viper.ReadInConfig(); err != nil {
		// Only return error if it's not a "file not found" error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Resolve relative paths relative to the config file location (if config file was found)
	// or relative to current working directory
	credentialsFile := viper.GetString("firestore.credentials_file")
	if credentialsFile != "" && !filepath.IsAbs(credentialsFile) {
		if configFile := viper.ConfigFileUsed(); configFile != "" {
			// Resolve relative to config file directory
			configDir := filepath.Dir(configFile)
			credentialsFile = filepath.Join(configDir, credentialsFile)
			// Clean the path (removes ./ and ../)
			credentialsFile = filepath.Clean(credentialsFile)
		}
		// If no config file was used, relative paths are resolved from current working directory
		// which is the default behavior, so no change needed
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
		},
		Firestore: FirestoreConfig{
			ProjectID:       viper.GetString("firestore.project_id"),
			DatabaseID:      viper.GetString("firestore.database_id"),
			CredentialsFile: credentialsFile,
		},
	}

	if config.Firestore.ProjectID == "" {
		return nil, fmt.Errorf("firestore.project_id is required (set via config file or OKRA_FIRESTORE_PROJECT_ID env var)")
	}

	return config, nil
}
