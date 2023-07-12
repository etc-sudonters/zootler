package set

type Id interface {
	comparable
}

type entry struct{}

type Hash[T Id] map[T]entry

type Constraint[T Id] interface {
	~map[T]entry
}

func DiscardUsing[T Id, U any](target map[T]U, using Hash[T]) {
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

func MappedFromSlice[T Id, E any](src []E, f func(E) T) Hash[T] {
	hash := make(Hash[T], len(src))
	for _, t := range src {
		hash[f(t)] = entry{}
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

type Operation[E Id, S Constraint[E], T Constraint[E]] func(s S, t T) Hash[E]

func Intersection[E Id, S Constraint[E], T Constraint[E]](s S, t T) Hash[E] {
	u := make(Hash[E], len(s))

	for k := range s {
		if _, ok := t[k]; ok {
			u[k] = struct{}{}
		}
	}

	return u
}

func Difference[E Id, S Constraint[E], T Constraint[E]](s S, t T) Hash[E] {
	d := make(Hash[E], len(s))

	for k := range s {
		if _, ok := t[k]; !ok {
			d[k] = entry{}
		}
	}

	return d
}

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

func AsSlice[E Id, S Constraint[E]](s S) []E {
	arr := make([]E, 0, len(s))
	for k := range s {
		arr = append(arr, k)
	}
	return arr
}
