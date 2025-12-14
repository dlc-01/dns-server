package model

type Message struct {
	Header    Header
	Questions []Question
	Answers   []Answer
}
