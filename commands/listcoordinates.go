package commands

import (
    "FoodDecider-TG-Bot/model"
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

func ListCoordinatesCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListCoordinates command called by " + ctx.EffectiveSender.Username())

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 1 {
        return utils.BasicReplyToUser(bot, ctx, "Food ID required\n\nFormat: /listcoordinates <food id>")
    }

    foodId, err := uuid.Parse(messageOpts[0])
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
    }

    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    // Get first 5 food results with status A
    foodGroups := repo.FindAllLocationsForFoodPaginated(foodId, 5, 0)
    food := repo.FindFoodById(foodId)
    message := populateListFoodLocationsMessage(foodGroups, food)

    return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend("coordinate-list+"+foodId.String()+"+", 0, true, true))
}

func populateListFoodLocationsMessage(groups []model.Locations, food *model.Food) string {
    foodName := "Unknown Food"
    if food != nil {
        foodName = food.Name
    }

    message := "No locations found for " + foodName
    if len(groups) > 0 {
        message = "Locations for " + foodName + ":\n\n"
        for _, group := range groups {
            name := group.Name
            if name == "" {
                name = "No name defined"
            }
            message += fmt.Sprintf("ID: %s\nLocation: %v, %v\nName: %s\n\n", group.ID, group.Latitude, group.Longitude, name)
        }
    }
    return message
}

func parseFoodLocationData(data string) (uuid.UUID, int, error) {
    // data: <prefix>-coordinate-list+<food_id>+-<page>
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

func ListCoordinatesCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListGroups previous button clicked by " + ctx.EffectiveSender.Username())

    cb := ctx.Update.CallbackQuery
    log.Println("Callback data: " + cb.Data)

    foodId, pageCnt, err := parseFoodLocationData(cb.Data)
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
    foodLocations := repo.FindAllLocationsForFoodPaginated(foodId, 5, pageCnt)
    food := repo.FindFoodById(foodId)
    message := populateListFoodLocationsMessage(foodLocations, food)
    _, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit("coordinate-list+"+foodId.String()+"+", pageCnt, true, true))

    return nil
}

func ListCoordinatesCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListGroups next button clicked by " + ctx.EffectiveSender.Username())

    cb := ctx.Update.CallbackQuery
    log.Println("Callback data: " + cb.Data)

    foodId, pageCnt, err := parseFoodLocationData(cb.Data)
    if err != nil {
        _, _ = cb.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
            Text: "An error occurred. Please try again later",
        })
        return fmt.Errorf("failed to parse data: %w", err)
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
        return fmt.Errorf("failed to answer callback: %w", err)
    }

    if !cont {
        return nil // End here
    }

    // Get next 5 food results with status A
    foodLocations := repo.FindAllLocationsForFoodPaginated(foodId, 5, pageCnt)
    food := repo.FindFoodById(foodId)
    message := populateListFoodLocationsMessage(foodLocations, food)
    _, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit("coordinate-list+"+foodId.String()+"+", pageCnt, true, true))

    return nil
}
