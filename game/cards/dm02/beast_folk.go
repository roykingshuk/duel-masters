package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BarkwhipTheSmasher ...
func BarkwhipTheSmasher(c *match.Card) {

	c.Name = "Barkwhip, the Smasher"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = family.BeastFolk
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetPowerEvent); ok {

			if card.Zone != match.BATTLEZONE || !card.Tapped || event.Card == card {
				return
			}

			if event.Card.Family == card.Family {
				event.Power += 2000
			}

		}

	})

}
