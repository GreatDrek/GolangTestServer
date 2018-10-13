package myLib

import (
	"encoding/json"
)

type User interface {
	ReturnInfo() []byte
	AddText(string)
}

type OneUser struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func (ou *OneUser) ReturnInfo() []byte {
	js, err := json.Marshal(*ou)
	if err != nil {
		//fmt.Fatal("Error:", err)
		panic("error")
	}
	return js
}

func (ou *OneUser) AddText(s string) {
	(*ou).Text = s
}
