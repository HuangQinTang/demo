package logicerror

import (
	"fmt"
	"chat/library/config"
)

func PrintError(err error) error {
	if err != nil && (config.GetConfig().App.ClientDebug || config.GetConfig().App.ServerDebug) {
		fmt.Println(err.Error())
	}
	return err
}
