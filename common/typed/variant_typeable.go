package typed

type VariantTypeable interface {
	Typeable
	GetVariantsTypes() []Typeable
}
