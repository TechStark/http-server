package main

import (
	"fmt"

	"github.com/techstark/http-server/cmd"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			switch err.(type) {
			case string:
				fmt.Println(err.(string))

			case error:
				fmt.Println(err.(error))

			default:
				panic(err)
			}
		}
	}()

	cmd.Execute()
}
