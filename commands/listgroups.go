package commands

import (
    "FoodDecider-TG-Bot/model"
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/services"
    "FoodDecider-TG-Bot/utils"
    "fmt"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/google/uuid"
    "log"
)

func ListGroupsCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListGroups command called by " + ctx.EffectiveSender.Username())

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 1 {
        return utils.BasicReplyToUser(bot, ctx, "Food ID required\n\nFormat: /listgroups <food id>")
    }

    foodId, err := uuid.Parse(messageOpts[0])
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
    }

    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    // Get first 5 food results with status A
    foodGroups := repo.FindAllGroupsForFoodPaginated(foodId, 5, 0)
    food := repo.FindFoodById(foodId)
    message := populateListFoodGroupsMessage(foodGroups, food)

    return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend("group-list+"+foodId.String()+"+", 0, true, true))
}

func populateListFoodGroupsMessage(groups []model.FoodGroups, food *model.Food) string {
    foodName := "Unknown Food"
    if food != nil {
        foodName = food.Name
    }

    message := "No groups found for " + foodName
    if len(groups) > 0 {
        message = "Groupings for " + foodName + ":\n\n"
        for _, group := range groups {
            message += fmt.Sprintf("%s\n", group.Name)
        }
    }
    return message
}

func ListGroupsCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListGroups previous button clicked by " + ctx.EffectiveSender.Username())

    cb := ctx.Update.CallbackQuery
    log.Println("Callback data: " + cb.Data)

    foodId, pageCnt, err := services.ParseFoodParameters(cb.Data)
    if err != nil {
        _, _ = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
            Text: "An error occurred. Please try again later",
        })
        return fmt.Errorf("failed to parse data: %w", err)
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
        return fmt.Errorf("failed to answer callback: %w", err)
    }

    if !cont {
        return nil // End here
    }

    // Get previous 5 food results with status A
    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    foodGroups := repo.FindAllGroupsForFoodPaginated(foodId, 5, pageCnt)
    food := repo.FindFoodById(foodId)
    message := populateListFoodGroupsMessage(foodGroups, food)
    _, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit("group-list+"+foodId.String()+"+", pageCnt, true, true))

    return nil
}

func ListGroupsCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListGroups next button clicked by " + ctx.EffectiveSender.Username())

    cb := ctx.Update.CallbackQuery
    log.Println("Callback data: " + cb.Data)

    err, foodId, pageCnt := services.HandleFoodNextCommands(bot, cb)
    if err != nil || foodId == nil || pageCnt == nil {
        return err // End here
    }

    // Get next 5 food results with status A
    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    foodGroups := repo.FindAllGroupsForFoodPaginated(*foodId, 5, *pageCnt)
    food := repo.FindFoodById(*foodId)
    message := populateListFoodGroupsMessage(foodGroups, food)
    _, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit("group-list+"+foodId.String()+"+", *pageCnt, true, true))

    return nil
}
