package cmd

import (
	"fmt"
	"log"

	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type Struct struct {
	Name   string
	Fields []Item
}

type Item struct {
	Name   string
	Type   string
	Fields []Item
}

func Generate() {
	schema := gqlparser.MustLoadSchema(&ast.Source{
		Input: `
				schema { query: Query }
				type Query {
					me: User!
					user(id: ID): User!
				}
				type User {
					id: ID!
					name: String
				}
			`,
	})
	query := gqlparser.MustLoadQuery(schema, `
		query Hi($id: ID) {
			user(id: $id) {
				id
				name 
			}
			me {
				id
			}
		}
	`)

	log.Printf("schema: %+v", schema.Types)

	var types []string
	var structs []Struct

	for _, t := range schema.Types {
		if t.Fields == nil {
			var typeName string
			switch t.Name {
			case "String":
				typeName = "string"
			case "Boolean":
				typeName = "bool"
			case "Float":
				typeName = "float"
			case "Int":
				typeName = "int"
			default:
				typeName = "string"
			}
			types = append(types, fmt.Sprintf("type %s %s\n", t.Name, typeName))
		} else {
			log.Printf("ext: %+v", t.Name)
			s := Struct{
				Name: t.Name,
			}

			for _, sel := range t.Fields {
				s.Fields = append(s.Fields, Item{
					Name: sel.Name,
					Type: sel.Type.Name(),
				})
			}

			structs = append(structs, s)
		}
	}

	for _, op := range query.Operations {
		s := Struct{
			Name: op.Name,
		}

		for _, sel := range op.SelectionSet {
			switch t := sel.(type) {
			case *ast.Field:
				s.Fields = append(s.Fields, Item{
					Name: t.Name,
					Type: t.Definition.Type.Name(),
				})
			default:
				panic("no matching type for selection set field")
			}
		}

		structs = append(structs, s)
	}

	log.Printf("types: %+v", types)
	log.Printf("structs: %+v", structs)

	str := "package gqlclient\n\n"
	for _, t := range types {
		str += fmt.Sprintf("%s", t)
	}
	for _, s := range structs {
		str += fmt.Sprintf("type %s struct {\n", s.Name)
		for _, f := range s.Fields {
			str += fmt.Sprintf("  %s %s\n", gocase.To(strcase.ToCamel(f.Name)), f.Type)
		}
		str += fmt.Sprintf("}\n\n")
	}
	log.Printf("file: %s", str)
}
