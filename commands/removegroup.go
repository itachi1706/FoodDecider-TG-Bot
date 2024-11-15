package commands

import (
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func RemoveGroupCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RemoveGroup command called by " + ctx.EffectiveSender.Username())

	userId, foodId, groupName, err := services.GroupHandlingParameter(bot, ctx, "/removegroup <id> <group name>")
	if err != nil {
		return err
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if food name already exists
	food := repo.FindFoodById(*foodId)
	group := repo.GetActiveFoodGroup(*groupName)
	message := "An error has occurred. Please try again later"
	if food == nil {
		// New Food
		message = "Food with ID " + foodId.String() + " does not exist"
	} else if group == nil {
		message = "Group " + *groupName + " does not exist"
	} else {
		link := repo.GetActiveFoodGroupLink(food.ID, group.ID)
		if link == nil {
			message = "Group " + *groupName + " is not linked to food " + food.Name
		} else {
			link.Status = "D"
			link.UpdatedBy = *userId
			db.Save(&link)
			message = "Group " + *groupName + " removed from food " + food.Name
		}
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
