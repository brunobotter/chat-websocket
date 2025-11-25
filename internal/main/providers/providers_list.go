package providers

func List() []any {
	return []any{
		NewCliServiceProvider(),
	}
}
