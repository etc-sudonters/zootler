package set

type Id interface {
	comparable
}

type entry interface{}

type Hash[T Id] map[T]entry

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
		hash[g(u)] = nil
	}

	return hash
}

func FromSlice[T Id](from []T) Hash[T] {
	hash := make(Hash[T], len(from))
	for _, t := range from {
		hash[t] = nil
	}
	return hash
}

func FromMap[T Id, U any](from map[T]U) Hash[T] {
	hash := make(Hash[T], len(from))

	for k := range from {
		hash[k] = nil
	}

	return hash
}

func New[T Id]() Hash[T] {
	return make(map[T]entry)
}

func (s Hash[T]) Add(t T) {
	s[t] = nil
}

func (s Hash[T]) Exists(t T) bool {
	_, ok := s[t]
	return ok
}

type Operation[T Id] func(s, t Hash[T]) Hash[T]

func Intersection[T Id](s, t Hash[T]) Hash[T] {
	u := make(Hash[T], len(s))

	for k := range s {
		if _, ok := t[k]; ok {
			u[k] = struct{}{}
		}
	}

	return u
}

func Difference[T Id](s, t Hash[T]) Hash[T] {
	d := make(Hash[T], len(s))

	for k, v := range s {
		if _, ok := t[k]; !ok {
			d[k] = v
		}
	}

	return d
}

func Union[T Id](s, t Hash[T]) Hash[T] {
	d := make(Hash[T], len(s)+len(t))

	for k, v := range s {
		d[k] = v
	}

	for k, v := range t {
		d[k] = v
	}

	return d
}
