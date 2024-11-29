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
	"github.com/google/uuid"
	"log"
)

const (
	CoordinateListGrp = "coordinate-list+"
)

func ListCoordinatesCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListCoordinates command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)
	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Food ID required\n\nFormat: /listcoordinates <food id>")
	}

	foodId, err := uuid.Parse(messageOpts[0])
	if err != nil {
		return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get first 5 food results with status A
	foodGroups := repo.FindAllLocationsForFoodPaginated(foodId, 5, 0)
	food := repo.FindFoodById(foodId)
	message := populateListFoodLocationsMessage(foodGroups, food)

	return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend(CoordinateListGrp+foodId.String()+"+", 0, true, true))
}

func ListCoordinatesCommandTrigger(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListCoordinates trigger button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)
	foodIdStr := cb.Data[len("list-coordinates-"):]
	foodId, err := uuid.Parse(foodIdStr)
	if err != nil {
		_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: constants.RerollError,
		})
		return err
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get first 5 food results with status A
	foodGroups := repo.FindAllLocationsForFoodPaginated(foodId, 5, 0)
	food := repo.FindFoodById(foodId)
	message := populateListFoodLocationsMessage(foodGroups, food)

	return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend(CoordinateListGrp+foodId.String()+"+", 0, true, true))
}

func populateListFoodLocationsMessage(groups []model.Locations, food *model.Food) string {
	foodName := "Unknown Food"
	if food != nil {
		foodName = food.Name
	}

	message := "No locations found for " + foodName
	if len(groups) > 0 {
		message = "Locations for " + foodName + ":\n\n"
		for _, group := range groups {
			name := group.Name
			if name == "" {
				name = "No name defined"
			}
			message += fmt.Sprintf("ID: %s\nLocation: %v, %v\nName: %s\n\n", group.ID, group.Latitude, group.Longitude, name)
		}
	}
	return message
}

func getCoordinatesData(foodId *uuid.UUID, pageCnt *int, bot *gotgbot.Bot, cb *gotgbot.CallbackQuery) error {
	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	foodLocations := repo.FindAllLocationsForFoodPaginated(*foodId, 5, *pageCnt)
	food := repo.FindFoodById(*foodId)
	message := populateListFoodLocationsMessage(foodLocations, food)
	_, _, err := cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit(CoordinateListGrp+foodId.String()+"+", *pageCnt, true, true))

	return err
}

func ListCoordinatesCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListCoordinates previous button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)

	err, foodId, pageCnt := services.HandleFoodPrevCommands(bot, cb)
	if err != nil || foodId == nil || pageCnt == nil {
		return err // End here
	}

	// Get previous 5 food results with status A
	return getCoordinatesData(foodId, pageCnt, bot, cb)
}

func ListCoordinatesCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListGroups next button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	cb := ctx.Update.CallbackQuery
	log.Println(constants.CallbackDataLog + cb.Data)

	err, foodId, pageCnt := services.HandleFoodNextCommands(bot, cb)
	if err != nil || foodId == nil || pageCnt == nil {
		return err // End here
	}

	// Get next 5 food results with status A
	return getCoordinatesData(foodId, pageCnt, bot, cb)
}
