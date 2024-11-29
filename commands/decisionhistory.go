package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func DecisionHistoryCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DecisionHistory command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	db := utils.GetDbConnection()
	repo := repository.NewRollsRepository(db)
	history := repo.FindAllRollsByChatIdOrderRecentPaginated(ctx.EffectiveChat.Id, 5, 0)
	message := populateListDecisionsMessage(history, db)

	return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend(constants.DecisionHistoryList, 0, true, true))
}

func populateListDecisionsMessage(decisions []model.Rolls, db *gorm.DB) string {
	if len(decisions) <= 0 {
		return "No decisions found"
	}

	message := "List of Decisions in this chat:\n\n"
	for _, decision := range decisions {
		msgFmt := "Decision ID: %s\nDecision: %s\n"
		if decision.DecidedLocationID != nil {
			msgFmt += "Location: %s (%f,%f)\n"
		}
		msgFmt += "Requested Location: %s\nRequested Group Name: %s\n"
		msgFmt += "Requestor: %s (%s)\nDate: %s\n\n"

		userRepo := repository.NewUserRepository(db)
		user := userRepo.FindUser(decision.UpdatedBy)
		foodRepo := repository.NewFoodsRepository(db)
		food := foodRepo.FindFoodById(decision.DecidedFoodID)

		reqLocation := "No location requested"
		reqGroup := "No groupings requested"
		if decision.Latitude != 0 && decision.Longitude != 0 {
			reqLocation = fmt.Sprintf("%f,%f", decision.Latitude, decision.Longitude)
		}
		if decision.GroupName != "" {
			reqGroup = decision.GroupName
		}

		updatedAt := decision.UpdatedAt.Format(constants.DateTimeFormat)
		if decision.DecidedLocationID != nil {
			location := foodRepo.FindActiveLocationById(*decision.DecidedLocationID)
			message += fmt.Sprintf(msgFmt, decision.ID.String(), food.Name, location.Name, location.Latitude, location.Longitude, reqLocation, reqGroup, user.FullName, user.Username, updatedAt)
		} else {
			message += fmt.Sprintf(msgFmt, decision.ID.String(), food.Name, reqLocation, reqGroup, user.FullName, user.Username, updatedAt)
		}
	}

	return message
}

func DecisionHistoryCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DecisionHistory previous button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println("Callback data: " + cb.Data)

	pageCnt, err := strconv.Atoi(strings.Replace(cb.Data, "previous-decision-history-", "", -1))
	if err != nil {
		// Default to 0
		pageCnt = 0
	}

	answerMsg := constants.ErrorMessage
	cont := true
	if pageCnt <= 0 {
		// First page
		answerMsg = constants.FirstPage
		cont = false
	} else {
		answerMsg = constants.GoToPrevious
		pageCnt--
	}

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: answerMsg,
	})

	if err != nil {
		return fmt.Errorf("failed to answer callback: %w", err)
	}

	if !cont {
		return nil // End here
	}

	// Get previous 5 food results with status A
	db := utils.GetDbConnection()
	repo := repository.NewRollsRepository(db)
	history := repo.FindAllRollsByChatIdOrderRecentPaginated(ctx.EffectiveChat.Id, 5, pageCnt)

	message := populateListDecisionsMessage(history, db)
	_, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit(constants.DecisionHistoryList, pageCnt, true, true))

	return nil
}

func DecisionHistoryCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("DecisionHistory next button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println("Callback data: " + cb.Data)

	pageCnt, err := strconv.Atoi(strings.Replace(cb.Data, "next-decision-history-", "", -1))
	if err != nil {
		// Default to 0
		pageCnt = 0
	}

	db := utils.GetDbConnection()
	repo := repository.NewRollsRepository(db)
	// Get total number of food and find number of possible pages (including partial)
	count := repo.GetRollsCountForChatId(ctx.EffectiveChat.Id)
	totalPages := count / 5
	modulo := count % 5
	if modulo > 0 {
		totalPages++
	}

	// pagecnt to int64
	pageCnt64 := int64(pageCnt)

	answerMsg := constants.ErrorMessage
	cont := true
	if pageCnt64 >= totalPages-1 {
		// last page
		answerMsg = "You are already on the last page"
		cont = false
	} else {
		answerMsg = "Going to next page"
		pageCnt++
	}

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: answerMsg,
	})

	if err != nil {
		return fmt.Errorf("failed to answer callback: %w", err)
	}

	if !cont {
		return nil // End here
	}

	// Get next 5 food results with status A
	history := repo.FindAllRollsByChatIdOrderRecentPaginated(ctx.EffectiveChat.Id, 5, pageCnt)

	message := populateListDecisionsMessage(history, db)
	_, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit(constants.DecisionHistoryList, pageCnt, true, true))

	return nil
}
