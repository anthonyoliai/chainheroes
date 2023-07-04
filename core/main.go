package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anthonyoliai/chainheroes/character"
	"github.com/anthonyoliai/chainheroes/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

var hero = character.New("Adventurer X")
var train = 0

func training() {
	for {
		expedition := game.New("Tutorial", time.Second*5, 100)
		go hero.Train(expedition)
		time.Sleep(time.Second * 10)
	}
}

// Updates the logic state of the game
func (g *Game) Update() error {
	return nil
}

// Renders information e.g. derived from the logical state onto the screen
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hero name: %v", hero.Name()), 2, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hero status: %v", hero.CurrentStatus()), 2, 40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hero level: %v ", hero.Level()), 2, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hero experience: %.1f", hero.Experience()), 2, 70)

	if hero.CurrentStatus() == "Training" {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Currently on expedition: %v", hero.Expedition()), 2, 80)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Chain Heroes")
	go training()

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
