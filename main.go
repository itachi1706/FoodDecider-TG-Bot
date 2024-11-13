package main

import (
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/joho/godotenv"
    "log"
    "os"
    "time"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file. Might not exist")
    }

    botToken := os.Getenv("BOT_TOKEN")
    if botToken == "" {
        panic("No BOT_TOKEN provided")
    }

    bot, err := gotgbot.NewBot(botToken, nil)
    if err != nil {
        panic("bot creation failed: " + err.Error())
    }

    // Create dispatching commands
    dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
        // If an error is returned by a handler, log it and continue going.
        Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
            log.Println("an error occurred while handling update:", err.Error())
            return ext.DispatcherActionNoop
        },
        MaxRoutines: ext.DefaultMaxRoutines,
    })
    updater := ext.NewUpdater(dispatcher, nil)

    // TODO: Add commands here
    InitCommands(dispatcher)
    commandInstalled, err := bot.GetMyCommands(nil)
    if err != nil {
        panic("failed to get commands: " + err.Error())
    }

    err = updater.StartPolling(bot, &ext.PollingOpts{
        DropPendingUpdates: true,
        GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
            Timeout: 9,
            RequestOpts: &gotgbot.RequestOpts{
                Timeout: time.Second * 10,
            },
        },
    })

    if err != nil {
        panic("polling failed: " + err.Error())
    }

    log.Printf("Bot %s started\n", bot.User.Username)
    log.Printf("Bot Info: %#v\n", bot.User)
    log.Printf("Commands available: %#v\n", commandInstalled)

    updater.Idle() // Idle, to keep updates coming in, and avoid bot stopping.
}
