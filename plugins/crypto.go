package plugins

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

// Crypto holds data that is needed to pass around the crypto game functions.
type Crypto struct {
	words        []byte
	encoded      string
	lastEncoding string
	inPlay       bool
	roundStart   time.Time
	roundEnd     time.Time
	wasDecoded   bool
}

// NewCrypto returns a new struct of type crypto
func NewCrypto() *Crypto {
	return &Crypto{}
}

// Game launches a new crypto game or executes the check function.
func (c *Crypto) Game(req []string, msg *dg.MessageCreate) (string, error) {
	if len(req) < 2 {
		return fmt.Sprintf("Usage: `<prefix>crypto init` to start a game"), nil
	}
	if req[1] == "init" {
		if !c.inPlay {
			c.init()
		}
		return fmt.Sprintf("Game in progress. Current encoding:\n%s", c.encoded), nil
	}

	var userGuess string
	for idx, word := range req {
		log.Printf("%d: %s", idx, word)
		if idx == 1 {
			userGuess = fmt.Sprintf("%s", word)
		}
		if idx > 1 {
			userGuess = fmt.Sprintf("%s %s", userGuess, word)
		}
	}

	if isCorrect := c.checkGuess(userGuess); isCorrect {
		return fmt.Sprintf("%s won this round! Will you be next?", msg.Author.Username), nil
	}

	return fmt.Sprintf("%s sorry, that is incorrect :smirk:", msg.Author.Username), nil
}

// Init kicks off the crypto game
func (c *Crypto) init() (string, error) {
	words, err := c.callPaswdAPI()
	if err != nil {
		return "", err
	}

	return c.encode(words), err
}

func (c *Crypto) callPaswdAPI() ([]byte, error) {
	url := fmt.Sprintf("https://makemeapassword.ligos.net/api/v1/passphrase/plain")
	res, err := http.Get(url)
	if err != nil {
		return []byte(""), fmt.Errorf("crypto game failed to get data from provided url: %s", err)
	}

	secret, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return []byte(""), fmt.Errorf("crypto game failed to read res.Body from LoC: %s", err)
	}

	return secret, nil
}

func (c *Crypto) encode(src []byte) string {
	encodedStr := hex.EncodeToString(src)
	c.encoded = encodedStr

	c.roundStart = time.Now()
	return c.encoded
}

func (c *Crypto) checkGuess(guess string) bool {
	isCorrect := false

	if guess == string(c.words) {
		isCorrect = true
		c.wasDecoded = true
		c.roundEnd = time.Now()
	}

	return isCorrect
}
