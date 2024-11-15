package commands

import (
    "FoodDecider-TG-Bot/model"
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/google/uuid"
    "log"
    "strings"
)

func AddGroupCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("AddGroup command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 2 {
        return utils.BasicReplyToUser(bot, ctx, "Invalid Format\n\nFormat: /addgroup <food id> <group name>")
    }

    foodId, err := uuid.Parse(messageOpts[0])
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
    }

    groupName := strings.Trim(strings.Join(messageOpts[1:], " "), " ")

    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    food := repo.FindFoodById(foodId)
    // Check if food exists
    message := "An error has occurred. Please try again later"
    if food == nil {
        // New Food
        message = "Food not found"
        return utils.BasicReplyToUser(bot, ctx, message)
    }

    // Check if group exists
    group := repo.GetFoodGroup(groupName)
    if group == nil {
        log.Println("Creating new group " + groupName)
        // New Group, create the group
        group = &model.FoodGroups{
            Name:      groupName,
            CreatedBy: userId,
            UpdatedBy: userId,
        }
        db.Create(&group)
    } else {
        // Group exists, if deleted, reactivate group
        if group.Status != "A" {
            log.Println("Reactivating group " + groupName)
            group.Status = "A"
            group.UpdatedBy = userId
            db.Save(&group)
        }
    }

    // Add the link if it does not exists
    message = "Linked food " + food.Name + " to " + groupName
    link := repo.GetFoodGroupLink(foodId, group.ID)
    if link != nil {
        if link.Status == "A" {
            message = "Food " + food.Name + " is already linked to " + groupName
        } else {
            log.Println("Reactivating link between food " + food.Name + " and group " + groupName)
            link.Status = "A"
            link.UpdatedBy = userId
            db.Save(&link)
        }
    } else {
        log.Println("Linking food " + food.Name + " to group " + groupName)
        link = &model.FoodGroupsLink{
            FoodID:    foodId,
            GroupID:   group.ID,
            CreatedBy: userId,
            UpdatedBy: userId,
        }
        db.Create(&link)
    }

    return utils.BasicReplyToUser(bot, ctx, message)
}
