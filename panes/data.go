/*
 *	PROGRAM		: data.go
 *	DESCRIPTION		:
 *
 *		This program handles field definitions.
 *
 *	PARAMETERS		:
  *
 *	RETURNS			:
 *
*/
package panes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

/*
 *	  These values are stored in the config.json file
 */

// queues tab
var Queue string         // server message queue
var Queuepassword string // server message queue password

// authenticate  tab
var Password string     // encrypt file password
var Passwordhash string // hash value of password
//
var PasswordMinimumSize int         // set minimum password size
var PasswordMustContainNumber bool  // password must contain number
var PasswordMustContainLetter bool  // password must contain letter
var PasswordMustContainSpecial bool // password must contain special character

// Server tab
var Server string // server url

var UseJetstream bool   // if set to true uses jetstream protocol otherwise regular pub/sub
var UseTLS string       // use TLS to Authenticate  else use userid /password
var Caroot string       // for UseTLS = true CAROOT certificate for server authentication
var UserID string       // for UseTLS = false
var UserPassword string // for UseTLS = false

/*
 *	  These constants are set to establish a password schema for Local File Encryption and Queue password
 */

const Cipherkey = "asuperstrong32bitpasswordgohere!" // 32 byte string  for hash value of cipher key to decrypt json fields modify this field for your ntwork
const PasswordDefault = "123456"                     // default password shipped with app
/*
 *	  Confignats is used to hold config.json fields
 */
type Confignats struct {
	//jlkj
	Jserver                     string // server url nats://333.333.333.333:port
	Jqueue                      string // queue created for deployment
	Jqueuepassword              string // queue password created for deployment
	Jpasswordminimumsize        int    // set minimum password size
	Jpasswordmustcontainnumber  bool   // password must contain number
	Jpasswordmustcontainletter  bool   // password must contain letter
	Jpasswordmustcontainspecial bool   // password must contain special character
	Jusejetstream               bool   // if set to true uses jetstream protocol otherwise regular pub/sub
	Jusetls                     string // use TLS to Authenticate  else use userid /password
	Jcaroot                     string // for UseTLS = true CAROOT certificate for server authentication
	Juserid                     string // for UseTLS = false
	Juserpassword               string // for UseTLS = false
}

// Pane defines the data structure
type MyPane struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Panes defines the metadata
	MyPanes = map[string]MyPane{
		"password": {"Password", "", passwordScreen, true},
		"logon":    {"Logon", "", logonScreen, true},
		"messages": {"Messages", "", messagesScreen, true},
	}

	// PanesIndex  defines how our panes should be laid out in the index tree
	MyPanesIndex = map[string][]string{
		"": {"password", "logon", "messages"},
	}
)

var MM []string
var MMSize = 200

func AddMessage(message string) {

	MM = append(MM, message+"\n")
	if len(MM) > MMSize {
		MM = MM[:len(MM)-100]
	}

}

/*
 *	FUNCTION		: parseURL
 *	DESCRIPTION		:
 *		This function takes a string and parses it for validity
 *
 *	PARAMETERS		:
 *		urlStr  	: String of url to parse
 *
 *	RETURNS			:
 *		string link	: URL for linking
 */
func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

/*
 *	FUNCTION		: dumpglobals
 *	DESCRIPTION		:
 *		Helper function to debug
 *
 *	PARAMETERS		:
 *		from    	: Name of calling location
 *
 *	PRINTS			:
 *		        	: Dumps contents to log
 */
func dumpglobals(from string) {
	log.Println(from, Password)
	log.Println(from, Passwordhash)
	log.Println(from, Server)
	log.Println(from, Caroot)
	log.Println(from, Queue)
	log.Println(from, Queuepassword)

}

/*
 *	FUNCTION		: loadconfig
 *	DESCRIPTION		:
 *		This function loads the json config file
 *
 *	PARAMETERS		:
 *		        	: None
 *
 *	RETURNS			:
 *		map     	: Map of interface with encrypted field values
 */
func loadconfig() map[string]interface{} {
	log.Println("loadconfig Password", Password)
	data := map[string]interface{}{
		"server":        string(Server),
		"caroot":        string(Caroot),
		"queue":         string(Queue),
		"queuepassword": string(Queuepassword),
	}
	log.Println("loadconfig data", data)
	return data
}

/*
 *	FUNCTION		: loadhash
 *	DESCRIPTION		:
 *		This function loads the json hash file
 *
 *	PARAMETERS		:
 *		        	: None
 *
 *	RETURNS			:
 *		map     	: Map of interface with hash file for password checking
 */
func loadhash() map[string]interface{} {
	log.Println("loadhash Password hash", Passwordhash)
	log.Println("loadhash Cipherkey", Cipherkey)
	data := map[string]interface{}{
		"passwordhash": string(Passwordhash),
	}
	log.Println("loadhash data", data)
	return data
}

