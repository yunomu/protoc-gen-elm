package main

import (
	"github.com/yunomu/protoc-gen-elm/generate"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generate.GenerateFile(gen, f)
		}
		return nil
	})
}
