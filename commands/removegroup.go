package commands

import (
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/google/uuid"
    "log"
    "strings"
)

func RemoveGroupCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("RemoveGroup command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 2 {
        return utils.BasicReplyToUser(bot, ctx, "Invalid Format\n\nFormat: /removegroup <id> <group name>")
    }

    foodId, err := uuid.Parse(messageOpts[0])
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
    }

    groupName := strings.Trim(strings.Join(messageOpts[1:], " "), " ")

    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    // Check if food name already exists
    food := repo.FindFoodById(foodId)
    group := repo.GetActiveFoodGroup(groupName)
    message := "An error has occurred. Please try again later"
    if food == nil {
        // New Food
        message = "Food with ID " + foodId.String() + " does not exist"
    } else if group == nil {
        message = "Group " + groupName + " does not exist"
    } else {
        link := repo.GetActiveFoodGroupLink(food.ID, group.ID)
        if link == nil {
            message = "Group " + groupName + " is not linked to food " + food.Name
        } else {
            link.Status = "D"
            link.UpdatedBy = userId
            db.Save(&link)
            message = "Group " + groupName + " removed from food " + food.Name
        }
    }

    return utils.BasicReplyToUser(bot, ctx, message)
}
