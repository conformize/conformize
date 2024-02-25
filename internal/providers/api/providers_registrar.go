package sdk

type ProvidersRegistrar interface {
	Register(alias string, provider ConfigurationProvider) error
	Get(name string) (ConfigurationProvider, bool)
}
