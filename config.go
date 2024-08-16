package main

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvFromConfig loads environment variables from a .env file
func LoadEnvFromConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

type notionConfig struct {
	token             string
	tagDatabaseId     string
	tagNamePrefix     *string
	articleDatabaseId string
}

func getNotionConfig() notionConfig {
	tagNamePrefix := os.Getenv("NOTION_TAG_NAME_PREFIX")
	return notionConfig{
		token:             os.Getenv("NOTION_INTEGRATION_TOKEN"),
		articleDatabaseId: os.Getenv("NOTION_ARTICLE_DATABASE_ID"),
		tagDatabaseId:     os.Getenv("NOTION_TAG_DATABASE_ID"),
		tagNamePrefix:     &tagNamePrefix,
	}
}

type webhookConfig struct {
	URL string
}

func getWebhookConfig() webhookConfig {
	return webhookConfig{
		URL: os.Getenv("WEBHOOK_URL"),
	}
}
