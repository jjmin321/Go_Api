package naturallanguage

import (
	"context"
	"log"

	language "cloud.google.com/go/language/apiv1"
	"github.com/labstack/echo"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Text is parameter form struct
type Text struct {
	Text string `json:"text" form:"text" query:"text"`
}

// Analyze method return analyzed result of text
func Analyze(c echo.Context) error {
	ctx := context.Background()

	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	u := new(Text)
	if err := c.Bind(u); err != nil {
		return err
	}

	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: u.Text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	entities, err := client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: u.Text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	return c.JSON(200, map[string]interface{}{
		"status":          200,
		"message":         "감정 분석 성공",
		"Text":            u.Text,
		"Sentiment score": sentiment.DocumentSentiment.Score,
		// "Entities":        entities.Entities,
		"Language": entities.Language,
	})
}
