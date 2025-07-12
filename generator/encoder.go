package generator

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func elmEncoder(field *protogen.Field, fieldAccessor string) (string, error) {
	var e string
	switch field.Desc.Kind() {
	case protoreflect.StringKind:
		e = "Encode.string"
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		e = "Encode.int"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		e = "Encode.float"
	case protoreflect.BoolKind:
		e = "Encode.bool"
	case protoreflect.MessageKind:
		messageName := field.Message.GoIdent.GoName
		e = strings.ToLower(messageName[:1]) + messageName[1:] + "Encoder"
	default:
		return "", fmt.Errorf("unknown encoder for field %s: %s", field.Desc.Name(), field.Desc.Kind().String())
	}

	if field.Desc.IsList() {
		// The encoder function for the inner type doesn't take an argument, so we need to handle it slightly differently.
		encoderName := e
		if !strings.HasSuffix(e, "Encoder") {
			// For primitive types, the function name is just the type (e.g., Encode.string)
			// For message types, it's the function name we constructed (e.g., userEncoder)
			encoderName = e
		}
		return "(Encode.list " + encoderName + " " + fieldAccessor + ")", nil
	}
	return e + " " + fieldAccessor, nil
}

func (g *Generator) genEncoder(f *protogen.GeneratedFile, msg *protogen.Message) {
	messageName := msg.GoIdent.GoName
	encoderName := strings.ToLower(messageName[:1]) + messageName[1:] + "Encoder"
	instanceName := strings.ToLower(messageName[:1]) + messageName[1:]

	f.P(encoderName, " : ", messageName, " -> Encode.Value")
	f.P(encoderName, " ", instanceName, " =")
	f.P("    Encode.object")

	separator := "["
	for _, field := range msg.Fields {
		fieldName := camelCase(string(field.Desc.Name()))
		jsonName := field.Desc.JSONName()
		fieldEncoder, err := elmEncoder(field, instanceName+"."+fieldName)
		if err != nil {
			g.gen.Error(err)
			return
		}

		f.P("        ", separator, " (\"", jsonName, "\", ", fieldEncoder, ")")
		separator = ","
	}

	f.P("        ]")
	f.P("")
}
