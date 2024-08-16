package main

import (
	"context"
	"math/rand"

	"github.com/dstotijn/go-notion"
)

const (
	tagDBNameProperty    = "Name"
	tagDBArticleProperty = "My Links"
	articleDBTagProperty = "Tags"
)

type notionClient struct {
	client *notion.Client
}

func newNotionClient(token string) *notionClient {
	return &notionClient{
		client: notion.NewClient(token),
	}
}

type tag struct {
	id    string
	name  string
	url   string
	emoji string
}

type article struct {
	id    string
	title string
	url   string
}

func (n *notionClient) pickTag(databaseId string, tagPrefix string) (tag, error) {
	titleFilter := &notion.TextPropertyFilter{IsNotEmpty: true}

	if tagPrefix != "" {
		titleFilter = &notion.TextPropertyFilter{StartsWith: tagPrefix}
	}

	tagFilter := notion.DatabaseQueryFilter{
		Property: tagDBNameProperty,
		DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
			Title: titleFilter,
		},
	}

	nonEmptyFilter := notion.DatabaseQueryFilter{
		Property: tagDBArticleProperty,
		DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
			Relation: &notion.RelationDatabaseQueryFilter{IsNotEmpty: true},
		},
	}

	// query Cathegories
	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: []notion.DatabaseQueryFilter{
				tagFilter,
				nonEmptyFilter,
			},
		}}

	tags, err := n.client.QueryDatabase(context.TODO(), databaseId, query)
	if err != nil {
		return tag{}, err
	}

	// choose a tag
	tagPage := tags.Results[rand.Intn(len(tags.Results))]
	props, _ := tagPage.Properties.(notion.DatabasePageProperties)

	tag := tag{
		id:    tagPage.ID,
		url:   tagPage.URL,
		name:  props["Name"].Title[0].PlainText,
		emoji: *tagPage.Icon.Emoji,
	}

	return tag, nil
}

func (n *notionClient) pickArticle(databaseId string, tagId string) (article, error) {
	tagFilter := notion.DatabaseQueryFilter{
		Property: articleDBTagProperty,
		DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
			Relation: &notion.RelationDatabaseQueryFilter{
				Contains: tagId,
			},
		},
	}

	query := &notion.DatabaseQuery{
		Filter: &tagFilter,
	}

	pages, err := n.client.QueryDatabase(context.TODO(), databaseId, query)
	if err != nil {
		return article{}, err
	}

	page := pages.Results[rand.Intn(len(pages.Results))]
	props, ok := page.Properties.(notion.DatabasePageProperties)
	if !ok {
		return article{}, err
	}

	article := article{
		id:    page.ID,
		url:   *props["URL"].URL,
		title: props["Name"].Title[0].PlainText,
	}

	return article, nil
}
