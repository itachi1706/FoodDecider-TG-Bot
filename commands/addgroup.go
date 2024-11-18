package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func AddGroupCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("AddGroup command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	userId, foodId, groupName, err := services.GroupHandlingParameter(bot, ctx, "/addgroup <food id> <group name>")
	if err != nil {
		return err
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	food := repo.FindFoodById(*foodId)
	// Check if food exists
	message := constants.ErrorMessage
	if food == nil {
		// New Food
		message = "Food not found"
		return utils.BasicReplyToUser(bot, ctx, message)
	}

	// Check if group exists
	group := repo.GetFoodGroup(*groupName)
	if group == nil {
		log.Println("Creating new group " + *groupName)
		// New Group, create the group
		group = &model.FoodGroups{
			Name:      *groupName,
			CreatedBy: *userId,
			UpdatedBy: *userId,
		}
		db.Create(&group)
	} else {
		// Group exists, if deleted, reactivate group
		if group.Status != "A" {
			log.Println("Reactivating group " + *groupName)
			group.Status = "A"
			group.UpdatedBy = *userId
			db.Save(&group)
		}
	}

	// Add the link if it does not exists
	message = "Linked food " + food.Name + " to " + *groupName
	link := repo.GetFoodGroupLink(*foodId, group.ID)
	if link != nil {
		if link.Status == "A" {
			message = "Food " + food.Name + " is already linked to " + *groupName
		} else {
			log.Println("Reactivating link between food " + food.Name + " and group " + *groupName)
			link.Status = "A"
			link.UpdatedBy = *userId
			db.Save(&link)
		}
	} else {
		log.Println("Linking food " + food.Name + " to group " + *groupName)
		link = &model.FoodGroupsLink{
			FoodID:    *foodId,
			GroupID:   group.ID,
			CreatedBy: *userId,
			UpdatedBy: *userId,
		}
		db.Create(&link)
	}

	return utils.BasicReplyToUser(bot, ctx, message)
}
