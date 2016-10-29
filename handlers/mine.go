package handlers

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/SophisticaSean/meme_coin/interaction"
	"github.com/jmoiron/sqlx"
)

// MineResponse is a struct for possible events to the !mine action
type MineResponse struct {
	amount   int
	response string
	chance   int
}

// GenerateResponseList picks a response out of our possible mine responses
// and returns it
func GenerateResponseList() []MineResponse {
	mineResponses := []MineResponse{
		MineResponse{
			amount:   100,
			response: " mined for a while and managed to scrounge up $AMOUNT$ dusty memes",
			chance:   75,
		},
		MineResponse{
			amount:   300,
			response: " mined for a bit and found an uncommon pepe worth $AMOUNT$ memes!",
			chance:   60,
		},
		MineResponse{
			amount:   1000,
			response: " fell down a meme-shaft and found a very dank rare pepe worth $AMOUNT$ memes!",
			chance:   15,
		},
		MineResponse{
			amount:   5000,
			response: " wandered in the meme mine for what seems like forever, eventually stumbling upon a vintage 1980's pepe worth $AMOUNT$ memes!",
			chance:   5,
		},
		MineResponse{
			amount:   25000,
			response: "'s meme mining has made the Maymay gods happy, they rewarded them with a MLG-shiny-holofoil-dankasheck Pepe Diamond worth $AMOUNT$ memes!",
			chance:   1,
		}}

	responseList := []MineResponse{}

	for _, response := range mineResponses {
		counter := response.chance
		for counter > 0 {
			responseList = append(responseList, response)
			counter--
		}
	}
	return responseList
}

// Mine processes a mining interaction for a player
func Mine(s interaction.Session, m *interaction.MessageCreate, responseList []MineResponse, db *sqlx.DB) {
	author := UserGet(m.Author, db)
	difference := time.Now().Sub(author.MineTime)
	timeLimit := 5
	channel, err := s.Channel(m.ChannelID)
	_, originalProduction, _ := ProductionSum(m.Author, db)
	productionMultiplier := int(math.Floor(float64(originalProduction) / float64(rand.Intn(160)+40)))
	if productionMultiplier < 1 {
		productionMultiplier = 1
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if channel.IsPrivate {
		_, _ = s.ChannelMessageSend(m.ChannelID, "you think you're slick, eh? gotta mine in a public room, bro.")
		return
	}

	// check to make sure user is not trying to mine before timeLimit has passed
	if difference.Minutes() < float64(timeLimit) {
		waitTime := strconv.Itoa(int(math.Ceil((float64(timeLimit) - difference.Minutes()))))
		_, _ = s.ChannelMessageSend(m.ChannelID, author.Username+" is too tired to mine, they must rest their meme muscles for "+waitTime+" more minute(s)")
		return
	}

	// pick a response out of the responses in responseList
	pickedIndex := rand.Intn(len(responseList))
	mineResponse := responseList[pickedIndex]
	amount := (mineResponse.amount * productionMultiplier)
	MoneyAdd(&author, amount, "mined", db)
	amountRegex := regexp.MustCompile(`\$AMOUNT\$`)
	response := amountRegex.ReplaceAllString(mineResponse.response, strconv.Itoa(amount))
	_, _ = s.ChannelMessageSend(m.ChannelID, author.Username+response)
	fmt.Println(author.Username + " mined " + strconv.Itoa(amount))
	return
}
