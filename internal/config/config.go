package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/kafka/consumer"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/kafka/producer"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	DevMode  = "DEVELOPMENT"
	ProdMode = "PRODUCTION"
)

const (
	baseCfgFilename = "base.yaml"
	envFilename     = ".env"
)

var cfgFileMapper = map[string]string{
	DevMode:  "dev.yaml",
	ProdMode: "prod.yaml",
}

type Config struct {
	AppMode        string `env:"APP_MODE" default:"DEVELOPMENT"`
	Server         Server `yaml:"server"`
	PostgresConfig PostgresConfig
	RedisConfig    RedisConfig
	KafkaConfig    KafkaConfig
	MinioConfig    MinioConfig
	TokenConfig    TokenConfig `yaml:"token"`
}

type Server struct {
	Host string `yaml:"host" default:"localhost"`
	Port int    `yaml:"port" default:"8080"`
}

type PostgresConfig struct {
	Host string `env:"POSTGRES_HOST"`
	Port string `env:"POSTGRES_PORT"`
	User string `env:"POSTGRES_USER"`
	Pass string `env:"POSTGRES_PASS"`
	DB   string `env:"POSTGRES_DB"`
	SSL  string `env:"POSTGRES_SSL"`
}

type RedisConfig struct {
	Pass string `env:"REDIS_PASS" default:""`
	Host string `env:"REDIS_HOST" default:"localhost"`
	Port int    `env:"REDIS_PORT" default:"6379"`
}

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Topics map[string]string `yaml:"topics"`

	Consumers map[string]KafkaConsumer `yaml:"consumers"`
	Producers map[string]KafkaProducer `yaml:"producers"`
}

type KafkaConsumer struct {
	Topic string `yaml:"topic"`
	Group string `yaml:"group"`
}

type KafkaProducer struct {
	Topic string `yaml:"topic"`
}

type MinioConfig struct {
	Endpoint string `env:"MINIO_ENDPOINT"` // "localhost:9000" or "minio:9000"
	User     string `env:"MINIO_ROOT_USER"`
	Password string `env:"MINIO_ROOT_PASSWORD"`

	UseSSL bool `env:"MINIO_USE_SSL" default:"false"`

	Bucket               string `env:"MINIO_BUCKET" default:"chat-x"`
	PresignExpirySeconds int    `env:"MINIO_PRESIGN_EXPIRY" default:"3600"`
}

type TokenConfig struct {
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
}

func Load() (*Config, error) {
	cfg := new(Config)

	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %w", err)
	}

	configDir := fmt.Sprintf("%s/configs/", workDir)

	err = cleanenv.ReadConfig(workDir+"/"+envFilename, cfg)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read from %s: %w", envFilename, err)
	}

	err = cleanenv.ReadConfig(configDir+baseCfgFilename, cfg)
	if err != nil && !isEOFerr(err) {
		return nil, fmt.Errorf("failed to read from %s: %w", baseCfgFilename, err)
	}

	modeFilename, ok := cfgFileMapper[cfg.AppMode]
	if ok {
		err = cleanenv.ReadConfig(configDir+modeFilename, cfg)
		if err != nil && !isEOFerr(err) {
			return nil, fmt.Errorf("failed to read from %s: %w", modeFilename, err)
		}
	}

	return cfg, nil
}

func isEOFerr(err error) bool {
	return strings.HasSuffix(err.Error(), io.EOF.Error())
}

func BuildConsumers(k KafkaConfig) (map[string]*consumer.Consumer, error) {
	out := make(map[string]*consumer.Consumer)
	broker := k.BrokerAddr()

	for name, spec := range k.Consumers {
		topic, err := k.ResolveTopic(spec.Topic)
		if err != nil {
			return nil, fmt.Errorf("consumer %s: %w", name, err)
		}
		if spec.Group == "" {
			return nil, fmt.Errorf("consumer %s: group is empty", name)
		}
		out[name] = consumer.NewConsumer(broker, topic, spec.Group)
	}

	return out, nil
}

func BuildProducers(k KafkaConfig) (map[string]*producer.Producer, error) {
	out := make(map[string]*producer.Producer)
	broker := k.BrokerAddr()

	for name, spec := range k.Producers {
		topic, err := k.ResolveTopic(spec.Topic)
		if err != nil {
			return nil, fmt.Errorf("producer %s: %w", name, err)
		}
		out[name] = producer.NewProducer(broker, topic)
	}

	return out, nil
}

func (k KafkaConfig) ResolveTopic(nameOrLiteral string) (string, error) {
	if nameOrLiteral == "" {
		return "", fmt.Errorf("kafka: topic is empty")
	}
	if k.Topics != nil {
		if t, ok := k.Topics[nameOrLiteral]; ok && t != "" {
			return t, nil
		}
	}
	return nameOrLiteral, nil
}

func (k KafkaConfig) BrokerAddr() string {
	return fmt.Sprintf("%s:%d", k.Host, k.Port)
}
