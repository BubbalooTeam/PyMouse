package checkers

import "pymouse/pymouse/client"

func LoadModule(bS *client.BotStruct) {
	bS.Handler.Use(SaveUsers)
	bS.Handler.Use(SaveChats)
}
