package rest

type ApplicationContainer struct {
	Services map[ServiceWrapper]ApiWrapper
}

func NewApplicationContainer() ApplicationContainer {
	container := ApplicationContainer{Services: make(map[ServiceWrapper]ApiWrapper)}
	return container
}

func (c *ApplicationContainer) AddService(wrapper ServiceWrapper, impl ApiWrapper) {
	c.Services[wrapper] = impl
}
