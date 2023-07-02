package game

import "time"

type Expedition struct {
	name       string
	duration   time.Duration
	experience float64
}

func New(name string, duration time.Duration, experience float64) Expedition {
	return Expedition{name, duration, experience}
}

func (e *Expedition) Duration() time.Duration {
	return e.duration
}

func (e *Expedition) Experience() float64 {
	return e.experience
}

func (e *Expedition) Name() string {
	return e.name
}
