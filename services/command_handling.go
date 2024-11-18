package services

import "github.com/PaulSonOfLars/gotgbot/v2/ext"

func RunPreCommandScripts(ctx *ext.Context) {
	go LogUserFound(ctx)
}
