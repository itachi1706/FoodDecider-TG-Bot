package commands

import (
    "FoodDecider-TG-Bot/model"
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/utils"
    "fmt"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "log"
    "strconv"
    "strings"
)

func ListFoodsCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListFoods command called by " + ctx.EffectiveSender.Username())

    db := utils.GetDbConnection()
    // Get first 5 food results with status A
    foods := repository.FindAllActiveFoodPaginated(db, 5, 0)
    message := populateMessage(foods)

    return utils.ReplyUserWithOpts(bot, ctx, message, utils.GeneratePageKeysSend("food-list", 0, true, true))
}

func populateMessage(foods []model.Food) string {
    message := "No foods found"
    if len(foods) > 0 {
        message = "List of foods:\n\n"
        for _, food := range foods {
            desc := food.Description
            if desc == "" {
                desc = "No description provided"
            }
            message += fmt.Sprintf("ID: %s\nName: %s\nDescription: %s\n\n", food.ID.String(), food.Name, desc)
        }
    }
    return message

}

func ListFoodsCommandPrev(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListFoods previous button clicked by " + ctx.EffectiveSender.Username())

    cb := ctx.Update.CallbackQuery
    log.Println("Callback data: " + cb.Data)

    pageCnt, err := strconv.Atoi(strings.Replace(cb.Data, "previous-food-list-", "", -1))
    if err != nil {
        log.Printf("failed to convert page number: %w. Default 0\n", err)
        pageCnt = 0
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
    foods := repository.FindAllActiveFoodPaginated(db, 5, pageCnt)

    message := populateMessage(foods)
    _, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit("food-list", pageCnt, true, true))

    return nil
}

func ListFoodsCommandNext(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("ListFoods next button clicked by " + ctx.EffectiveSender.Username())

    cb := ctx.Update.CallbackQuery
    log.Println("Callback data: " + cb.Data)

    pageCnt, err := strconv.Atoi(strings.Replace(cb.Data, "next-food-list-", "", -1))
    if err != nil {
        log.Printf("failed to convert page number: %w. Default 0\n", err)
        pageCnt = 0
    }

    db := utils.GetDbConnection()
    // Get total number of food and find number of possible pages (including partial)
    count := repository.GetFoodCount(db)
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
    foods := repository.FindAllActiveFoodPaginated(db, 5, pageCnt)

    message := populateMessage(foods)
    _, _, err = cb.Message.EditText(bot, message, utils.GeneratePageKeysEdit("food-list", pageCnt, true, true))

    return nil
}
