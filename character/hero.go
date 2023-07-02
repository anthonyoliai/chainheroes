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
	return hero{name: name}
}

func (h *hero) Train(expedition game.Expedition) {
	h.status = Training
	fmt.Printf("Hero is training until %v", time.Now().Add(expedition.Duration()))
	<-time.After(expedition.Duration())
	h.status = Idle
	h.experience += expedition.Experience()
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
