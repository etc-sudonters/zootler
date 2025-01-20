package z2

import "sudonters/zootler/mido/objects"

type Token struct {
	proxy
	Name
}

type Tokens struct {
	Entities NamedEntities
}

func (this Tokens) Entity(name Name) Token {
	entity := this.Entities.Entity(name)
	entity.Attach(Collectable{}, objects.PtrToken)
	return Token{entity, name}
}

type TokenLoader struct {
	*Tokens
}

func (this *TokenLoader) Load(raw token) {
	token := this.Tokens.Entity(Name(raw.Name))
	var attachments attachments

	if raw.Advancement {
		attachments.add(PriorityAdvancement)
	} else if raw.Priority {
		attachments.add(PriorityMajor)
	} else if raw.Special != nil {
		if _, exists := raw.Special["junk"]; exists {
			attachments.add(PriorityJunk)
		}
	}

	switch raw.Type {
	case "BossKey", "bosskey":
		attachments.add(BossKey{})
		break
	case "Compass", "compass":
		attachments.add(Compass{})
		break
	case "Drop", "drop":
		attachments.add(Drop{})
		break
	case "DungeonReward", "dungeonreward":
		attachments.add(DungeonReward{})
		break
	case "Event", "event":
		attachments.add(Event{})
		break
	case "GanonBossKey", "ganonbosskey":
		attachments.add(GanonBossKey{})
		break
	case "HideoutSmallKey", "hideoutsmallkey":
		attachments.add(HideoutSmallKey{})
		break
	case "HideoutSmallKeyRing", "hideoutsmallkeyring":
		attachments.add(HideoutSmallKeyRing{})
		break
	case "Item", "item":
		attachments.add(Item{})
		break
	case "Map", "map":
		attachments.add(Map{})
		break
	case "Refill", "refill":
		attachments.add(Refill{})
		break
	case "Shop", "shop":
		attachments.add(Shop{})
		break
	case "SilverRupee", "silverrupee":
		attachments.add(SilverRupee{})
		break
	case "SmallKey", "smallkey":
		attachments.add(SmallKey{})
		break
	case "SmallKeyRing", "smallkeyring":
		attachments.add(SmallKeyRing{})
		break
	case "Song", "song":
		attachments.add(Song{})
		break
	case "TCGSmallKey", "tcgsmallkey":
		attachments.add(TCGSmallKey{})
		break
	case "TCGSmallKeyRing", "tcgsmallkeyring":
		attachments.add(TCGSmallKeyRing{})
		break
	case "GoldSkulltulaToken", "goldskulltulatoken":
		attachments.add(GoldSkulltulaToken{})
		break
	}

	if raw.Special != nil {
		for name, special := range raw.Special {
			// TODO turn this into more components
			_, _ = name, special
		}
	}

	token.AttachAll(attachments.v)
}
