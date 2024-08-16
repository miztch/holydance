package main

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

func sendDiscordMessage(article article, tag tag, webhookURL string) error {
	header := fmt.Sprintf(`[%s%s](%s)`, tag.emoji, tag.name, tag.url)
	msg := fmt.Sprintf("%s\n**%s**\n%s", header, article.title, article.url)

	// send to discord
	client, err := webhook.NewWithURL(webhookURL)
	if err != nil {
		return fmt.Errorf("failed to create webhook client: %v", err)
	}

	_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetContent(msg).Build())
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
