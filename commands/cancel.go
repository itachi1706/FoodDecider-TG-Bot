package commands

import (
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
    "log"
)

func CancelCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("Cancel command called by " + ctx.EffectiveSender.Username())
    _ = utils.BasicReplyToUser(bot, ctx, "Command cancelled")

    return handlers.EndConversation()

}
