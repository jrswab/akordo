package plugins

import (
	"crypto/md5"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

// Crypto holds data that is needed to pass around the crypto game functions.
type Crypto struct {
	words        []byte
	encoded      string
	lastEncoding int
	inPlay       bool
	//roundStart   time.Time
	roundEnd   time.Time
	wasDecoded bool
	waitTime   time.Duration
}

// NewCrypto returns a new struct of type crypto
func NewCrypto() *Crypto {
	return &Crypto{waitTime: 2}
}

// Game launches a new crypto game or executes the check function.
func (c *Crypto) Game(req []string, msg *dg.MessageCreate) (string, error) {
	if len(req) < 2 {
		return fmt.Sprintf("Usage: `<prefix>crypto init` to start a game"), nil
	}

	if req[1] == "init" {
		gameTimeout := c.roundEnd.Add(c.waitTime * time.Minute).Unix()
		initTime := time.Now().Unix()

		if c.inPlay {
			return fmt.Sprintf("Mining in progress...\nCurrent encoding:\n%s", c.encoded), nil
		}

		if gameTimeout >= initTime {
			return fmt.Sprintf("Please wait %d minutes to open a new mine.", c.waitTime), nil
		}

		// Start a new crypto game
		err := c.init()
		if err != nil {
			return "", fmt.Errorf("cryto game error: %s", err)
		}
		return fmt.Sprintf("A new mine has opened! \n%s", c.encoded), nil
	}

	var userGuess string
	for idx, word := range req {
		if idx == 1 {
			userGuess = fmt.Sprintf("%s", word)
		}
		if idx > 1 {
			userGuess = fmt.Sprintf("%s %s", userGuess, word)
		}
	}

	if isCorrect := c.checkGuess(userGuess); isCorrect {
		return fmt.Sprintf("%s won this round!", msg.Author.Username), nil
	}

	return fmt.Sprintf("%s sorry, that is incorrect :smirk:", msg.Author.Username), nil
}

// Init kicks off the crypto game
func (c *Crypto) init() error {
	url := fmt.Sprintf("https://makemeapassword.ligos.net/api/v1/passphrase/plain")
	if err := c.callPasswdAPI(url); err != nil {
		return err
	}

	c.encode()
	return nil
}

func (c *Crypto) callPasswdAPI(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("crypto game failed to get data from provided url: %s", err)
	}

	words, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("crypto game failed to read res.Body from LoC: %s", err)
	}

	// Remove new lines form received data
	c.words = []byte(strings.Trim(string(words), "\r\n"))

	return nil
}

func (c *Crypto) encode() {
	var n int
	for c.lastEncoding < 10 {
		rand.Seed(time.Now().UnixNano())
		n = rand.Intn(4)
		if n != c.lastEncoding {
			break
		}
	}

	switch n {
	case 0:
		c.encoded = hex.EncodeToString([]byte(c.words))
		c.lastEncoding = 0
	case 1:
		c.encoded = base32.StdEncoding.EncodeToString([]byte(c.words))
		c.lastEncoding = 1
	case 2:
		c.encoded = base64.StdEncoding.EncodeToString([]byte(c.words))
		c.lastEncoding = 2
	case 3:
		c.encoded = ""
		for _, bits := range []byte(c.words) {
			c.encoded = fmt.Sprintf("%s %d", c.encoded, bits)
		}
		c.lastEncoding = 3
	case 4:
		c.encoded = fmt.Sprintf("%x", md5.Sum([]byte(c.words)))
		c.lastEncoding = 4
	}
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
