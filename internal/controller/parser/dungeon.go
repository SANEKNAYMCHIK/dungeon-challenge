package parser

import (
	"dungeon-challenge/internal/domain"
	"dungeon-challenge/internal/dto"
	"encoding/json"
	"os"
)

type DungeonParser struct {
	dungeonFilename string
}

func NewDungeonParser(filename string) DungeonParser {
	return DungeonParser{
		dungeonFilename: filename,
	}
}

func (dp DungeonParser) ParseDungeon() (domain.Dungeon, error) {
	file, err := os.ReadFile(dp.dungeonFilename)
	if err != nil {
		return domain.Dungeon{}, err
	}

	var dungeon dto.Dungeon

	err = json.Unmarshal(file, &dungeon)
	if err != nil {
		return domain.Dungeon{}, err
	}
	return dungeon.ToDomain(), nil
}
