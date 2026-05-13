package output

import "dungeon-challenge/internal/domain"

var templates = map[domain.EventType]string{
	domain.EventRegistered:       "%s Player [%v] registered\n",
	domain.EventInDungeon:        "%s Player [%v] entered the dungeon\n",
	domain.EventKilledMonster:    "%s Player [%v] killed the monster\n",
	domain.EventNextFloor:        "%s Player [%v] went to the next floor\n",
	domain.EventPreviousFloor:    "%s Player [%v] went to the previous floor\n",
	domain.EventEnteredBossFloor: "%s Player [%v] entered the boss's floor\n",
	domain.EventKilledBoss:       "%s Player [%v] killed the boss\n",
	domain.EventLeftDungeon:      "%s Player [%v] left the dungeon\n",
	domain.EventFailed:           "%s Player [%v] cannot continue due to [%v]\n",
	domain.EventGetHealth:        "%s Player [%v] has restored [%v] of health\n",
	domain.EventGetDamage:        "%s Player [%v] recieved [%v] of damage\n",
	domain.EventDisqualified:     "%s Player [%v] disqualified\n",
	domain.EventDead:             "%s Player [%v] is dead\n",
	domain.EventImpossibleMove:   "%s Player [%v] makes imposible move [%v]\n",
}
