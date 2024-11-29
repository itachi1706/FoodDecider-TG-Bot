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
	"log"
	"strconv"
	"strings"
)

func ListFoodsCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListFoods command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get first 5 food results with status A
	foods := repo.FindAllActiveFoodPaginated(5, 0)
	message := populateListFoodMessage(foods)

	return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend(constants.FoodList, 0, true, true))
}

func populateListFoodMessage(foods []model.Food) string {
	message := "No foods found"
	if len(foods) > 0 {
		message = "List of foods:\n\n"
		for _, food := range foods {
			desc := food.Description
			if desc == "" {
				desc = "No description provided"
			}
			message += fmt.Sprintf("ID: %s\nName: %s\nDescription: %s\n\n", food.ID.String(), food.Name, desc)
		}
	}
	return message

}

func ListFoodsCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListFoods previous button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)

	pageCnt, err := strconv.Atoi(strings.Replace(cb.Data, "previous-food-list-", "", -1))
	if err != nil {
		log.Printf("failed to convert page number: %w. Default 0\n", err)
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
	repo := repository.NewFoodsRepository(db)
	foods := repo.FindAllActiveFoodPaginated(5, pageCnt)

	message := populateListFoodMessage(foods)
	_, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit(constants.FoodList, pageCnt, true, true))

	return nil
}

func ListFoodsCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListFoods next button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)

	pageCnt, err := strconv.Atoi(strings.Replace(cb.Data, "next-food-list-", "", -1))
	if err != nil {
		log.Printf("failed to convert page number: %w. Default 0\n", err)
		pageCnt = 0
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get total number of food and find number of possible pages (including partial)
	count := repo.GetFoodCount()
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
	foods := repo.FindAllActiveFoodPaginated(5, pageCnt)

	message := populateListFoodMessage(foods)
	_, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit(constants.FoodList, pageCnt, true, true))

	return nil
}
