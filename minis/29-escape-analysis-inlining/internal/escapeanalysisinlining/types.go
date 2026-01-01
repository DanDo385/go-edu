package escapeanalysisinlining

// Rectangle is a small struct used in receiver/inlining exercises.
type Rectangle struct {
	Width  float64
	Height float64
}

// Config is a tiny config struct used in escape-analysis examples.
type Config struct {
	Host string
	Port int
}
