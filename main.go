package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func handler(ctx context.Context) {
	notionConfig := getNotionConfig()

	n := newNotionClient(notionConfig.token)
	tag, err := n.pickTag(notionConfig.tagDatabaseId, *notionConfig.tagNamePrefix)
	if err != nil {
		slog.Error("failed to pick a tag from database", "Error", err.Error())
		return
	}

	article, err := n.pickArticle(notionConfig.articleDatabaseId, tag.id)
	if err != nil {
		slog.Error("failed to pick an article from database", "Error", err.Error())
		return
	}

	err = sendDiscordMessage(article, tag, getWebhookConfig().URL)
	if err != nil {
		slog.Error("failed to send article to discord", "Error", err.Error())
		return
	}
}

func main() {
	if isRunningOnLambda() {
		lambda.Start(handler)
	} else {
		err := LoadEnvFromConfig()
		if err != nil {
			slog.Error("failed to load Environment Variables.", "Error", err.Error())
			return
		}
		handler(context.Background())
	}
}
