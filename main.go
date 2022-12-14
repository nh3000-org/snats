/*
* Modify cipherkey for your installation
 */

package panes

import (
	//	"natsgui/pkg/cmd/main"

	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	//	"errors"
	"encoding/base64"
	"encoding/json"

	//	"io/ioutil"
	"fyne.io/fyne/v2"
	//	"fyne.io/fyne/v2/canvas"
	//	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	//	"encoding/hex"
	"io"

	"golang.org/x/crypto/bcrypt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	//	"reflect"
	//   "github.com/ilyakaznacheev/cleanenv"
	//	"github.com/gookit/config/v2"
	//	"github.com/gookit/config/v2/json"
)

var Password string      // encrypt file password
var Caroot string        // CAROOT certificate for server authentication
var Queue string         // server message queue
var Queuepassword string // server message queue password
var Server string        // server url

var Passwordhash string                              // hash value of password
const Cipherkey = "asuperstrong32bitpasswordgohere!" //string    // hash value of cipher key to decrypt json fields modify this field for your ntwork

//var config ConfigNats

type Confignats struct {
	Jserver        string
	Jcaroot        string
	Jqueue         string
	Jqueuepassword string
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
		Server, _ = encrypt([]byte("None"), Cipherkey)
		Caroot, _ = encrypt([]byte("None"), Cipherkey)
		Queue, _ = encrypt([]byte("None"), Cipherkey)
		Queuepassword, _ = encrypt([]byte("None"), Cipherkey)
		Encmessage, _ = encrypt([]byte("None"), Cipherkey)
		configfile, configfileerr := os.Create("config.json")
		if configfileerr == nil {
			enc := json.NewEncoder(configfile)

			MyCrypt("ENCRYPT")
			enc.Encode(loadconfig())
		}
		configfile.Close()
	}
	if action == "LOAD" {
		//var c Confignats
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
	if action == "ENCRYPT" {
		var newvalue, _ = encrypt([]byte(Server), Cipherkey)
		Server = newvalue
		newvalue, _ = encrypt([]byte(Caroot), Cipherkey)
		Caroot = newvalue
		newvalue, _ = encrypt([]byte(Queue), Cipherkey)
		Queue = newvalue
		newvalue, _ = encrypt([]byte(Queuepassword), Cipherkey)
		Queuepassword = newvalue

	}
	if action == "DECRYPT" {
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

func logonScreen(_ fyne.Window) fyne.CanvasObject {

	_, configfileerr := os.Stat("config.json")
	if configfileerr != nil {

		myjson("CREATE")
	}
	myjson("LOAD")

	password := widget.NewEntry()
	password.SetPlaceHolder("Enter Password For Encryption")
	password.SetText(Password)
	server := widget.NewEntry()
	server.SetPlaceHolder("URL: nats://xxxxxx:4332")
	server.Disable()
	caroot := widget.NewMultiLineEntry()
	caroot.SetPlaceHolder("CAROOT For nats://xxxxxx:4332")
	caroot.Disable()
	queue := widget.NewEntry()
	queue.SetPlaceHolder("Message Queue for Pub/Sub")
	queue.Disable()
	queuepassword := widget.NewEntry()
	queuepassword.SetPlaceHolder("Message Queue Password")
	queuepassword.Disable()
	// try the password
	tpbutton := widget.NewButton("Try Password", func() {
		Password = password.Text
		log.Println("TP password", Password)
		pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		Passwordhash = string(pwh)
		log.Println("TP Passwordhash", Passwordhash)
		if err != nil {
			log.Fatal(err)
		}
		_, confighasherr := os.Stat("config.hash")
		if confighasherr != nil {

			MyHash("CREATE", Passwordhash)
		}

		Password = password.Text
		MyHash("LOAD", "NONE")
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {
			// TODO: Properly handle error
			log.Fatal(err)
		}

		myjson("LOAD")
		log.Println("save password config", Password)
		log.Println("save password gui", password.Text)

		server.SetText(Server)
		caroot.SetText(Caroot)
		queue.SetText(Queue)
		queuepassword.SetText(Queuepassword)
		password.Disable()
		server.Enable()
		caroot.Enable()
		queue.Enable()
		queuepassword.Enable()
		//dumpglobals("myjson try password")
	})
	// save the server
	ssbutton := widget.NewButton("Save Server", func() {
		Password = password.Text
		Server = server.Text
		Caroot = caroot.Text
		Queue = queue.Text
		Queuepassword = queuepassword.Text
		password.Disable()
		server.Disable()
		caroot.Disable()
		queue.Disable()
		queuepassword.Disable()

		myjson("SAVE")

		//dumpglobals("myjson save")
	})
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		server,
		caroot,
		queue,
		queuepassword,
		ssbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.og", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

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
 *	RETURNS			:
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
