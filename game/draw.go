package game

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen,
		"debug: "+g.Debug.Foo+
			"\nsuccessful server roundtrips: "+strconv.Itoa(g.Debug.Roundtrips)+
			"\ng.Debug.Frame: "+strconv.Itoa(g.Debug.Frame)+
			"\nconnected players: "+strconv.Itoa(g.Debug.ConnectedPlayers)+
			"\n\ntime: "+time.Now().Format(time.RFC3339))
}
