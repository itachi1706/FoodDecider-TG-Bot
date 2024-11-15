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

func GeneratePageKeys(cmdType string, currentPage int, showPrev bool, showNext bool) gotgbot.InlineKeyboardMarkup {
    var keys [][]gotgbot.InlineKeyboardButton
    var row []gotgbot.InlineKeyboardButton
    if showPrev {
        row = append(row, gotgbot.InlineKeyboardButton{
            Text:         "⬅️",
            CallbackData: fmt.Sprintf("previous-%s-%d", cmdType, currentPage),
        },
        )
    }
    if showNext {
        row = append(row, gotgbot.InlineKeyboardButton{
            Text:         "➡️",
            CallbackData: fmt.Sprintf("next-%s-%d", cmdType, currentPage),
        },
        )
    }
    keys = append(keys, row)

    return gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: keys,
    }
}

func GeneratePageKeysSend(cmdType string, currentPage int, showPrev bool, showNext bool) *gotgbot.SendMessageOpts {
    return &gotgbot.SendMessageOpts{ReplyMarkup: GeneratePageKeys(cmdType, currentPage, showPrev, showNext)}
}

func GeneratePageKeysEdit(cmdType string, currentPage int, showPrev bool, showNext bool) *gotgbot.EditMessageTextOpts {
    return &gotgbot.EditMessageTextOpts{ReplyMarkup: GeneratePageKeys(cmdType, currentPage, showPrev, showNext)}
}
