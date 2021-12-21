package battle

type Group struct {
	timelines []TimelineObj
}

func (g *Group) TimelineTick(interval float64) {
	if len(g.timelines) <= 0 {
		return
	}
	ori := len(g.timelines)

	for idx := 0; idx < ori; {
		wasTimeElapsed := g.timelines[idx].timeElapsed
		g.timelines[idx].timeElapsed += interval * g.timelines[idx].timeScale
		if g.timelines[idx].model.chargeGoBack.atDuration < g.timelines[idx].timeScale &&
			g.timelines[idx].model.chargeGoBack.atDuration >= wasTimeElapsed {

			// todo

		}

		for i := range g.timelines[idx].model.nodes {
			if g.timelines[idx].model.nodes[i].timeElapsed < g.timelines[idx].timeElapsed &&
				g.timelines[idx].model.nodes[i].timeElapsed >= wasTimeElapsed {
				g.timelines[idx].model.nodes[i].doEvent(
					g.timelines[idx],
					g.timelines[idx].model.nodes[i].eveParams,
				)
			}
		}
		if g.timelines[idx].model.duration <= g.timelines[idx].timeElapsed {
			g.timelines = append(g.timelines[:idx], g.timelines[idx+1:]...)
			ori--
		} else {
			idx++
		}
	}
}