/*
 *	FUNCTION		: myjson
 *	DESCRIPTION		:
 *		This function handles file actions for config.json to load memory
 *
 *	PARAMETERS		        :
 *		action string   	: CREATE, LOAD or SAVE encrypted fields
 *
 *	RETURNS			:
 *		         	: None
 */
func myjson(action string) {

	if action == "CREATE" {
		log.Println("MyJson Create ", Password)
		Server = string("None")
		Caroot = string("None")
		Queue = string("None")
		Queuepassword = string("None")
		Encmessage = string("None")
		configfile, configfileerr := os.Create("config.json")
		if configfileerr == nil {
			enc := json.NewEncoder(configfile)

			MyCrypt("ENCRYPT")
			enc.Encode(loadconfig())
		}
		configfile.Close()
	}
	if action == "LOAD" {

		var c map[string]interface{}

		jf, errf := os.Open("config.json")
		if errf != nil {
			log.Println("LOAD Error file", errf)
		}
		jc, je := ioutil.ReadAll(jf)

		if je != nil {
			log.Println("LOAD Error read all", je)
		}
		jf.Close()

		json.Unmarshal([]byte(jc), &c)
		for k, v := range c {
			// decrypt all fields ater loading
			//fmt.Println(k, "=>", v)

			if k == "server" {

				Server = v.(string)
			}
			if k == "caroot" {
				Caroot = v.(string)
			}
			if k == "queue" {
				Queue = v.(string)
			}
			if k == "queuepassword" {
				Queuepassword = v.(string)
			}
		}
		MyCrypt("DECRYPT")

	}
	if action == "SAVE" {
		e := os.Remove("config.json")
		if e != nil {
			log.Fatal(e)
		}
		sc, se := os.Create("config.json")

		if se == nil {
			enc := json.NewEncoder(sc)

			MyCrypt("ENCRYPT")
			enc.Encode(loadconfig())

		}

		sc.Close()
	}
}

/*
 *	FUNCTION		: MyCrypt Public function to be used by message encryption/decryption
 *	DESCRIPTION		:
 *		This function handles fiedd encryption/decryption of memory
 *
 *	PARAMETERS		        :
 *		action string   	: ENCRYPT or DECRYPT
 *
 *	RETURNS			:
 *		         	: None
 */
func MyCrypt(action string) {
	if action == "ENCRYPTNOTNOW" {
		var newvalue, _ = encrypt([]byte(Cipherkey), Server)
		Server = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Caroot)
		Caroot = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Queue)
		Queue = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Queuepassword)
		Queuepassword = newvalue

	}
	if action == "DECRYPTNOTNOW" {
		var newvalue, _ = decrypt([]byte(Server), Cipherkey)
		Server = newvalue
		newvalue, _ = decrypt([]byte(Caroot), Cipherkey)
		Caroot = newvalue
		newvalue, _ = decrypt([]byte(Queue), Cipherkey)
		Queue = newvalue
		newvalue, _ = decrypt([]byte(Queuepassword), Cipherkey)
		Queuepassword = newvalue

	}
}

/*
 *	FUNCTION		: myhash
 *	DESCRIPTION		:
 *		This function save the password hash to config.hash
 *
 *	PARAMETERS		        :
 *		action string   	: CREATE, LOAD or SAVE
 *      hash string         : Value of password hash
 *
 *	RETURNS			:
 *		         	: None
 */
func MyHash(action string, hash string) {

	if action == "CREATE" {
		log.Println("create Hash", hash)

		confighash, confighasherr := os.Create("config.hash")
		if confighasherr == nil {
			enc := json.NewEncoder(confighash)
			//cipherKey := []byte("!99099jjhhnniikjkjilhh7dDDDkillp") //32 bit key for AES-256

			log.Println("myhash save config", loadhash())
			enc.Encode(loadhash())
		}
		confighash.Close()
	}
	if action == "LOAD" {

		var c map[string]interface{}
		jf, errf := os.Open("config.hash")
		if errf != nil {
			log.Println("LOAD Hash Error file", errf)
		}
		jc, je := ioutil.ReadAll(jf)

		if je != nil {
			log.Println("LOAD Hash Error read all", je)
		}
		jf.Close()

		json.Unmarshal([]byte(jc), &c)
		for k, v := range c {
			fmt.Println(k, "=>", v)
			if k == "passwordhash" {
				Passwordhash = v.(string)
			}
		}

	}
	if action == "SAVE" {
		e := os.Remove("config.hash")
		if e != nil {
			log.Fatal(e)
		}
		sc, se := os.Create("config.hash")

		if se == nil {
			enc := json.NewEncoder(sc)
			log.Println("myhash save", loadhash())
			enc.Encode(loadhash())

		}

		sc.Close()
	}
}

