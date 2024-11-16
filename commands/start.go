package commands

import (
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func StartCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("Start command called by " + ctx.EffectiveSender.Username())

	return utils.BasicReplyToUser(bot, ctx, "This bot is a WIP that will decide food in the future")
}
