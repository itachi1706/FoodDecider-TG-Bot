package commands

import (
    "FoodDecider-TG-Bot/model"
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/google/uuid"
    "log"
    "strings"
)

func UpdateFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("UpdateFood command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 3 {
        return utils.BasicReplyToUser(bot, ctx, "Invalid update food format\n\nFormat: /updatefood <id> <name/description> [value]")
    }

    foodId, err := uuid.Parse(messageOpts[0])
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
    }

    updateType := strings.ToLower(messageOpts[1])
    log.Printf("Update type: '%s'\n", updateType)
    if updateType != "name" && updateType != "description" {
        return utils.BasicReplyToUser(bot, ctx, "Invalid update type provided. Valid types: name, description")
    }

    updateValue := strings.Trim(strings.Join(messageOpts[2:], " "), " ")

    db := utils.GetDbConnection()
    var food model.Food
    // Check if food name already exists
    result := db.Where("id = ? AND status = ?", foodId, "A").First(&food)
    message := "An error has occurred. Please try again later"
    if result.Error != nil {
        // New Food
        message = "Food with ID " + foodId.String() + " does not exist"
    } else {
        if updateType == "name" {
            food.Name = updateValue
        } else {
            food.Description = updateValue
        }

        food.UpdatedBy = userId
        db.Save(&food)
        message = "Food " + food.Name + " updated in database"
    }

    return utils.BasicReplyToUser(bot, ctx, message)
}
