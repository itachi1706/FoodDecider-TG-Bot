package main

import (
    "FoodDecider-TG-Bot/commands"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
    "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
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
    dispatcher.AddHandler(handlers.NewCommand("removegroup", commands.RemoveGroupCommand))
    dispatcher.AddHandler(handlers.NewConversation(
        []ext.Handler{handlers.NewCommand("renamegroup", commands.RenameGroupCommand)},
        map[string][]ext.Handler{
            commands.NewGroupName: {handlers.NewMessage(noCommands, commands.RenameGroupCommandNewName)},
        },
        &handlers.ConversationOpts{
            Exits:        []ext.Handler{handlers.NewCommand("cancel", commands.CancelCommand)},
            StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
            AllowReEntry: true,
        }))
    dispatcher.AddHandler(handlers.NewCommand("listgroups", commands.ListGroupsCommand))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("previous-group-list"), commands.ListGroupsCommandPrev))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("next-group-list"), commands.ListGroupsCommandNext))

    dispatcher.AddHandler(handlers.NewCommand("addcoordinate", commands.AddCoordinateCommand))
    dispatcher.AddHandler(handlers.NewCommand("listcoordinate", commands.ListCoordinatesCommand))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("previous-coordinate-list"), commands.ListCoordinatesCommandPrev))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("next-coordinate-list"), commands.ListCoordinatesCommandNext))

    log.Println("Commands initialized")
}

func noCommands(msg *gotgbot.Message) bool {
    return message.Text(msg) && !message.Command(msg)
}
