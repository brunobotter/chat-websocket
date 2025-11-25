package providers

func List() []any {
	return []any{
		NewConfigServiceProvider(),
		NewCliServiceProvider(),
	}
}
