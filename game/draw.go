package game

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/plyr4/go-ebiten-multiplayer/constants"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return constants.WINDOW_WIDTH, constants.WINDOW_HEIGHT
}

func (g *Game) Draw(screen *ebiten.Image) {
	// main background
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	// debug
	debugObjectSize := constants.WINDOW_WIDTH / 15
	// debugRed := color.RGBA{0xff, 0, 0, 0xff}
	opts := &ebiten.DrawImageOptions{}

	// box
	var pointerImage = ebiten.NewImage(debugObjectSize, debugObjectSize)
	pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	screen.DrawImage(pointerImage, opts)

	ebitenutil.DebugPrint(screen,
		"debug: "+g.Debug.Foo+
			"\nsuccessful server roundtrips: "+strconv.Itoa(g.Debug.Roundtrips)+
			"\ng.Debug.Frame: "+strconv.Itoa(g.Debug.Frame)+
			"\n"+fmt.Sprintf("connected players: %v", g.Debug.ConnectedPlayers)+
			"\ninput: "+fmt.Sprintf("%v", *g.Input)+
			"\n\ntime: "+time.Now().Format(time.RFC3339))
	for _, e := range g.entities {
		e.Draw(screen)
	}

	// todo: just convert the other players to entities... avoid this mess
	// todo: when we receive a server update, we make sure those entities exist in our list
	// entities could include other things like enemies, items, etc.
	// entities probably need sub entities or drawing order or something like that
	// the player should be drawn last, on top of everything else
	// maybe not all entities should exist in the same list
	// it would make updating the list of other players easier
	// assign a new player a color and maintain that color, draw it here

	for _, p := range g.Debug.ConnectedPlayers {
		if p.UUID == g.uuid {
			continue
		}
		pp := g.player
		x, y := pp.X, pp.Y
		pp.X = p.X
		pp.Y = p.Y
		pp.Draw(screen)
		pp.X = x
		pp.Y = y
	}
}
