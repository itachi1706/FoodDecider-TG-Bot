package services

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/utils"
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strings"
)

func RunPreCommandScripts(ctx *ext.Context) {
	RunPreCommandScriptCustomType(ctx, constants.COMMAND)
}

func RunPreCommandScriptCustomType(ctx *ext.Context, messageType constants.MessageType) {
	go func() {
		uuid := LogUserFound(ctx)
		message := ctx.EffectiveMessage
		user := ctx.EffectiveSender
		chatId := ctx.EffectiveChat.Id

		messageString, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message object: %s. Skipping", err)
		}

		commandLog := &model.CommandsLog{
			UserID:    uuid,
			ChatId:    chatId,
			Type:      messageType,
			CreatedBy: user.Id(),
			RawData:   string(messageString),
			Command:   "Unknown Command",
		}

		if messageType == constants.COMMAND {
			// Get arguments
			args := utils.GetArgumentsFromMessage(ctx)
			cmd := utils.GetCommandFromMessage(ctx)
			commandLog.Command = cmd
			commandLog.Arguments = strings.Join(args, ", ")
		} else if messageType == constants.CALLBACK {
			cb := ctx.Update.CallbackQuery
			commandLog.Command = cb.Data

			callbackString, err := json.Marshal(cb)
			if err != nil {
				log.Printf("Failed to marshal callback object: %s. Skipping", err)
			}
			commandLog.ExtraData = string(callbackString)
		} else if messageType == constants.CONVERSATION {
			if ctx.EffectiveMessage.Text != "" {
				commandLog.Command = message.Text
			} else if ctx.EffectiveMessage.Location != nil {
				commandLog.Command = "Location"
				commandLog.Arguments = fmt.Sprintf("%v,%v", message.Location.Latitude, message.Location.Longitude)
			} else {
				commandLog.Command = "Unknown Conversation Type. See Raw Data"
			}
		}

		db := utils.GetDbConnection()
		db.Save(commandLog)
	}()
}
