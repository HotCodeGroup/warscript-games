package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/HotCodeGroup/warscript-utils/utils"
	"github.com/pkg/errors"
)

func TestGetGameBySlug(t *testing.T) {
	m := &GamesManager{}

	Games = &gameTest{
		games: map[string]*GameModel{
			"pong": {
				ID:          1,
				Slug:        "pong",
				Title:       "Pong",
				Description: "Very cool game(net)",
				Rules:       "Do not cheat, please",
				CodeExample: "const a = 5;",
				BotCode:     "const a = 5;",
			},
		},
	}

	cases := []struct {
		slug          string
		expected      *models.InfoGame
		expectedError error
	}{
		{
			slug: "pong",
			expected: &models.InfoGame{
				Slug:        "pong",
				Title:       "Pong",
				Description: "Very cool game(net)",
				Rules:       "Do not cheat, please",
				CodeExample: "const a = 5;",
				BotCode:     "const a = 5;",
			},
		},
		{
			slug:          "ping-pong",
			expectedError: utils.ErrNotExists,
		},
	}

	for i, c := range cases {
		req := &models.GameSlug{Slug: c.slug}
		resp, err := m.GetGameBySlug(context.Background(), req)
		if errors.Cause(err) != c.expectedError {
			t.Errorf("[%d] GetGameBySlug got unexpected error: %v, expected: %v", i, err, c.expectedError)
		}
		if !reflect.DeepEqual(resp, c.expected) {
			t.Errorf("[%d] GetGameBySlug returns: %v, wanted: %v", i, resp, c.expected)
		}
	}
}
