package main

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/yunomu/protoc-gen-elm/generator"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		g := generator.New(gen)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			g.GenerateFile(f)
		}
		return nil
	})
}
