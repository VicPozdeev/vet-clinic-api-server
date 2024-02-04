package config

import (
	"embed"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm/logger"
	"os"
)

type Config struct {
	Database struct {
		Dialect   string `default:"sqlite3"`
		Host      string `default:"sqlite.db"`
		Port      string
		Dbname    string
		Username  string
		Password  string
		Migration bool `default:"false"`
	}
	Redis struct {
		Enabled            bool `default:"false"`
		ConnectionPoolSize int  `yaml:"connection_pool_size" default:"10"`
		Host               string
		Port               string
	}
	Extension struct {
		MasterGenerator bool `yaml:"master_generator" default:"false"`
		SecurityEnabled bool `yaml:"security_enabled" default:"false"`
		CorsEnabled     bool `yaml:"cors_enabled" default:"false"`
		CsrfEnabled     bool `yaml:"csrf_enabled" default:"false"`
	}
	StaticContents struct {
		Enabled bool `default:"false"`
	}
	Swagger struct {
		Enabled bool `default:"false"`
		Path    string
	}
	Logger struct {
		GormConfig logger.Config     `json:"gorm_config" yaml:"gorm_config"`
		ZapConfig  zap.Config        `json:"zap_config" yaml:"zap_config"`
		LogRotate  lumberjack.Logger `json:"log_rotate" yaml:"log_rotate"`
	}
}

func LoadConfig(configFile embed.FS) *Config {
	var env *string
	if value := os.Getenv("VET_CLINIC_ENV"); value != "" {
		env = &value
	} else {
		env = flag.String("env", "develop", "To switch configurations.")
		flag.Parse()
	}

	file, err := configFile.ReadFile(fmt.Sprintf(AppConfigPath, *env))
	if err != nil {
		fmt.Printf("Failed to read %s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		fmt.Printf("Failed to read %s.yml: %s", *env, err)
		os.Exit(ErrExitStatus)
	}

	return config
}
