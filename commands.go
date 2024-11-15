package main

import (
    "FoodDecider-TG-Bot/commands"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
    "log"
)

func InitCommands(dispatcher *ext.Dispatcher) {
    log.Println("Initializing commands")

    dispatcher.AddHandler(handlers.NewCommand("start", commands.StartCommand))
    dispatcher.AddHandler(handlers.NewCommand("debuginfo", commands.DebugInfoCommand))
    dispatcher.AddHandler(handlers.NewCommand("userinfo", commands.UserInfoCommand))

    dispatcher.AddHandler(handlers.NewCommand("addadmin", commands.AddAdminCommand))
    dispatcher.AddHandler(handlers.NewCommand("deladmin", commands.DelAdminCommand))
    dispatcher.AddHandler(handlers.NewCommand("listadmins", commands.ListAdminsCommand))

    dispatcher.AddHandler(handlers.NewCommand("addfood", commands.AddFoodCommand))
    dispatcher.AddHandler(handlers.NewCommand("updatefood", commands.UpdateFoodCommand))
    dispatcher.AddHandler(handlers.NewCommand("delfood", commands.DelFoodCommand))
    dispatcher.AddHandler(handlers.NewCommand("listfoods", commands.ListFoodsCommand))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("previous-food-list-"), commands.ListFoodsCommandPrev))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("next-food-list-"), commands.ListFoodsCommandNext))

    dispatcher.AddHandler(handlers.NewCommand("addgroup", commands.AddGroupCommand))

    log.Println("Commands initialized")
}
