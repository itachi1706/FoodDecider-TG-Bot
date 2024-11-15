package utils

import (
    "fmt"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func BasicReplyToUser(bot *gotgbot.Bot, ctx *ext.Context, message string) error {
    return replyUser(bot, ctx, message, nil)
}

func BasicReplyToUserWithMarkdown(bot *gotgbot.Bot, ctx *ext.Context, message string) error {
    return replyUser(bot, ctx, message, &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
}

func BasicReplyToUserWithHTML(bot *gotgbot.Bot, ctx *ext.Context, message string) error {
    return replyUser(bot, ctx, message, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
}

func ReplyUserWithOpts(bot *gotgbot.Bot, ctx *ext.Context, message string, opts *gotgbot.SendMessageOpts) error {
    return replyUser(bot, ctx, message, opts)
}

func replyUser(bot *gotgbot.Bot, ctx *ext.Context, message string, opt *gotgbot.SendMessageOpts) error {
    _, err := ctx.EffectiveMessage.Reply(bot, message, opt)
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }

    return nil
}

func GetArgumentsFromMessage(ctx *ext.Context) []string {
    message := ctx.EffectiveMessage.Text
    messageList := SplitString(message)
    // Remove the command from the list
    return messageList[1:]
}
