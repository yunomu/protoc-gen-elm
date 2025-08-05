package generator

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func elmType(field *protogen.Field) (string, error) {
	var t string
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		t = "String"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		t = "Int"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		t = "String"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		t = "Float"
	case protoreflect.BoolKind:
		t = "Bool"
	case protoreflect.BytesKind:
		t = "Bytes.Bytes"
	case protoreflect.MessageKind:
		t = field.Message.GoIdent.GoName
	default:
		return "", fmt.Errorf("unknown type for field %s: %s", field.Desc.Name(), field.Desc.Kind().String())
	}

	if field.Desc.IsList() {
		return "List " + t, nil
	}

	if field.Desc.HasOptionalKeyword() {
		return "Maybe " + t, nil
	}

	return t, nil
}

func (g *Generator) genMessage(f *GeneratedFile, msg *protogen.Message) {
	messageName := msg.GoIdent.GoName
	f.P("type alias ", messageName, " =")

	separator := "{"
	for _, field := range msg.Fields {
		fieldName := camelCase(string(field.Desc.Name()))
		fieldType, err := elmType(field)
		if err != nil {
			g.gen.Error(err)
			return
		}

		if fieldType == "Bytes.Bytes" || fieldType == "(Maybe Bytes.Bytes)" {
			g.hasBytes = true
		}

		if field.Desc.HasOptionalKeyword() {
			g.hasOptional = true
		}

		f.P("    ", separator, " ", fieldName, " : ", fieldType)
		separator = ","
	}
	f.P("    }")

	f.Exposing(messageName)
}
