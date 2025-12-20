package config

import (
	"log"

	"github.com/spf13/viper"
)

// 설정 값을 담을 구조체 정의 (YAML 구조와 일치해야 함)
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		File string `mapstructure:"file"`
	} `mapstructure:"database"`

	Log struct {
		Level      string `mapstructure:"level"`
		Path       string `mapstructure:"path"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
	} `mapstructure:"log"`
}

// 전역 설정 변수
var AppConfig *Config

// 설정 로드 함수
func LoadConfig() {
	viper.AddConfigPath(".")      // 현재 디렉토리에서 찾기
	viper.SetConfigName("config") // 파일 이름 (config.yaml)
	viper.SetConfigType("yaml")   // 파일 형식

	viper.AutomaticEnv() // 환경 변수도 읽을 수 있게 설정

	// 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 읽은 값을 구조체에 매핑 (Unmarshalling)
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
