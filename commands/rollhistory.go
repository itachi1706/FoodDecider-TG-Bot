package commands

import (
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
)

func RollHistoryCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RollHistory command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	// Get roll information
	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)

	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Please provide a roll id\n\nFormat: /rollhistory <roll id>")
	}

	db := utils.GetDbConnection()
	repo := repository.NewRollsRepository(db)

	rollIdStr := messageOpts[0]
	rollId, err := uuid.Parse(rollIdStr)
	if err != nil {
		return utils.BasicReplyToUser(bot, ctx, "Invalid roll ID")
	}

	history := repo.GetAllHistoryForRolls(rollId)

	if len(history) == 0 {
		return utils.BasicReplyToUser(bot, ctx, "No history found for this roll")
	}

	foodRepo := repository.NewFoodsRepository(db)
	message := "Roll History for " + utils.EscapeMarkdownV2(rollIdStr) + " \\(Quoted is final decision\\)\n\n"
	first := true
	for _, h := range history {
		subMsgFmt := ""
		if first {
			subMsgFmt += ">>> ***%s***\n"
			first = false
		} else {
			subMsgFmt += "%s\n"
		}

		food := foodRepo.FindFoodById(h.DecidedFoodID)
		message += fmt.Sprintf(subMsgFmt, utils.EscapeMarkdownV2(food.Name))
	}

	return utils.BasicReplyToUserWithMarkdownV2(bot, ctx, message)
}
