package entities

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/plyr4/go-ebiten-multiplayer/internal"
	"github.com/plyr4/go-ebiten-multiplayer/resources/images"
	"github.com/plyr4/go-ebiten-multiplayer/sprite"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
)

type NetworkPlayer struct {
	Entity
	*ebiten.Image
	*ws.PlayerData
}

func NewNetworkPlayer(g *internal.Game, pd *ws.PlayerData) (*NetworkPlayer, error) {
	p := new(NetworkPlayer)

	p.Game = g

	// this entity renders using network data
	p.PlayerData = pd

	// sprite
	img, _, err := image.Decode(bytes.NewReader(images.Runner))
	if err != nil {
		return nil, err
	}

	p.Image = ebiten.NewImageFromImage(img)

	return p, nil
}

func (p *NetworkPlayer) Draw(renderTarget *ebiten.Image) error {
	frame := 0

	// player is moving
	if p.DX != 0 || p.DY != 0 {
		frame = p.Game.Frame
	}

	// store this in the entity
	frameOpts := &sprite.FrameOpts{
		CurrentGameFrame: frame,
		FrameOX:          0,
		FrameOY:          32,
		FrameWidth:       32,
		FrameHeight:      32,
		FrameCount:       8,
	}

	opts := &colorm.DrawImageOptions{}

	// center the sprite
	opts.GeoM.Translate(
		-float64(frameOpts.FrameWidth)/2,
		-float64(frameOpts.FrameHeight)/2,
	)

	// flip the sprite facing left
	if p.Dir == LEFT {
		opts.GeoM.Scale(-1, 1)
	}

	// move the sprite to the player's position
	opts.GeoM.Translate(p.PlayerData.X, p.PlayerData.Y)

	// apply the player hue
	cm := colorm.ColorM{}
	cm.Reset()
	if !p.PlayerData.Connected {
		cm.Scale(0.5, 0.5, 0.5, 0.8)
	}

	// draw the sprite
	colorm.DrawImage(renderTarget,
		sprite.Sprite(p.Image, frameOpts),
		cm, opts)

	return nil
}

func (p *NetworkPlayer) Update() error {
	return nil
}
