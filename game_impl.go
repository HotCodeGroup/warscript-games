package main

import (
	"github.com/google/uuid"
)

func getGameBySlugImpl(slug string) (*GameFull, error) {
	game, err := Games.GetGameBySlug(slug)

	if err != nil {
		return nil, err
	}

	return &GameFull{
		Game: Game{
			Slug:           game.Slug.String,
			Title:          game.Title.String,
			BackgroundUUID: uuid.UUID(game.BackgroundUUID.Bytes).String(), // точно 16 байт
		},
		Description: game.Description.String,
		Rules:       game.Rules.String,
		CodeExample: game.CodeExample.String,
		BotCode:     game.BotCode.String,
		LogoUUID:    uuid.UUID(game.LogoUUID.Bytes).String(), // точно 16 байт
	}, nil
}
