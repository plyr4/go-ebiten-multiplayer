package game

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

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
	// screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	renderTarget := ebiten.NewImage(constants.WINDOW_WIDTH, constants.WINDOW_HEIGHT)

	// server status
	serverStatusSize := 8
	var serverStatusImg = ebiten.NewImage(serverStatusSize, serverStatusSize)
	serverStatus := "down"
	serverStatusColor := color.RGBA{0xff, 0, 0, 0xff}
	if g.wsClient != nil {
		if g.wsClient.IsConnected() {
			serverStatus = "up"
			serverStatusColor = color.RGBA{0, 0xff, 0, 0xff}
		}
	}
	serverStatusImg.Fill(serverStatusColor)

	renderTarget.DrawImage(serverStatusImg, nil)

	ebitenutil.DebugPrint(renderTarget,
		fmt.Sprintf("fps: %v", math.Round(ebiten.ActualFPS()))+
			"\n input: "+fmt.Sprintf("%v", *g.Input)+
			"\nserver_status: "+serverStatus+
			"\nserver_roundtrips: "+strconv.Itoa(g.Roundtrips)+
			"\n"+fmt.Sprintf(" connected_players: %v", g.ConnectedPlayers)+
			"\n current_frame: "+strconv.Itoa(g.Frame)+
			"\n num_entities: "+fmt.Sprintf("%v", len(g.entities)))

	// entities
	for _, e := range g.entities {
		e.Draw(renderTarget)
	}

	// draw render target to screen
	// apply world scale, zoom etc
	g.DrawImageOptions.GeoM.Scale(constants.WORLD_SCALE, constants.WORLD_SCALE)
	screen.DrawImage(renderTarget, g.DrawImageOptions)
	g.DrawImageOptions.GeoM.Reset()

}
