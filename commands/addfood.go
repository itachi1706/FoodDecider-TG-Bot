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

func AddFoodCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("AddFood command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 1 {
        return utils.BasicReplyToUser(bot, ctx, "Please provide a food name\n\nFormat: /addfood <name>")
    }

    foodName := strings.Trim(strings.Join(messageOpts[0:], " "), " ")

    db := utils.GetDbConnection()
    // Check if food name already exists
    food := repository.FindFoodByNameAll(db, foodName)
    message := "An error has occurred. Please try again later"
    if food == nil {
        // New Food
        log.Println("Adding new food " + foodName)
        food = model.Food{
            ID:        uuid.New(),
            Name:      foodName,
            CreatedBy: userId,
            UpdatedBy: userId,
        }
        db.Create(&food)
        message = "Food " + foodName + " added to database.\n\nID: " + food.ID.String() + "\n\nUse the other commands to add more details to the food"
    } else {
        // Check if status is A
        if food.Status == "A" {
            // Food already exists
            message = "Food " + foodName + " already exists. Modify food with /updatefood command"
        } else {
            // Food exists but is inactive. Update it to active
            food.Status = "A"
            food.UpdatedBy = userId
            db.Save(&food)
            message = "Food " + foodName + " added to database.\n\nID: " + food.ID.String() + "\n\nUse the other commands to add more details to the food"
        }
    }

    return utils.BasicReplyToUser(bot, ctx, message)
}
