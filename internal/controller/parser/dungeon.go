package parser

import (
	"dungeon-challenge/internal/domain"
	"dungeon-challenge/internal/dto"
	"encoding/json"
	"log"
	"os"
)

type DungeonParser interface {
	ParseDungeon() (domain.Dungeon, error)
}

type DungeonParserStruct struct {
	dungeonFilename string
}

func NewDungeonParser(filename string) DungeonParser {
	return &DungeonParserStruct{
		dungeonFilename: filename,
	}
}

func (dp *DungeonParserStruct) ParseDungeon() (domain.Dungeon, error) {
	file, err := os.ReadFile(dp.dungeonFilename)
	if err != nil {
		log.Printf("error of opening file: %v", err)
		return domain.Dungeon{}, err
	}

	var dungeon dto.Dungeon

	err = json.Unmarshal(file, &dungeon)
	if err != nil {
		log.Printf("error of decoding file: %v", err)
		return domain.Dungeon{}, err
	}
	return dungeon.ToDomain(), nil
}
