package config

type Config struct {
	ScreenWidth  int
	ScreenHeight int
	WorldWidth   int
	WorldHeigth  int
}

var C *Config

const (
	MapWidth  int = 80
	MapHeigth int = 45
	BlockSize int = 32
)

func init() {
	C = &Config{
		ScreenWidth:  640,
		ScreenHeight: 360,
		WorldWidth:   MapWidth * BlockSize,
		WorldHeigth:  MapHeigth * BlockSize,
	}
}
