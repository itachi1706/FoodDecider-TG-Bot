package commands

import (
    "FoodDecider-TG-Bot/model"
    "FoodDecider-TG-Bot/repository"
    "FoodDecider-TG-Bot/utils"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/google/uuid"
    "log"
    "strconv"
    "strings"
)

func AddCoordinateCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
    log.Println("AddCoordinate command called by " + ctx.EffectiveSender.Username())

    userId := ctx.EffectiveSender.Id()
    // Make sure guy is an admin to run
    if utils.CheckIfAdmin(userId) == false {
        return utils.BasicReplyToUser(bot, ctx, "This command can only be ran by an administrator")
    }

    messageOpts := utils.GetArgumentsFromMessage(ctx)
    log.Printf("Message options: %v\n", messageOpts)
    if len(messageOpts) < 3 {
        return utils.BasicReplyToUser(bot, ctx, "Invalid Format\n\nFormat: /addcoordinate <food id> <latitude> <longitude> [name]")
    }

    foodId, err := uuid.Parse(messageOpts[0])
    if err != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid food id provided")
    }

    latitude, err1 := strconv.ParseFloat(messageOpts[1], 64)
    longitude, err2 := strconv.ParseFloat(messageOpts[2], 64)
    if err1 != nil || err2 != nil {
        return utils.BasicReplyToUser(bot, ctx, "Invalid latitude or longitude provided")
    }

    // Ensure latitude and longitude within range
    if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
        return utils.BasicReplyToUser(bot, ctx, "Invalid latitude or longitude provided. Out of range")
    }

    var friendlyName string
    if len(messageOpts) > 3 {
        // Has friendly name
        friendlyName = strings.Trim(strings.Join(messageOpts[3:], " "), " ")
    }
    log.Printf("Food ID: %v, Latitude: %v, Longitude: %v, Name: %v\n", foodId, latitude, longitude, friendlyName)
    db := utils.GetDbConnection()
    repo := repository.NewFoodsRepository(db)
    location := repo.GetFoodLocation(foodId, latitude, longitude)
    message := "An error has occurred. Please try again later"
    if location == nil {
        // New location
        log.Println("Creating new location for food " + foodId.String())
        location = &model.Locations{
            FoodID:    foodId,
            Name:      friendlyName,
            Latitude:  latitude,
            Longitude: longitude,
            CreatedBy: userId,
            UpdatedBy: userId,
            ID:        uuid.New(),
        }
        db.Create(&location)
        message = "Location added for food " + foodId.String()
    } else {
        location.Name = friendlyName
        location.UpdatedBy = userId
        message = "Location updated for food " + foodId.String()
        if location.Status != "A" {
            log.Println("Reactivating location for food " + foodId.String())
            location.Status = "A"
            message = "Location added for food " + foodId.String()
        }
        db.Save(&location)
    }

    return utils.BasicReplyToUser(bot, ctx, message)
}
