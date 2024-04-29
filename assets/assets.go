package assets

import (
	"embed"
	"encoding/json"
)

func Load() error {

	sprites := &spriteConfig{}
	mustReadJSON("config/sprites.json", sprites)

	loadSprites(sprites)
	loadAnimations(sprites)

	return nil
}

//go:embed img/*.png config/*.json
var fs embed.FS

func mustRead(name string) []byte {
	b, err := fs.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return b
}

func mustReadJSON(name string, v interface{}) {
	b := mustRead(name)
	if err := json.Unmarshal(b, v); err != nil {
		panic(err)
	}
}
