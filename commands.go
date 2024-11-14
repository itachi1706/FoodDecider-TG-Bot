package main

import (
    "FoodDecider-TG-Bot/commands"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
    "log"
)

func InitCommands(dispatcher *ext.Dispatcher) {
    log.Println("Initializing commands")

    dispatcher.AddHandler(handlers.NewCommand("start", commands.StartCommand))
    dispatcher.AddHandler(handlers.NewCommand("debuginfo", commands.DebugInfoCommand))

    dispatcher.AddHandler(handlers.NewCommand("addadmin", commands.AddAdminCommand))

    log.Println("Commands initialized")
}
