package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("UserName").Unique(),

		field.Bytes("PasswordHash"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("UserName").Unique(),
	}
}
