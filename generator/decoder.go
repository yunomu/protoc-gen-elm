package generator

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func elmDecoder(field *protogen.Field) (string, error) {
	var d string
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		d = "Decode.string"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		d = "Decode.int"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		d = "Decode.string"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		d = "Decode.float"
	case protoreflect.BoolKind:
		d = "Decode.bool"
	case protoreflect.BytesKind:
		d = "decodeBytes"
	case protoreflect.MessageKind:
		messageName := field.Message.GoIdent.GoName
		d = strings.ToLower(messageName[:1]) + messageName[1:] + "Decoder"
	default:
		return "", fmt.Errorf("unknown decoder for field %s: %s", field.Desc.Name(), field.Desc.Kind().String())
	}

	if field.Desc.IsList() {
		return "(Decode.list " + d + ")", nil
	}

	return d, nil
}

func (g *Generator) genDecoder(f *GeneratedFile, msg *protogen.Message) {
	messageName := msg.GoIdent.GoName
	decoderName := strings.ToLower(messageName[:1]) + messageName[1:] + "Decoder"

	f.P(decoderName, " : Decode.Decoder ", messageName)
	f.P(decoderName, " =")

	fieldsNum := len(msg.Fields)
	switch {
	case fieldsNum > 8:
		g.gen.Error(fmt.Errorf("Maximum number of fields supported is 8"))
		return
	case fieldsNum == 0:
		panic("not reached: number of fields is not zero")
	case fieldsNum == 1:
		f.P("    Decode.map ", messageName)
	default:
		f.P("    Decode.map", fieldsNum, " ", messageName)
	}

	for _, field := range msg.Fields {
		jsonName := field.Desc.JSONName()
		fieldDecoder, err := elmDecoder(field)
		if err != nil {
			g.gen.Error(err)
			return
		}
		labelDecoder := "Decode.field"
		if field.Desc.HasOptionalKeyword() {
			labelDecoder = "Decode.maybe <| Decode.field"
		}
		f.P("        (", labelDecoder, " \"", jsonName, "\" ", fieldDecoder, ")")
	}

	f.Exposing(decoderName)
}
