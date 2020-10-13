package gofm

import (
	"testing"
)

func Test_getRoomInfo(t *testing.T) {
	room, err := getRoomInfo(209202328)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(room.Info.Creator.Username)
}
