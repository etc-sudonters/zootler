package set

type Id interface {
	~int
}

type entry struct{}

var marker entry = struct{}{}

type Hash[T Id] map[T]struct{}

func New[T Id]() Hash[T] {
	return make(map[T]struct{})
}

func (s Hash[T]) Add(t T) {
	s[t] = marker
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
