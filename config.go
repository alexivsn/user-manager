package main

import (
	"errors"

	"git.containerum.net/ch/user-manager/clients"
	"git.containerum.net/ch/user-manager/models"
	"git.containerum.net/ch/user-manager/models/postgres"
	"git.containerum.net/ch/user-manager/server"
	"git.containerum.net/ch/user-manager/server/impl"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupLogger() error {
	switch gin.Mode() {
	case gin.TestMode, gin.DebugMode:
	case gin.ReleaseMode:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	viper.SetDefault("log_level", logrus.InfoLevel)
	level := logrus.Level(viper.GetInt("log_level"))
	if level > logrus.DebugLevel || level < logrus.PanicLevel {
		return errors.New("invalid log level")
	}
	logrus.SetLevel(level)
	return nil
}

func getListenAddr() string {
	viper.SetDefault("listen_addr", ":8111")
	return viper.GetString("listen_addr")
}

func getDB() (models.DB, error) {
	viper.SetDefault("db", "postgres")
	switch viper.GetString("db") {
	case "postgres":
		return postgres.DBConnect(viper.GetString("pg_url"))
	default:
		return nil, errors.New("invalid db")
	}
}

func getMailClient() (clients.MailClient, error) {
	viper.SetDefault("mail", "http")
	switch viper.GetString("mail") {
	case "http":
		return clients.NewHTTPMailClient(viper.GetString("mail_url")), nil
	default:
		return nil, errors.New("invalid mail client")
	}
}

func getReCaptchaClient() (clients.ReCaptchaClient, error) {
	viper.SetDefault("recaptcha", "http")
	switch viper.GetString("recaptcha") {
	case "http":
		return clients.NewHTTPReCaptchaClient(viper.GetString("recaptcha_key")), nil
	case "dummy":
		return clients.NewDummyReCaptchaClient(), nil
	default:
		return nil, errors.New("invalid reCaptcha client")
	}
}

func oauthClientsSetup() error {
	viper.SetDefault("oauth_clients", "http")
	switch viper.GetString("oauth_clients") {
	case "http":
		clients.RegisterOAuthClient(clients.NewGithubOAuthClient(viper.GetString("github_app_id"), viper.GetString("github_secret")))
		clients.RegisterOAuthClient(clients.NewGoogleOAuthClient(viper.GetString("google_app_id"), viper.GetString("google_secret")))
		clients.RegisterOAuthClient(clients.NewFacebookOAuthClient(viper.GetString("facebook_app_id"), viper.GetString("facebook_secret")))
	default:
		return errors.New("invalid oauth clients kind")
	}
	return nil
}

func getWebAPIClient() (clients.WebAPIClient, error) {
	viper.SetDefault("web_api", "http")
	switch viper.GetString("web_api") {
	case "http":
		return clients.NewHTTPWebAPIClient(viper.GetString("web_api_url")), nil
	default:
		return nil, errors.New("invalid web_api client")
	}
}

func getAuthClient() (clients.AuthClientCloser, error) {
	viper.SetDefault("auth", "grpc")
	switch viper.GetString("auth") {
	case "grpc":
		return clients.NewGRPCAuthClient(viper.GetString("auth_grpc_addr"))
	default:
		return nil, errors.New("invalid auth client")
	}
}

func getUserManager(services server.Services) (server.UserManager, error) {
	viper.SetDefault("user_manager", "impl")
	switch viper.Get("user_manager") {
	case "impl":
		return impl.NewUserManagerImpl(services), nil
	default:
		return nil, errors.New("invalid user manager impl")
	}
}

func getResourceServiceClient() (clients.ResourceServiceClient, error) {
	viper.SetDefault("resource_service", "http")
	switch viper.GetString("resource_service") {
	case "http":
		return clients.NewHTTPResourceServiceClient(viper.GetString("resource_service_url")), nil
	default:
		return nil, errors.New("invalid resource service client")
	}
}
