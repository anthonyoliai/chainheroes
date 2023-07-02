package character

import (
	"fmt"
	"time"

	"github.com/anthonyoliai/chainheroes/game"
)

type hero struct {
	tokenID    uint64
	level      uint64
	hunger     uint64
	health     uint64
	experience float64
	name       string
	mood       Mood
	status     Status
	expedition *game.Expedition
}

type Mood uint64

type Status uint64

const (
	Happy Mood = iota
	Neutral
	Upset
	Adventurous
	Hungry
	Sleepy
)

const (
	Idle Status = iota
	Training
	Sleeping
)

func New(name string) hero {
	return hero{name: name, level: 1}
}

func (h *hero) Train(expedition game.Expedition) {
	h.status = Training
	h.expedition = &expedition
	<-time.After(expedition.Duration())
	h.handleLevelUp(expedition.Experience())
	h.status = Idle
	h.expedition = nil
}

func (h *hero) handleLevelUp(experience float64) {
	remainingExp := experience

	for {
		requiredExp := float64(h.level*h.level) - h.experience

		if requiredExp <= remainingExp {
			h.level++
			h.experience = 0
			remainingExp = remainingExp - requiredExp
			fmt.Println(remainingExp)
			continue
		}
		h.experience += remainingExp
		return
	}

}

func (h *hero) CurrentStatus() string {
	switch h.status {
	case Idle:
		return "Idle"
	case Training:
		return "Training"
	case Sleeping:
		return "Sleeping"
	default:
		return ""
	}
}

func (h *hero) Name() string {
	return h.name
}

func (h *hero) Level() uint64 {
	return h.level
}

func (h *hero) Experience() float64 {
	return h.experience
}

func (h *hero) Expedition() string {
	if h.expedition == nil {
		return ""
	}
	return h.expedition.Name()
}
