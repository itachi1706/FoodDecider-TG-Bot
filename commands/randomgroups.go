package commands

import (
	"FoodDecider-TG-Bot/constants"
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/services"
	"FoodDecider-TG-Bot/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/google/uuid"
	"log"
	"strings"
)

const (
	RandomGroupSpecify = "random-group-specify"
)

func RandomGroupsCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomGroups command called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScripts(ctx)

	// Get group name information
	messageArgs := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageArgs)

	group := ""
	if len(messageArgs) >= 1 {
		group = strings.Trim(strings.Join(messageArgs, " "), " ")
	}

	if group == "" {
		// No group specified, start convo
		err := utils.BasicReplyToUser(bot, ctx, "Please reply to this message with a list of groups seperated by a new line")
		if err != nil {
			log.Println("Error replying to user: " + err.Error())
			return handlers.EndConversation()
		}
		return handlers.NextConversationState(RandomGroupSpecify)
	}

	message, rollInfo, hasLoc := searchAndReplyWithGroups(ctx, []string{group})
	if rollInfo == nil {
		return utils.BasicReplyToUser(bot, ctx, message)
	}

	// Send message to user with reroll button
	messageOpts := utils.GenerateRerollKeysSend(constants.GROUP, *rollInfo, hasLoc)

	return utils.ReplyUserWithOpts(bot, ctx, message, messageOpts)
}

func searchAndReplyWithGroups(ctx *ext.Context, groups []string) (string, *model.Rolls, bool) {
	// Search groups and reply randomly
	sender := ctx.EffectiveSender
	userId := sender.Id()
	chatId := ctx.EffectiveChat.Id

	db := utils.GetDbConnection()

	groupsNL := strings.Join(groups, "\n")

	// Store in database
	rollInfo := model.Rolls{
		ID:                uuid.New(),
		Type:              constants.GROUP,
		DecidedLocationID: nil,
		CreatedBy:         userId,
		UpdatedBy:         userId,
		ChatId:            chatId,
		GroupName:         groupsNL,
	}

	rollHistory, count, err := services.RollForFood(rollInfo)
	if err != nil {
		return utils.Capitalize(err.Error()), nil, false
	}
	rollInfo.DecidedFoodID = rollHistory.DecidedFoodID

	// Save both roll and roll history
	db.Save(&rollInfo)
	db.Save(&rollHistory)

	// groups to string seperated by comma
	groupsStr := strings.Join(groups, ", ")

	// Send message to user with reroll button
	message, hasLoc := sendWithRerollButtonGroups(rollInfo, sender, count, false, groupsStr)

	return message, &rollInfo, hasLoc
}

func sendWithRerollButtonGroups(rollInfo model.Rolls, trigger *gotgbot.Sender, count int64, reroll bool, groups string) (string, bool) {
	foodRepo := repository.NewFoodsRepository(utils.GetDbConnection())
	food := foodRepo.FindFoodById(rollInfo.DecidedFoodID)
	locationCnt := foodRepo.FindAllLocationsForFoodCount(rollInfo.DecidedFoodID)

	messageFmt := "Food Decision made 🎉\n\n"
	messageFmt += "Selected food: %s\n"
	messageFmt += "Description: %s\n\n"
	if locationCnt > 0 {
		messageFmt += "There are %d locations found for this option. Click the button below to view more\n\n"
	} else {
		messageFmt += "There are no locations found for this option. Please go online to find your nearest location yourself!\n\n"
	}
	messageFmt += "This was randomized from a list of %d food options with the following groups:\n%s\n\n"
	if reroll {
		messageFmt += "Decision was re-ran on %s by %s (%s)"
	} else {
		messageFmt += "Decision was ran on %s by %s (%s)"
	}

	updatedTime := rollInfo.UpdatedAt
	// Format the time to be more readable
	updatedTimeStr := updatedTime.Format(constants.DateTimeFormat)

	var message string
	if locationCnt > 0 {
		message = fmt.Sprintf(messageFmt, food.Name, food.Description, locationCnt, count, groups, updatedTimeStr, trigger.Name(), trigger.Username())
	} else {
		message = fmt.Sprintf(messageFmt, food.Name, food.Description, count, groups, updatedTimeStr, trigger.Name(), trigger.Username())
	}

	return message, locationCnt > 0
}

func RandomGroupsCommandReroll(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomGroups reroll button clicked by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CALLBACK)

	rollInfo, count, cb, err := services.RerollCommon(bot, ctx, constants.GROUP)
	if err != nil {
		return err
	}

	groups := rollInfo.GroupName
	groupsComma := strings.Join(strings.Split(groups, "\n"), ", ")

	// Send message to user with reroll button
	message, hasLoc := sendWithRerollButtonGroups(*rollInfo, ctx.EffectiveSender, count, true, groupsComma)
	messageOpts := utils.GenerateRerollKeysEdit(constants.GROUP, *rollInfo, hasLoc)

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: constants.RerollSuccess,
	})

	_, _, err = cb.Message.EditText(bot, message, messageOpts)

	return nil
}

func RandomGroupsCommandGroupList(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("RandomGroups command with group list called by " + ctx.EffectiveSender.Username())
	services.RunPreCommandScriptCustomType(ctx, constants.CONVERSATION)

	message := ctx.EffectiveMessage.Text
	groups := strings.Split(message, "\n")

	// trim whitespace
	for i, group := range groups {
		groups[i] = strings.Trim(group, " ")
	}

	message, rollInfo, hasLoc := searchAndReplyWithGroups(ctx, groups)
	// Send message to user with reroll button
	messageOpts := utils.GenerateRerollKeysSend(constants.GROUP, *rollInfo, hasLoc)

	err := utils.ReplyUserWithOpts(bot, ctx, message, messageOpts)
	if err != nil {
		log.Println("Error replying to user: " + err.Error())
	}

	return handlers.EndConversation()
}
