package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Secret struct {
	ent.Schema
}

func (Secret) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").Unique(),

		field.Bytes("value"),
	}
}