/*
 *	FUNCTION		: NATSConnect
 *	DESCRIPTION		:
 *		This function connects to the nats server and populates mm in data using a go thread
 *
 *	PARAMETERS		:
 *
 *	RETURNS		!	:
 */
func NATSConnect() {

	//nc, err := nats.Connect(Server, nats.RootCAs(strings.Trim(Caroot,"\r")))
	// use jetstream
	// save caroot in json to file
	if UseJetstream == true {
		nc, err := nats.Connect(Server, nats.RootCAs("./ca-nats.pem"))
		// delete caroot from file system after connect
		if err == nil {

			nc.Subscribe(Queue+">", func(msg *nats.Msg) {

				log.Println("mymessage   - ", msg.Subject)
				mymessage := msg.Subject + string(msg.Data)

				AddMessage(mymessage)

			})
		} else {
			log.Println("mymessage error   - ", err)
		}
	} else {

	}
}

/*
 *	FUNCTION		: encrypt
 *	DESCRIPTION		:
 *		This function takes a string and a cipher key and uses AES to encrypt the message
 *
 *	PARAMETERS		:
 *		byte[] key	: Byte array containing the cipher key
 *		string message	: String containing the message to encrypt
 *
 *	RETURNS		New Horizons 3000 Secure Communications	:
 *		string encoded	: String containing the encoded user input
 *		error err	: Error message
 */
func encrypt(key []byte, message string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

/*
 *	FUNCTION		: decrypt
 *	DESCRIPTION		:
 *		This function takes a string and a key and uses AES to decrypt the string into plain text
 *
 *	PARAMETERS		:
 *		byte[] key	: Byte array containing the cipher key
 *		string secure	: String containing an encrypted message
 *
 *	RETURNS			:
 *		string decoded	: String containing the decrypted equivalent of secure
 *		error err	: Error message
 */
func decrypt(key []byte, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

/*
 *	FUNCTION		: hashAndSalt
 *	DESCRIPTION		:
 *		This function takes a byte array and creates a slted password hash
 *
 *	PARAMETERS		:
 *		byte[] password	: Byte array containing the password
 *
 *	RETURNS			:
 *		string hash	: String containing the hashed and salted value
 */
func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

/*
 *	FUNCTION		: comparePasswords
 *	DESCRIPTION		:
 *		This function takes a hashed password and compare it to the provided password calculated hash
 *
 *	PARAMETERS		:
 *		string hashedpassword	: password hash
 *      string password         : password to compare
 *
 *	RETURNS			:
 *		bool match	: value of compare
 */
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

/*
 *	FUNCTION		: editEntry
 *	DESCRIPTION		:
 *		This function takes a string and edits it according to action
 *
 *	PARAMETERS		:
 *		string action       	: URL or STRING or CERTIFICATE
 *      string value            : password to compare
 *
 *	RETURNS			:
 *		bool valid	: value of edit
 */
func editEntry(action string, value string) bool {

	if action == "URL" {

		valid := strings.Contains(value, "nats://")
		if valid == false {
			return false
		}
		valid1 := strings.Contains(value, "NATS://")
		if valid1 == false {
			return false
		}
		valid2 := strings.Contains(value, ".")
		if valid2 == false {
			return false
		}
		valid3 := strings.Contains(value, ":")
		if valid3 == false {
			return false
		}

		return true
	}
	if action == "STRING" {
		valid3 := strings.Contains(value, "None")
		if valid3 == false {
			return false
		}
		valid4 := strings.Contains(value, "NONE")
		if valid4 == false {package panes
2
​
3
import (
4
        "fyne.io/fyne/v2"
5
)
6
​
7
// Pane defines the data structure
8
type MyPane struct {
9
        Title, Intro string
10
        View         func(w fyne.Window) fyne.CanvasObject
11
        SupportWeb   bool
12
}
13
​
14
var (
15
        // Panes defines the metadata
16
        MyPanes = map[string]MyPane{
17
                "logon":    {"Logon", "", logonScreen, true},
18
                "messages": {"Messages", "", messagesScreen, true},
19
        }
20
​
21
        // PanesIndex  defines how our panes should be laid out in the index tree
22
        MyPanesIndex = map[string][]string{
23
                "": {"logon", "messages"},
24
        }
25
)
26

			return false
		}

		if len(value) == 0 {
			return false
		}
	}
	if action == "CERTIFICATE" {
		valid := strings.Contains(value, "-----BEGIN CERTIFICATE-----")
		if valid == false {
			return false
		}
		valid2 := strings.Contains(value, "-----END CERTIFICATE-----")
		if valid2 == false {
			return false
		}
	}
	return true

}
