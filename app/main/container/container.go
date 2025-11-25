package container

import "github.com/brunobotter/chat-websocket/main/container/golobby"

type Container interface {
	NamedSingleton(name string, resolver interface{})
	Singleton(resolver interface{})
	Transcient(resolver interface{})
	NamedTranscient(name string, resolver interface{})
	Resolve(abstraction interface{})
	NamedResolve(abstraction interface{}, name string)
	Call(function interface{}) any
	Fill(structure interface{})
}

type golobbyContainerAdapter struct {
	container golobby.Container
}

func (a *golobbyContainerAdapter) NamedSingleton(name string, resolver interface{}) {
	golobby.MustNamedSingleton(a.container, name, resolver)
}

func (a *golobbyContainerAdapter) Singleton(resolver interface{}) {
	golobby.MustSingleton(a.container, resolver)
}

func (a *golobbyContainerAdapter) Transcient(resolver interface{}) {
	golobby.MustTransient(a.container, resolver)
}
func (a *golobbyContainerAdapter) NamedTranscient(name string, resolver interface{}) {
	golobby.MustNamedTransient(a.container, name, resolver)
}
func (a *golobbyContainerAdapter) Resolve(abstraction interface{}) {
	golobby.MustResolve(a.container, abstraction)
}
func (a *golobbyContainerAdapter) NamedResolve(abstraction interface{}, name string) {
	golobby.MustNamedResolve(a.container, abstraction, name)
}

func (a *golobbyContainerAdapter) Fill(structure interface{}) {
	golobby.MustFill(a.container, structure)
}

func (a *golobbyContainerAdapter) Call(function interface{}) any {
	return golobby.MustCall(a.container, function)
}

func NewContainer() Container {
	c := &golobbyContainerAdapter{container: golobby.New()}

	c.Singleton(func() Container {
		return c
	})
	return c
}
