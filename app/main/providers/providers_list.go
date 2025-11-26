package providers

func List() []any {
	return []any{
		NewConfigServiceProvider(),
		NewRedisServiceProvider(),
		NewHubServiceProvider(),
		NewCliServiceProvider(),
	}
}
