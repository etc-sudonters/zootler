package logic

type State interface {
}

type AccessRule interface {
	CanAccess() bool
}
