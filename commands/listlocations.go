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

func ListLocationsCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListLocations command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	_, foodId, _, err := services.FoodValidationParameterChecks(bot, ctx, 1, "Food ID required\n\nFormat: /listlocations <food id>")
	if err != nil {
		return err
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get first 5 food results with status A
	foodGroups := repo.FindAllLocationsForFoodPaginated(*foodId, 5, 0)
	food := repo.FindFoodById(*foodId)
	message := populateListFoodLocationsMessage(foodGroups, food)

	return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend(CoordinateListGrp+foodId.String()+"+", 0, true, true))
}

func ListLocationsCommandTrigger(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListLocations trigger button clicked by " + ctx.EffectiveSender.Username())
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

func populateListFoodLocationsMessage(locations []model.Locations, food *model.Food) string {
	foodName := "Unknown Food"
	if food != nil {
		foodName = food.Name
	}

	message := "No locations found for " + foodName
	if len(locations) > 0 {
		message = "Locations for " + foodName + ":\n\n"
		for _, location := range locations {
			name := location.Name
			if name == "" {
				name = "No name defined"
			}
			message += fmt.Sprintf("ID: %s\nLocation: %v, %v\nName: %s\nAddress: %s\nPlus Code: %s\n\n", location.ID, location.Latitude, location.Longitude, name, location.Address, location.PlusCode)
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

func ListLocationsCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListLocations previous button clicked by " + ctx.EffectiveSender.Username())
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

func ListLocationsCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("ListLocations next button clicked by " + ctx.EffectiveSender.Username())
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
