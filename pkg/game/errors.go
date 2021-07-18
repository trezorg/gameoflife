package game

type stopIteration struct{}

func (e stopIteration) Error() string {
	return "stopIteration"
}
