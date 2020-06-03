package main

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type basic struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func loadedSprite(url string) *common.Texture {
	tex, err := common.LoadedSprite(url)
	if err != nil {
		handleErr(err)
	}

	return tex
}

func handleErr(err error) {
	if err != nil {
		log.Fatal("error: ", err)
	}
}
