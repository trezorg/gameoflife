package writer

type item interface {
	Positions() (int, int)
}
