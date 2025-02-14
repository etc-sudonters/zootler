package settings

type Model struct {
	Logic      Logic
	Cosmetics  Cosmetics
	Generation Generation
	Rom        Rom
}

func Finalize(zootr *Zootr) (Model, error) {
	var m Model
	return m, notImpled
}

func FromString(encoded string) (Model, error) {
	z, err := decodeSettingStr(encoded)
	if err != nil {
		return Model{}, err
	}

	return Finalize(&z)
}

func Default() Model {
	panic(notImpled)
}
