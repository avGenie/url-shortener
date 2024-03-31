package model

type Batch []BatchObject

type BatchObject struct {
	ID       string
	InputURL string
	ShortURL string
}
