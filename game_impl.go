package main

func getGameBySlugImpl(slug string) (*GameFull, error) {
	game, err := Games.GetGameBySlug(slug)

	if err != nil {
		return nil, err
	}

	return &GameFull{
		Game: Game{
			Slug:           game.Slug,
			Title:          game.Title,
			BackgroundUUID: game.GetBackgroundUUID(),
		},
		Description: game.Description,
		Rules:       game.Rules,
		CodeExample: game.CodeExample,
		BotCode:     game.BotCode,
		LogoUUID:    game.GetLogoUUID(), // точно 16 байт
	}, nil
}
