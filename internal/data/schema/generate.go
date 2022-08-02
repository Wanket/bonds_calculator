//go:build pregenerate

package schema

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate . --target=../entgenerated
