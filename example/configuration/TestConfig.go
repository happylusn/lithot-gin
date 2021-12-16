package configuration

import "github.com/happylusn/lithot-gin/example/controllers"

type TestConfig struct {}

func NewTestConfig() *TestConfig {
	return new(TestConfig)
}
func (t *TestConfig) PersonDefault() *controllers.Person {
	return &controllers.Person{Name:"luu"}
}
func (t *TestConfig) PersonWithArgs(n string) *controllers.Person {
	return &controllers.Person{Name: n}
}
