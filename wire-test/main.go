//go:build wireinject

package wire_test

import "github.com/google/wire"

type Animal interface {
	Eat() string
}

type Dog struct{}

func (d *Dog) Eat() string {
	panic("implement me")
}

func NewDog() Animal {
	return &Dog{}
}

type Cat struct{}

func (c *Cat) Eat() string {
	panic("implement me")
}

func NewCat() Animal {
	return &Cat{}
}

type Ha struct {
	animal Animal
}

func NewHa(animal Animal) *Ha {
	return &Ha{animal: animal}
}

type He struct {
	animal Animal
}

func NewHe(animal Animal) *He {
	return &He{animal: animal}
}

type End struct {
	He *He
	Ha *Ha
}

func NewEnd(he *He, ha *Ha) *End {
	return &End{He: he, Ha: ha}
}

func Init() *End {
	wire.Build(NewDog, NewCat, NewHa, NewHa, NewEnd)
	return new(End)
}
