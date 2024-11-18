package services

import (
	"FoodDecider-TG-Bot/model"
	"FoodDecider-TG-Bot/repository"
	"FoodDecider-TG-Bot/utils"
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"log"
)

func LogUserFound(ctx *ext.Context) {
	// Log user found
	user := ctx.EffectiveSender
	log.Printf("User found: %#v", user)

	db := utils.GetDbConnection()
	repo := repository.NewUserRepository(db)

	// Find user if exists
	userObj := repo.FindUser(user.Id())

	// Convert user object to string
	userString, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to marshal user object: %s. Skipping", err)
	}

	if userObj != nil {
		log.Printf("User found: %#v", user)
		// Check if username and full name is identical
		if userObj.Username != user.Username() || userObj.FullName != user.Name() {
			log.Printf("User data is not identical. Updating user data.")
			userObj.Username = user.Username()
			userObj.FullName = user.Name()
			userObj.RawData = string(userString)
			db.Save(userObj)

			// Create history object too
			historyObj := createHistoryObject(userObj)

			// Save history
			db.Save(historyObj)
		}
		// If user data is identical, do nothing
	} else {
		log.Printf("User not found. Creating new user.")
		userObj = &model.Users{
			TelegramID: user.Id(),
			Username:   user.Username(),
			FullName:   user.Name(),
			CreatedBy:  user.Id(),
			ID:         uuid.New(),
			RawData:    string(userString),
		}

		// Create history object too
		historyObj := createHistoryObject(userObj)

		// Save user and history
		db.Save(userObj)
		db.Save(historyObj)
	}
}

func createHistoryObject(userObj *model.Users) *model.PastHistory {
	return &model.PastHistory{
		UserID:   userObj.ID,
		Username: userObj.Username,
		FullName: userObj.FullName,
	}
}
