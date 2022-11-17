package store

type Registry struct {
	UserStore UserStore
}

func NewRegistry(pool Pool) *Registry {
	return &Registry{
		UserStore: NewUserStore(pool),
	}
}
