package commands

import (
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"strings"
)

func SearchCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("Search command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	// Get the search term
	messageOpts := utils.GetArgumentsFromMessage(ctx)
	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Invalid Format\n\nFormat: /search <food>")
	}

	searchTerm := strings.Join(messageOpts, " ")

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get first 5 food results with status A
	foods := repo.FindAllActiveFoodBySearchTerm(searchTerm)
	message := populateSearchResults(foods, repo)

	return utils.BasicReplyToUser(bot, ctx, message)
}

func populateSearchResults(foods []model.Food, repo repository.FoodRepository) string {
	message := "No results found"
	if len(foods) > 0 {
		message = fmt.Sprintf("%d results found:\n\n", len(foods))
		for _, food := range foods {
			desc := food.Description
			if desc == "" {
				desc = "No description provided"
			}

			// Get groups and locations
			groupCnt := repo.GetFoodGroupForFoodCount(food.ID)
			locCnt := repo.FindAllLocationsForFoodCount(food.ID)

			message += fmt.Sprintf("ID: %s\nName: %s\nDescription: %s\nNumber of Groups: %d\nNumber of locations: %d\n\n", food.ID.String(), food.Name, desc, groupCnt, locCnt)
		}
	}
	return message

}
