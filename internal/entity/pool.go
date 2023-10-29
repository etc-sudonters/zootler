package entity

// responsible for the total administration of a population of models
type Pool interface {
	Queryable
	Manager
}

// responsible for creation and destruction of models
type Manager interface {
	Create() (View, error)
}

// responsible for looking either individual models or creating a subset of the
// population that matches the provided selectors
type Queryable interface {
	// return a subset of the population that matches the provided selectors
	Query([]Selector) ([]View, error)
	// load the specified components from the specified model, if a component
	// isn't attached to the model its pointer should be set to nil
	Get(Model, []interface{})
	// return the specific model from the pool
	Fetch(Model) (View, error)
}
