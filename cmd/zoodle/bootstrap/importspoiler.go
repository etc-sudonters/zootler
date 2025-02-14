package bootstrap

import (
	"context"
	"fmt"
	"io"
	"sudonters/libzootr/internal/json"
	"sudonters/libzootr/internal/settings"
	"sudonters/libzootr/magicbean/tracking"

	"github.com/etc-sudonters/substrate/dontio"
)

func malformed(err error) error {
	return fmt.Errorf("malformed spoilers document: %w", err)
}

func LoadSpoilerData(
	ctx context.Context,
	r io.Reader,
	these *settings.Zootr,
	nodes *tracking.Nodes,
	tokens *tracking.Tokens) error {

	if std, err := dontio.StdFromContext(ctx); err == nil {
		std.WriteLineOut("Loading spoiler data...")
	}

	reader := json.NewParser(json.NewScanner(r))
	obj, err := reader.ReadObject()
	if err != nil {
		return malformed(err)
	}

	for obj.More() {
		property, err := obj.ReadPropertyName()
		if err != nil {
			return malformed(err)
		}

		switch property {
		case "settings":
			obj, err := obj.ReadObject()
			if err != nil {
				return malformed(err)
			}
			if err := CopySettings(these, obj); err != nil {
				return malformed(err)
			}
		case "locations":
			obj, err := obj.ReadObject()
			if err != nil {
				return malformed(err)
			}
			if err := CopyPlacements(tokens, nodes, obj); err != nil {
				return malformed(err)
			}
		case ":skipped_locations":
			obj, err := obj.ReadObject()
			if err != nil {
				return malformed(err)
			}
			if err := CopySkippedLocations(nodes, obj); err != nil {
				return malformed(err)
			}
		default:
			if err := obj.DiscardValue(); err != nil {
				return malformed(err)
			}
		}
	}

	if err := obj.ReadEnd(); err != nil {
		return malformed(err)
	}

	return nil
}

func CopySettings(these *settings.Zootr, obj *json.ObjectParser) error {
	return obj.DiscardRemaining()
}

func CopyPlacements(tokens *tracking.Tokens, nodes *tracking.Nodes, obj *json.ObjectParser) error {
	return obj.DiscardRemaining()
}

func CopySkippedLocations(nodes *tracking.Nodes, obj *json.ObjectParser) error {
	return obj.DiscardRemaining()
}

/*
func notimpled() error {
	panic("not implemented")
}

func (this spoiler) CopySettings(these *settings.Zootr) error {
	for name, value := range this.Settings {
		_, _ = name, value
		notimpled()
	}

	return nil
}

func (this spoiler) CopyPlacements(tokens *tracking.Tokens, nodes *tracking.Nodes) error {
	for where, placed := range this.Placements {
		placement := nodes.Placement(components.Name(where))
		var token tracking.Token
		switch placed := placed.(type) {
		case string:
			token = tokens.MustGet(components.Name(placed))
		case map[string]any:
			name := placed["item"].(string)
			token = tokens.MustGet(components.Name(name))
			if price, exists := placed["price"]; exists {
				token.Attach(components.Price(price.(int)))
			}
		default:
			return errors.New("unexpected token")
		}

		placement.Holding(token)
	}

	return notimpled()
}

func (this spoiler) MarkSkippedLocations(nodes *tracking.Nodes) error {
	for where := range this.SkippedLocations {
		placement := nodes.Placement(components.Name(where))
		placement.Attach(components.Collected{})
	}

	return nil
}
*/
