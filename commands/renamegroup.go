package commands

import (
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/services"
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
    "log"
    "strings"
)

const (
    NewGroupName = "newgroupname-rename"
)

func RenameGroupCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("RenameGroup command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 1 {
        return utils.BasicReplyToUser(bot, ctx, "Please enter old group name\n\nFormat: /updatefood <group name>")
    }

    groupName := strings.Trim(strings.Join(messageOpts[0:], " "), " ")

    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    // Check if foodGroup name already exists
    foodGroup := repo.GetActiveFoodGroup(groupName)
    message := "An error has occurred. Please try again later"
    advance := false
    if foodGroup == nil {
        // New Food
        message = "Group " + groupName + " does not exist"
    } else {
        // Enter conversation to get new group name
        advance = true
        message = "Please enter new group name. Run /cancel to cancel this operation"
    }

    id, _ := conversation.KeyStrategySenderAndChat(ctx)
    err := utils.BasicReplyToUser(bot, ctx, message)
    if !advance {
        return err
    } else {
        services.SetString("groupId-"+id, foodGroup.ID)
        return handlers.NextConversationState(NewGroupName)
    }
}

func RenameGroupCommandNewName(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("RenameGroup command with new name called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        _ = utils.BasicReplyToUser(bot, ctx, "This operation can only be done by an administrator. Use /cancel to cancel this operation")
        return handlers.NextConversationState(NewGroupName)
    }

    newGroupName := ctx.EffectiveMessage.Text
    cid, _ := conversation.KeyStrategySenderAndChat(ctx)

    idIf, get := services.GetString("groupId-" + cid)
    if !get {
        _ = utils.BasicReplyToUser(bot, ctx, "An error has occurred. Please try again later")
        return handlers.EndConversation()
    }
    services.DeleteString("groupId-" + cid)
    id := idIf.(int)

    log.Println("New group name: " + newGroupName)
    log.Printf("Group ID: %d\n", id)

    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    // Check if foodGroup name already exists
    foodGroup := repo.GetActiveFoodGroupById(id)
    message := "An error has occurred. Please try again later"
    if foodGroup == nil {
        // New Food
        message = "Group does not exist"
    } else {
        // Update group name
        foodGroup.Name = newGroupName
        foodGroup.UpdatedBy = userId
        db.Save(&foodGroup)
        message = "Group has been renamed successfully"
    }

    _ = utils.BasicReplyToUser(bot, ctx, message)
    return handlers.EndConversation()
}
