package rediservice

import "testing"

func TestZhenaiDuplicate(t *testing.T) {
	_, err := RedisServer.ZhenaiDuplicate("test")
	if err != nil {
		t.Errorf(err.Error())
	}
}
