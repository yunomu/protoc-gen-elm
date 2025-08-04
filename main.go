package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yunomu/protoc-gen-elm/generator"
)

func main() {
	var golden bool
	flag.BoolVar(&golden, "golden", false, "golden")

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		g := generator.New(gen, generator.WithGolden(golden))
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			g.GenerateFile(f)
		}
		return nil
	})
}
