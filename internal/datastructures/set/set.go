package set

type Id interface {
	comparable
}

type entry struct{}

type Hash[T Id] map[T]entry

type Constraint[T Id] interface {
	~map[T]entry
}

// removes keys from the target map that do not exist in the provided set
func IntersectMap[T Id, U any](target map[T]U, using Hash[T]) {
	discards := make(Hash[T], len(target))

	for k := range target {
		if !using.Exists(k) {
			discards.Add(k)
		}
	}

	if len(discards) >= 1 {
		for k := range discards {
			delete(target, k)
		}
	}
}

// creates a new hashset by mapping over a slice to return the values
func MapFromSlice[T Id, U any](from []U, g func(U) T) Hash[T] {
	hash := make(Hash[T], len(from))

	for _, u := range from {
		hash[g(u)] = entry{}
	}

	return hash
}

func FromSlice[T Id](from []T) Hash[T] {
	hash := make(Hash[T], len(from))
	for _, t := range from {
		hash[t] = entry{}
	}
	return hash
}

func FromMap[T Id, U any](from map[T]U) Hash[T] {
	hash := make(Hash[T], len(from))

	for k := range from {
		hash[k] = entry{}
	}

	return hash
}

// creates a new hashset w/ default capacity
func New[T Id]() Hash[T] {
	return make(map[T]entry)
}

func (s Hash[T]) Add(t T) {
	s[t] = entry{}
}

func (s Hash[T]) Exists(t T) bool {
	_, ok := s[t]
	return ok
}

// allows operations between types that wrap a hashset
type Operation[E Id, S Constraint[E], T Constraint[E]] func(s S, t T) Hash[E]

// return a new set with only items present in both sets
func Intersection[E Id, S Constraint[E], T Constraint[E]](s S, t T) Hash[E] {
	u := make(Hash[E], len(s))

	for k := range s {
		if _, ok := t[k]; ok {
			u[k] = struct{}{}
		}
	}

	return u
}

// return a new set that contains all keys in S that are not in T
func Difference[E Id, S Constraint[E], T Constraint[E]](s S, t T) Hash[E] {
	d := make(Hash[E], len(s))

	for k := range s {
		if _, ok := t[k]; !ok {
			d[k] = entry{}
		}
	}

	return d
}

// returns new set that contains all keys from S and T
func Union[E Id, S Constraint[E], T Constraint[E]](s S, t T) Hash[E] {
	d := make(Hash[E], len(s)+len(t))

	for k, v := range s {
		d[k] = v
	}

	for k, v := range t {
		d[k] = v
	}

	return d
}

// returns true if S and T have the same elements recorded
func Equal[E Id, S Constraint[E], T Constraint[E]](s S, t T) bool {
	if len(s) != len(t) {
		return false
	}

	for k := range s {
		if _, ok := t[k]; !ok {
			return false
		}
	}

	return true
}

// transforms a hashset into a slice, order is not gauranteed
func AsSlice[E Id, S Constraint[E]](s S) []E {
	arr := make([]E, 0, len(s))
	for k := range s {
		arr = append(arr, k)
	}
	return arr
}
