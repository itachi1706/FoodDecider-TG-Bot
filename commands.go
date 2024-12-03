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
	dispatcher.AddHandler(handlers.NewCommand("delcoordinate", commands.DelCoordinateCommand))
	dispatcher.AddHandler(handlers.NewCommand("listcoordinates", commands.ListCoordinatesCommand))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list-coordinates"), commands.ListCoordinatesCommandTrigger))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("previous-coordinate-list"), commands.ListCoordinatesCommandPrev))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("next-coordinate-list"), commands.ListCoordinatesCommandNext))
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("addlocation", commands.AddLocationCommand)},
		map[string][]ext.Handler{
			commands.AddFoodLocation: {handlers.NewMessage(locationPin, commands.AddLocationCommandLocationPin)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", commands.CancelCommand)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		}))

	dispatcher.AddHandler(handlers.NewCommand("randomfood", commands.RandomFoodCommand))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("reroll-GENERAL"), commands.RandomFoodCommandReroll))
	dispatcher.AddHandler(handlers.NewCommand("lastdecision", commands.LastDecisionCommand))
	dispatcher.AddHandler(handlers.NewCommand("decisionhistory", commands.DecisionHistoryCommand))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("previous-decision-history-"), commands.DecisionHistoryCommandPrev))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("next-decision-history-"), commands.DecisionHistoryCommandNext))
	dispatcher.AddHandler(handlers.NewCommand("rollhistory", commands.RollHistoryCommand))

	// Random group
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("randomgroups", commands.RandomGroupsCommand)},
		map[string][]ext.Handler{
			commands.RandomGroupSpecify: {handlers.NewMessage(noCommands, commands.RandomGroupsCommandGroupList)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", commands.CancelCommand)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("reroll-GROUP"), commands.RandomGroupsCommandReroll))

	// Random location
	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("randomlocation", commands.RandomLocationCommand),
			handlers.NewCommand("randomnearme", commands.RandomLocationCommand)},
		map[string][]ext.Handler{
			commands.SelectLocationRandom: {handlers.NewMessage(locationPin, commands.RandomLocationCommandLocationPin)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", commands.CancelCommand)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("reroll-LOCATION"), commands.RandomLocationCommandReroll))

	log.Println("Commands initialized")
}

func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

func locationPin(msg *gotgbot.Message) bool {
	return message.Location(msg)
}
