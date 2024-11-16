package services

import (
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

func ParseFoodParameters(data string) (uuid.UUID, int, error) {
	// split by + and -
	splitData := strings.Split(data, "+")
	if len(splitData) < 3 {
		return uuid.Nil, 0, fmt.Errorf("invalid data provided")
	}

	foodId, err := uuid.Parse(splitData[1])
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("invalid food id provided")
	}

	// Remove the first "-" from splitData[2]
	splitData[2] = strings.Replace(splitData[2], "-", "", 1)
	pageCnt, err := strconv.Atoi(splitData[2])
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("invalid page number provided")
	}

	return foodId, pageCnt, nil
}

func HandleFoodPrevCommands(bot *gotgbot.Bot, cb *gotgbot.CallbackQuery) (error, *uuid.UUID, *int) {
	foodId, pageCnt, err := ParseFoodParameters(cb.Data)
	if err != nil {
		_, _ = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: "An error occurred. Please try again later",
		})
		return fmt.Errorf("failed to parse data: %w", err), nil, nil
	}

	answerMsg := "An error occurred. Please try again later"
	cont := true
	if pageCnt <= 0 {
		// First page
		answerMsg = "You are already on the first page"
		cont = false
	} else {
		answerMsg = "Going to previous page"
		pageCnt--
	}

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: answerMsg,
	})

	if err != nil {
		return fmt.Errorf("failed to answer callback: %w", err), nil, nil
	}

	if !cont {
		return nil, nil, nil // End here
	}

	return nil, &foodId, &pageCnt
}

func HandleFoodNextCommands(bot *gotgbot.Bot, cb *gotgbot.CallbackQuery) (error, *uuid.UUID, *int) {
	foodId, pageCnt, err := ParseFoodParameters(cb.Data)
	if err != nil {
		_, _ = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: "An error occurred. Please try again later",
		})
		return fmt.Errorf("failed to parse data: %w", err), nil, nil
	}

	db := utils.GetDbConnection()
	repo := repository.NewFoodsRepository(db)
	// Get total number of food and find number of possible pages (including partial)
	count := repo.GetFoodGroupForFoodCount(foodId)
	totalPages := count / 5
	modulo := count % 5
	if modulo > 0 {
		totalPages++
	}

	// pagecnt to int64
	pageCnt64 := int64(pageCnt)

	answerMsg := "An error occurred. Please try again later"
	cont := true
	if pageCnt64 >= totalPages-1 {
		// last page
		answerMsg = "You are already on the last page"
		cont = false
	} else {
		answerMsg = "Going to next page"
		pageCnt++
	}

	_, err = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text: answerMsg,
	})

	if err != nil {
		return fmt.Errorf("failed to answer callback: %w", err), nil, nil
	}

	if !cont {
		return nil, nil, nil // End here
	}

	return nil, &foodId, &pageCnt
}

func GroupHandlingParameter(bot *gotgbot.Bot, ctx *ext.Context, format string) (*int64, *uuid.UUID, *string, error) {
	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return nil, nil, nil, utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)
	if len(messageOpts) < 2 {
		return nil, nil, nil, utils.BasicReplyToUser(bot, ctx, "Invalid Format\n\nFormat: "+format)
	}

	foodId, err := uuid.Parse(messageOpts[0])
	if err != nil {
		return nil, nil, nil, utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
	}

	groupName := strings.Trim(strings.Join(messageOpts[1:], " "), " ")

	return &userId, &foodId, &groupName, nil
}

func FoodValidationParameterChecks(bot *gotgbot.Bot, ctx *ext.Context, argLen int, errorMsg string) (*int64, *uuid.UUID, []string, error) {
	userId := ctx.EffectiveSender.Id()
	// Make sure guy is an admin to run
	if utils.CheckIfAdmin(userId) == false {
		return nil, nil, nil, utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
	}

	messageOpts := utils.GetArgumentsFromMessage(ctx)
	log.Printf("Message options: %v\n", messageOpts)
	if len(messageOpts) < argLen {
		return nil, nil, messageOpts, utils.BasicReplyToUser(bot, ctx, errorMsg)
	}

	foodId, err := uuid.Parse(messageOpts[0])
	if err != nil {
		return nil, nil, messageOpts, utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
	}

	return &userId, &foodId, messageOpts, nil
}
