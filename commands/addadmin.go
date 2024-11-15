package commands

import (
    "FoodDecider-TG-Bot/model"
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "log"
    "strconv"
    "strings"
)

func AddAdminCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("AddAdmin command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    // Get user information to add to admin list from message
    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)

    if len(messageOpts) < 2 {
        return utils.BasicReplyToUser(bot, ctx, "Please provide a user id and nickname\n\nFormat: /addadmin <telegram id> <nickname>")
    }

    userIdToAddStr := messageOpts[0]
    // The rest of array add as nickname
    nickname := strings.Join(messageOpts[1:], " ")

    userIdToAdd, err := strconv.ParseInt(userIdToAddStr, 10, 64)
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid user ID format")
    }

    db := utils.GetDbConnection()
    admin := repository.FindAdmin(db, userIdToAdd)
    message := "An error has occurred. Please try again later"
    if admin == nil {
        // New user
        log.Println("Adding new user " + userIdToAddStr + " to admin list")
        admin = &model.Admins{
            TelegramID: userIdToAdd,
            Name:       nickname,
            CreatedBy:  userId,
            UpdatedBy:  userId,
        }
        db.Create(admin)
        message = "User " + nickname + " (" + userIdToAddStr + ") added as admin"
    } else {
        // User already exists
        log.Println("User already exists in admin list. Check if is admin already")
        if admin.Status == "A" {
            message = "User " + nickname + " (" + userIdToAddStr + ") already an admin"
            return utils.BasicReplyToUser(bot, ctx, message)
        } else {
            log.Println("User is not an admin. Updating user to admin")
            admin.Status = "A"
            admin.UpdatedBy = userId

            db.Save(&admin)
            message = "User " + nickname + " (" + userIdToAddStr + ") reinstated as admin"
        }
    }

    return utils.BasicReplyToUser(bot, ctx, message)
}
