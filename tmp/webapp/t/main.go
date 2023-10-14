package main

type Redis struct {
	endpoint string
	port     int
}

type DBer interface {
	Open() error
	Close() (*Redis, error, int)
}

func main() {
}
