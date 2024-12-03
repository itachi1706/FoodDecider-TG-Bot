package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
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
	services.RunPreCommandScripts(ctx)

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)
	if len(messageOpts) < 1 {
		return utils.BasicReplyToUser(bot, ctx, "Please enter old group name\n\nFormat: /renamegroup <group name>")
	}

	groupName := strings.Trim(strings.Join(messageOpts[0:], " "), " ")

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if foodGroup name already exists
	foodGroup := repo.GetActiveFoodGroup(groupName)
	message := constants.ErrorMessage
	advance := false
	if foodGroup == nil {
		// New Food
		message = "Group " + groupName + " does not exist"
	} else {
		// Enter conversation to get new group name
		advance = true
		message = "Please reply to this message with the new group name. Run /cancel to cancel this operation"
	}

	id, _ := conversation.KeyStrategySenderAndChat(ctx)
	err := utils.BasicReplyToUser(bot, ctx, message)
	if !advance {
		return err
	} else {
		services.SetString(fmt.Sprint("groupId-%s", id), foodGroup.ID)
		return handlers.NextConversationState(NewGroupName)
	}
}

func RenameGroupCommandNewName(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RenameGroup command with new name called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CONVERSATION)

	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		_ = utils.BasicReplyToUser(bot, ctx, "This operation can only be done by an administrator. Use /cancel to cancel this operation")
		return handlers.NextConversationState(NewGroupName)
	}

	newGroupName := ctx.EffectiveMessage.Text
	cid, _ := conversation.KeyStrategySenderAndChat(ctx)

	idKey := fmt.Sprint("groupId-%s", cid)

	idIf, get := services.GetString(idKey)
	if !get {
		_ = utils.BasicReplyToUser(bot, ctx, constants.ErrorMessage)
		return handlers.EndConversation()
	}
	services.DeleteString(idKey)
	id := idIf.(int)

	log.Println("New group name: " + newGroupName)
	log.Printf("Group ID: %d\n", id)

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Check if foodGroup name already exists
	foodGroup := repo.GetActiveFoodGroupById(id)
	message := constants.ErrorMessage
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
