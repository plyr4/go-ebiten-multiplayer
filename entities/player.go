package entities

import (
	"bytes"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/plyr4/go-ebiten-multiplayer/constants"
	"github.com/plyr4/go-ebiten-multiplayer/internal"
	"github.com/plyr4/go-ebiten-multiplayer/resources/images"
	"github.com/plyr4/go-ebiten-multiplayer/sprite"
	"github.com/plyr4/go-ebiten-multiplayer/ws"
)

type Player struct {
	Entity
	*ebiten.Image
	Movement
}

const (
	LEFT  = 0
	RIGHT = 1
)

type Movement struct {
	speed float64
	Dir   int
	DX    float64
	DY    float64
}

func NewPlayer(g *internal.Game) (*Player, error) {
	p := new(Player)

	p.Game = g

	// starting position
	p.X = constants.SCREEN_WIDTH / 2
	p.Y = constants.SCREEN_HEIGHT / 2

	// movement
	p.speed = 2

	// sprite
	img, _, err := image.Decode(bytes.NewReader(images.Runner))
	if err != nil {
		return nil, err
	}
	p.Image = ebiten.NewImageFromImage(img)
	// p.Hue = 2.0 * math.Pi * (rand.Float64() * 50)

	return p, nil
}

func (p *Player) Draw(renderTarget *ebiten.Image) error {
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
	opts.GeoM.Translate(p.X, p.Y)

	// apply the player hue
	cm := colorm.ColorM{}
	cm.Reset()
	// cm.RotateHue(p.Hue)

	// draw the sprite
	colorm.DrawImage(renderTarget,
		sprite.Sprite(p.Image, frameOpts),
		cm, opts)

	return nil
}

func (p *Player) Update() error {
	var dx, dy float64

	if p.Game.Input.Up {
		dy -= 1
	}

	if p.Game.Input.Down {
		dy += 1
	}

	if p.Game.Input.Left {
		dx -= 1
		p.Dir = LEFT
	}

	if p.Game.Input.Right {
		dx += 1
		p.Dir = RIGHT
	}

	// normalize
	if dx != 0 && dy != 0 {
		length := math.Sqrt(dx*dx + dy*dy)
		dx /= length
		dy /= length
	}

	p.DX = dx
	p.DY = dy

	p.X += p.DX * p.speed
	p.Y += p.DY * p.speed

	return nil
}

func (p *Player) ToMultiplayerData() ws.PlayerData {
	return ws.PlayerData{
		UUID: p.UUID,
		X:    p.X,
		Y:    p.Y,
		DX:   p.DX,
		DY:   p.DY,
		Dir:  p.Dir,
	}
}
