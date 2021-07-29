package game

type ruleFunc func(cell, *Game) (cellType, bool)

type composeRule struct {
	rules []ruleFunc
}

func newComposeRule(rules ...ruleFunc) composeRule {
	return composeRule{rules: rules}
}

func (r composeRule) getCellType(c cell, game *Game) cellType {
	for _, rule := range r.rules {
		tp, changed := rule(c, game)
		if changed {
			return tp
		}
	}
	return c.cellType
}

func aliveLessThanTwo(c cell, game *Game) (cellType, bool) {
	alive := 0
	for _, c := range game.getNeighbors(c) {
		if c.cellType == ALIVE {
			alive++
		}
	}
	if c.cellType == ALIVE && alive < 2 {
		return DEAD, true
	}
	return c.cellType, false
}

func aliveTwoOrThree(c cell, game *Game) (cellType, bool) {
	alive := 0
	for _, c := range game.getNeighbors(c) {
		if c.cellType == ALIVE {
			alive++
		}
	}
	if c.cellType == ALIVE && (alive == 2 || alive == 3) {
		return ALIVE, true
	}
	return c.cellType, false
}

func aliveMoreThanThree(c cell, game *Game) (cellType, bool) {
	alive := 0
	for _, c := range game.getNeighbors(c) {
		if c.cellType == ALIVE {
			alive++
		}
	}
	if c.cellType == ALIVE && alive > 3 {
		return DEAD, true
	}
	return c.cellType, false
}

func deadThreeAlive(c cell, game *Game) (cellType, bool) {
	alive := 0
	for _, c := range game.getNeighbors(c) {
		if c.cellType == ALIVE {
			alive++
		}
	}
	if c.cellType == DEAD && alive == 3 {
		return ALIVE, true
	}
	return c.cellType, false
}

var defaultRule = newComposeRule(aliveLessThanTwo, aliveTwoOrThree, aliveMoreThanThree, deadThreeAlive)
