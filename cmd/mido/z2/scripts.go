package z2

import "strings"

type ScriptDecl string

type Scripts struct {
	entities NamedEntities
}

func (this Scripts) Load(decl, body string) {
	name, _, _ := strings.Cut(decl, "(")
	entity := this.entities.Entity(Name(name))
	entity.Attach(ScriptDecl(decl), StringSource(body))
}
