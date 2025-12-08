package golobby

func MustSingleton(c Container, resolver interface{}) {
	if err := c.Singleton(resolver); err != nil {
		panic(err)
	}
}

// MustCall wraps the `Call` method and panics on errors instead of returning the errors.
func MustCall(c Container, receiver interface{}) any {
	result, err := c.Call(receiver)
	if err != nil {
		panic(err)
	}
	return result
}

// MustResolve wraps the `Resolve` method and panics on errors instead of returning the errors.
func MustResolve(c Container, abstraction interface{}) {
	if err := c.Resolve(abstraction); err != nil {
		panic(err)
	}
}

// MustNamedResolve wraps the `NamedResolve` method and panics on errors instead of returning the errors.
func MustNamedResolve(c Container, abstraction interface{}, name string) {
	if err := c.NamedResolve(abstraction, name); err != nil {
		panic(err)
	}
}
