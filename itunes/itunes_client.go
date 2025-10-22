package client

import (
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Client struct {
	HTTP   *resty.Client
	Logger *zap.Logger
	Config Config
}

type Config struct {
	Timeout    time.Duration
	RetryCount int
	RetryDelay time.Duration
	UserAgent  string
	Debug      bool
}

func NewClient(config Config) *Client {
	var logger *zap.Logger
	var err error

	if config.Debug {
		// Development config with colors and console encoder
		developmentConfig := zap.NewDevelopmentConfig()
		developmentConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		developmentConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		developmentConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		logger, err = developmentConfig.Build()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		logger = zap.NewNop()
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}
	if config.UserAgent == "" {
		config.UserAgent = "go-api-sdk-apple/1.0.0"
	}

	httpClient := resty.New().
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryDelay).
		SetHeader("User-Agent", config.UserAgent)

	if config.Debug {
		httpClient.SetDebug(true)
	}

	return &Client{
		HTTP:   httpClient,
		Logger: logger,
		Config: config,
	}
}

func NewDefaultClient() *Client {
	return NewClient(Config{})
}

func (c *Client) Close() {
	if c.Logger != nil {
		c.Logger.Sync()
	}
}
