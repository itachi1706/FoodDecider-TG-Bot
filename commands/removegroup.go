package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
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
	message := constants.ErrorMessage
	if food == nil {
		// New Food
		message = fmt.Sprintf("Food with ID %s does not exist", foodId.String())
	} else if group == nil {
		message = fmt.Sprintf("Group %s does not exist", *groupName)
	} else {
		link := repo.GetActiveFoodGroupLink(food.ID, group.ID)
		if link == nil {
			message = fmt.Sprintf("Group %s is not linked to food %s", *groupName, food.Name)
		} else {
			link.Status = "D"
			link.UpdatedBy = *userId
			db.Save(&link)
			message = fmt.Sprintf("Group %s removed from food %s", *groupName, food.Name)
		}
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
