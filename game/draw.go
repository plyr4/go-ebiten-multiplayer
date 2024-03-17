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
	// todo: obviously clean this up

	// main background
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	// server status
	serverStatusSize := 8
	var pointerImage = ebiten.NewImage(serverStatusSize, serverStatusSize)
	serverStatusColor := color.RGBA{0xff, 0, 0, 0xff}

	if g.wsClient != nil {
		if g.wsClient.IsConnected() {
			serverStatusColor = color.RGBA{0, 0xff, 0, 0xff}
		}
	}

	pointerImage.Fill(serverStatusColor)

	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(pointerImage, opts)

	ebitenutil.DebugPrint(screen,
		"debug: "+g.Debug.Foo+
			"\n"+fmt.Sprintf("fps: %v", ebiten.ActualFPS())+
			"\nsuccessful server roundtrips: "+strconv.Itoa(g.Debug.Roundtrips)+
			"\ng.Debug.Frame: "+strconv.Itoa(g.Debug.Frame)+
			"\n"+fmt.Sprintf("connected players: %v", g.Debug.ConnectedPlayers)+
			"\ninput: "+fmt.Sprintf("%v", *g.Input)+
			"\n\ntime: "+time.Now().Format(time.RFC3339))
	// todo: just convert the other players to entities... avoid this mess
	// todo: when we receive a server update, we make sure those entities exist in our list
	// entities could include other things like enemies, items, etc.
	// entities probably need sub entities or drawing order or something like that
	// the player should be drawn last, on top of everything else
	// maybe not all entities should exist in the same list
	// it would make updating the list of other players easier
	// assign a new player a color and maintain that color, draw it here

	// todo: fix: when these overlap they are not happy
	// its possible that the map is getting jumbled every time we draw
	for _, p := range g.Debug.ConnectedPlayers {
		if p.UUID == g.uuid || !p.Connected {
			continue
		}

		pp := g.player

		x, y, hue := pp.X, pp.Y, pp.Hue
		pp.X = p.X
		pp.Y = p.Y
		pp.Hue = p.Hue

		pp.Draw(screen)

		pp.X = x
		pp.Y = y
		pp.Hue = hue
	}

	for _, e := range g.entities {
		e.Draw(screen)
	}
}
