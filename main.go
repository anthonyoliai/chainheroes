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

var hero = character.New("Tony soyboy")
var train = 0

// Updates the logic state of the game
func (g *Game) Update() error {
	if train == 500 {
		expedition := game.New("Tutorial", time.Second*10, 10)
		go hero.Train(expedition)
	}
	train++
	return nil
}

// Renders information e.g. derived from the logical state onto the screen
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("You hero %v is currently %v with current level: %v with %.1f experience.", hero.Name(), hero.CurrentStatus(), hero.Level(), hero.Experience()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
