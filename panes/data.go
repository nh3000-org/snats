/*
 *	PROGRAM		: data.go
 *	DESCRIPTION		:
 *
 *		This program handles field definitions and common functions
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

	//"encoding/base64"
	"encoding/base64"
	//"encoding/hex"

	//"encoding/hex"

	"github.com/google/uuid"

	//"encoding/json"
	//"errors"
	"fmt"

	"github.com/goccy/go-json"

	"log"
	"net"
	"net/url"
	"os"

	"strconv"
	"strings"

	"github.com/nats-io/nats.go"

	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"

	//"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

/*
*  The following fields need to be modified for you production
*  Environment to provide maximum security
*
*  These fields are meant to be distributed at compile time and
*  editable in the gui.
*
 */
var MyBytes = []byte{35, 46, 57, 66, 32, 14, 05}

const MySecret string = "abd&^7%c34"
const MyDurable string = "snatsdurable"
const PasswordDefault = "123456" // default password shipped with app
var xCaroot = string("-----BEGIN CERTIFICATE-----HO1DwrlkTzUimG5PoiswB6swCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAgFw0yMzAzMzExNzI5MDBaGA8yMDUzMDMyMzE3MjkwMFowZjELMAkG\nA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UEChMDU0VD\nMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9y\nZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHXwMUfMXiJix3tuzFymcA+3RkeY\nZE7urUzVgaqkv/Oef3jhqhtf1XzK/qVYGxWWmpvADGB252PG1Mp7Z5wmzqyjRTBD\nMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMB0GA1UdDgQWBBQm\nFA5caanuqxGFOf9DtZkVYv5dCzAKBggqhkjOPQQDAgNHADBEAiB3BheNP4XdBZ27\nxVBQ7ztMJqK7wDi1V3LuMy5jmXr7rQIgHCse0oaiAwcl4VwF00aSshlV+T/da0Tx\n1ANkaM+rie4=\n-----END CERTIFICATE-----\n")
var xClientcert = string("-----BEGIN CERTIFICATE-----UUyhlJt8mp1XApRbSkdrUS55LGV8wCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAeFw0yMzAzMzExNzI5MDBaFw0yODAzMjkxNzI5MDBaMHIxCzAJBgNV\nBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlDcmVzdHZpZXcxGjAY\nBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMSEwHwYDVQQLExhuYXRzLm5ld2hvcml6\nb25zMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDFttVH\nQ131JYwazAQMm0XAQvRvTjTjOY3aei1++mmQ+NQ9mrOFk6HlZFoKqsy6+HPXsB9x\nQbWlYvUOuqBgb9xFQZoL8jiKskLLrXoIxUAlIBTlyf76r4SV+ZpxJYoGzXNTedaU\n0EMTyAiUQ6nBbFMXiehN5q8VzxtTESk7QguGdAUYXYsCmYBvQtBXoFYO5CHyhPqu\nOZh7PxRAruYypEWVFBA+29+pwVeaRHzpfd/gKLY4j2paInFn7RidYUTqRH97BjdR\nSZpOJH6fD7bI4L09pnFtII5pAARSX1DntS0nWIWhYYI9use9Hi/B2DRQLcDSy1G4\n0t1z4cdyjXxbFENTAgMBAAGjgawwgakwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQM\nMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFAzgPVB2/sfT7R0U\ne3iXRSvUkfoQMB8GA1UdIwQYMBaAFCYUDlxpqe6rEYU5/0O1mRVi/l0LMDQGA1Ud\nEQQtMCuCGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZ4IJMTI3LDAsMCwxhwTAqABn\nMAoGCCqGSM49BAMCA0kAMEYCIQCDlUH2j69mJ4MeKvI8noOmvLHfvP4qMy5nFW2F\nPT5UxgIhAL6pHFyEbANtSkcVJqxTyKE4GTXcHc4DB43Z1F7VxSJj\n-----END CERTIFICATE-----\n")
var xClientkey = string("-----BEGIN RSA PRIVATE KEY----AQEAxbbVR0Nd9SWMGswEDJtFwEL0b0404zmN2notfvppkPjUPZqz\nhZOh5WRaCqrMuvhz17AfcUG1pWL1DrqgYG/cRUGaC/I4irJCy616CMVAJSAU5cn+\n+q+ElfmacSWKBs1zU3nWlNBDE8gIlEOpwWxTF4noTeavFc8bUxEpO0ILhnQFGF2L\nApmAb0LQV6BWDuQh8oT6rjmYez8UQK7mMqRFlRQQPtvfqcFXmkR86X3f4Ci2OI9q\nWiJxZ+0YnWFE6kR/ewY3UUmaTiR+nw+2yOC9PaZxbSCOaQAEUl9Q57UtJ1iFoWGC\nPbrHvR4vwdg0UC3A0stRuNLdc+HHco18WxRDUwIDAQABAoIBACe0XMZP4Al//c/P\n0qxZbjt69q13jiVnhHYwfPx3+0UywySP8adMi4GOkop73Ftb05+n7diHspvA8KeB\nkP1s2VZLI01s2i/4NnPCpbQnMIeEFs5Cr2LWZpDbrEk2ma5eCd/kotQFssLBM//a\nSrfeMh2TA0TJo7WEft9Cnf4ZeEkKnycplfvwTyv286iFZCYo2dv66BfTej6kkVCo\nAi+ZVCe2zSqRYyr0u4/j/kE3b3eSkCnY2IVcqlP7epuEGVOZyxeFLwM5ljbWL816\npA6WIJgQo2EQ1N7L531neg5WjXQ/UwTQoXP1jvuuVtKtOBFqm1IshEyFk3WpsfpD\nr16OTdECgYEA6FB6NYxYtnWPaIYAOqP7GtMKoJujH8MtZy6J33LkxI7nPkMkn8Mv\nva32tvjU4Bu1FVNp9k5guC+b+8ixXK0URj25IOhDs6K57tck22W9WiTZlmnkCO01\nJOavrelWbvYt5xNWIdnPualoPfGB0iJKXsKY/bpH4eVfhWwpNPI5sMkCgYEA2d9G\nEPuWN6gUjZ+JfdS+0WHK1yGD7thXs7MPUlhGqDzBryh2dkywyo8U8+tMLuDok1RZ\njnT3PYkLQEpzoV0qBkpFFShL6ubaGmDz1UZsozl0YcIg4diZeuPHnIAeXOFrhgYf\n825163LmT3jYHCROFEMLtTYyIQP0EznE+qFT3TsCgYEApgtvbfqkJbWdDL5KR5+R\nCLky7VyQmVEtkIRI8zbxoDPrwCrJcI9X/iDrKBhuPshPA7EdGXkn1D3jJXFqo6zp\nwtK3EXgxe6Ghd766jz4Guvl/s+x3mpHA3GEtzAXtS14VrQW7GHLP8AnPggauHX14\n3oYER8XvPtxtC7YlNbyz01ECgYAe2b7SKM3ck7BVXYHaj4V1oKNYUyaba4b/qxtA\nTb+zkubaJqCfn7xo8lnFMExZVv+X3RnRUj6wN/ef4ur8rnSE739Yv5wAZy/7DD96\ns74uXrRcI2EEmechv59ESeACxuiy0as0jS+lZ1+1YSc41Os5c0T1I/d1NVoaXtPF\nqZJ2gQKBgBp/XavdULBPzC7B8tblySzmL01qJZV7MSSVo2/1vJ7gPM0nQPZdTDog\nTfA5QKSX9vFTGC9CZHSJ+fabYDDd6+3UNYUKINfr+kwu9C2cysbiPaM3H27WR5mW\n5LhStAfwuRRYBDsG2ndjraxcBrrPdtkbS0dpeQUDJxvkMIuLHnhQ\n-----END RSA PRIVATE KEY-----\n")

// version
const Version = "snats-beta"

// messages from nats
var NatsMessages []MessageStore

// logon status

var LoggedOn bool = false
var PasswordValid bool = false

// error message
var ErrorMessage = "None"

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
var PasswordMinimumSize string      // set minimum password size
var PasswordMustContainNumber bool  // password must contain number
var PasswordMustContainLetter bool  // password must contain letter
var PasswordMustContainSpecial bool // password must contain special character

// Server tab
var Server string // server url

var IdUUID string       // unique message id
var Caroot string       // CAROOT certificate for server authentication
var Clientcert string   // Client cert signed by caroot
var Clientkey string    // Client key signed by carooit
var UserID string       // NATS user id
var UserPassword string // NATS password for crypto operations
var Alias string        // name the queue user
var NodeUUID string     // nodeuuid created on logon
//var Nonce string        // nonce for encrypt/decrypt
//var Noncesize int       // nonce for encrypt/decrypt
/*
 *	  These constants are set to establish a password schema for Local File Encryption and Queue password
 */

/*
 *	  Confignats is used to hold config.json fields
 */
type Confignats struct {
	Jserver                     string `json:"server"`                     // server url nats://333.333.333.333:port
	Jqueue                      string `json:"queue"`                      // queue created for deployment
	Jqueuepassword              string `json:"queuepassword"`              // queue password created for deployment
	Jpasswordminimumsize        string `json:"passwordminimumsize"`        // set minimum password size
	Jpasswordmustcontainnumber  bool   `json:"passwordmustcontainnumber"`  // password must contain number
	Jpasswordmustcontainletter  bool   `json:"passwordmustcontainletter"`  // password must contain letter
	Jpasswordmustcontainspecial bool   `json:"passwordmustcontainspecial"` // password must contain special character
	//Jusejetstream               bool   `json:"usejetstream"`               // if set to true uses jetstream protocol otherwise regular pub/sub
	//Jusetls                     bool   `json:"usetls"`                   // caroot self signed ca
	Jcaroot       string `json:"caroot"`       // CAROOT certificate for server authentication
	Jclientcert   string `json:"clientcert"`   // Client certificate for server authentication
	Jclientkey    string `json:"clientkey"`    // client key for server authentication
	Juserid       string `json:"userid"`       // user id
	Juserpassword string `json:"userpassword"` // user password
	Jalias        string `json:"alias"`        // user alias
	Jnodeuuid     string `json:"nodeuuid"`     // node id created on logon

}

// Pane defines the data structure
type MyPane struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

type MessageStore struct {
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
}

var (
	// Panes defines the metadata
	MyPanes = map[string]MyPane{
		"password": {"Password", "", passwordScreen, true},
		"logon":    {"Logon", "", logonScreen, true},
		"settings": {"Settings", "", settingsScreen, true},

		"messages": {"Messages", "", messagesScreen, true},
	}

	// PanesIndex  defines how our panes should be laid out in the index tree
	MyPanesIndex = map[string][]string{
		"": {"password", "logon", "settings", "messages"},
	}
)

/*
 *	FUNCTION		: DataStore
 *	DESCRIPTION		:
 *		Handle access to storage
 *
 *	PARAMETERS		: filename
 *		        	:
 *
 *	RETURNS			: uri
 *
 */
func DataStore(myfile string) fyne.URI {

	DataLocation, dlerr := storage.Child(fyne.CurrentApp().Storage().RootURI(), myfile)
	if dlerr != nil {
		log.Println("DataStore error ", dlerr)
	}

	return DataLocation
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
 *	FUNCTION		: myjson
 *	DESCRIPTION		:
 *		This function handles file actions for config.json to load memory
 *
 *	PARAMETERS		        :
 *		action string   	: CREATE, LOAD or SAVE encrypted fields


 *	RETURNS			:
 *		         	: None
 */
func MyJson(action string) {

	if action == "CREATE" {
		// backwards server first, encrypt inot struct write
		var c = Confignats{}
		Server = string("nats://192.168.0.103:4222")

		Caroot = strings.ReplaceAll(xCaroot, "\n", "<>")

		Clientcert = strings.ReplaceAll(xClientcert, "\n", "<>")

		Clientkey = strings.ReplaceAll(xClientkey, "\n", "<>")

		Queue = string("MESSAGES")

		Alias = string("Myalias")

		Queuepassword = string("123456")

		PasswordMinimumSize = string("6")

		PasswordMustContainNumber = bool(false)

		PasswordMustContainLetter = bool(false)

		PasswordMustContainSpecial = bool(false)

		UserID = string("None")

		UserPassword = string("None")

		MyCrypt("ENCRYPT")
		c.Jserver = Server
		c.Jcaroot = Caroot
		c.Jclientcert = Clientcert
		c.Jclientkey = Clientkey
		c.Jqueue = Queue
		c.Jalias = Alias
		c.Jqueuepassword = Queuepassword
		c.Jpasswordminimumsize = PasswordMinimumSize
		c.Jpasswordmustcontainnumber = PasswordMustContainNumber
		c.Jpasswordmustcontainletter = PasswordMustContainLetter
		c.Jpasswordmustcontainspecial = PasswordMustContainSpecial
		c.Juserid = UserID
		c.Juserpassword = UserPassword

		wrt, errwrite := storage.Writer(DataStore("config.json"))
		if errwrite != nil {
			log.Println("SaveConfig Error Writing", DataStore("config.json"))
		}

		enc := json.NewEncoder(wrt)
		err2 := enc.Encode(c)

		if err2 != nil {
			log.Println("SaveConfig Error Saving", DataStore("config.json"))
		}

		wrt.Close()
		c = Confignats{}
		MyCrypt("DECRYPT")
	}

	if action == "LOAD" {
		log.Println("MyJson LOAD")
		jsonfile, errf := os.ReadFile(DataStore("config.json").Path())
		//log.Println("MyJson LOAD file", jsonfile)
		if errf != nil {
			log.Println("MyJson LOAD Error file", errf)
		}
		myc := Confignats{}
		err := json.Unmarshal([]byte(jsonfile), &myc)
		if err != nil {
			log.Println("MyJson LOAD Unmarshall: ", err)
		}

		Server = myc.Jserver
		Caroot = myc.Jcaroot
		Clientkey = myc.Jclientkey
		Clientcert = myc.Jclientcert
		Queue = myc.Jqueue
		Queuepassword = myc.Jqueuepassword
		PasswordMinimumSize = myc.Jpasswordminimumsize
		PasswordMustContainLetter = myc.Jpasswordmustcontainletter
		PasswordMustContainNumber = myc.Jpasswordmustcontainnumber
		PasswordMustContainSpecial = myc.Jpasswordmustcontainspecial
		//UseJetstream = myc.Jusejetstream
		//UseTLS = myc.Jusetls
		UserID = myc.Juserid
		UserPassword = myc.Juserpassword
		Alias = myc.Jalias
		NodeUUID = myc.Jnodeuuid

		MyCrypt("DECRYPT")
		var xCaroot = strings.ReplaceAll(Caroot, "<>", "\n")
		Caroot = xCaroot
		var xClientcert = strings.ReplaceAll(Clientcert, "<>", "\n")
		Clientcert = xClientcert
		var xClientkey = strings.ReplaceAll(Clientkey, "<>", "\n")
		Clientkey = xClientkey

	}
	if action == "SAVE" {
		log.Println("MyJson SAVE")
		MyCrypt("ENCRYPT")

		myc := Confignats{}

		myc.Jserver = Server
		myc.Jcaroot = strings.ReplaceAll(Caroot, "\n", "<>")
		myc.Jclientcert = strings.ReplaceAll(Clientcert, "\n", "<>")
		myc.Jclientkey = strings.ReplaceAll(Clientkey, "\n", "<>")

		myc.Jqueue = Queue
		myc.Jqueuepassword = Queuepassword
		myc.Jpasswordminimumsize = PasswordMinimumSize
		myc.Jpasswordmustcontainletter = PasswordMustContainLetter
		myc.Jpasswordmustcontainnumber = PasswordMustContainNumber
		myc.Jpasswordmustcontainspecial = PasswordMustContainSpecial
		//myc.Jusejetstream = UseJetstream
		//myc.Jusetls = UseTLS
		myc.Juserid = UserID
		myc.Juserpassword = UserPassword
		myc.Jalias = Alias
		myc.Jnodeuuid = NodeUUID

		err := storage.Delete(DataStore("config.json"))
		if err != nil {
			log.Println("SaveConfig Error Deleting", DataStore("config.json"))
		}

		wrt, errwrite := storage.Writer(DataStore("config.json"))
		if errwrite != nil {
			log.Println("SaveConfig Error Writing", DataStore("config.json"))
		}
		enc := json.NewEncoder(wrt)
		err2 := enc.Encode(myc)
		if err2 != nil {
			log.Println("SaveConfig Error Saving", DataStore("config.json"))
		}
		MyCrypt("DECRYPT")
		//SaveCarootToFS()

	}
}

/*
 *	FUNCTION		: Encode bytes to base64 encoding
 *	DESCRIPTION		:
 *		This function encodes a string
 *
 *	PARAMETERS		:  bytes to encode

 *
 *	RETURNS			: bytes as encoded string
 *		         	: None
 */
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

/*
 *	FUNCTION		: Decode string into bytes
 *	DESCRIPTION		:
 *		This function decodes a string
 *
 *	PARAMETERS		:  string to encode

 *
 *	RETURNS			: string as bytes
 *		         	: None
 */
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

/*
 *	FUNCTION		: MyCrypt
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

		cryptoText, _ := Encrypt(Server, MySecret)

		Server = cryptoText

		cryptoText1, _ := Encrypt(Caroot, MySecret)
		Caroot = cryptoText1

		cryptoText2, _ := Encrypt(Clientcert, MySecret)
		Clientcert = cryptoText2
		cryptoText3, _ := Encrypt(Clientkey, MySecret)
		Clientkey = cryptoText3

		cryptoText4, _ := Encrypt(Queue, MySecret)
		Queue = cryptoText4

		cryptoText5, _ := Encrypt(Queuepassword, MySecret)
		Queuepassword = cryptoText5

		cryptoText6, _ := Encrypt(UserID, MySecret)
		UserID = cryptoText6

		cryptoText7, _ := Encrypt(UserPassword, MySecret)
		UserPassword = cryptoText7

		cryptoText8, _ := Encrypt(Alias, MySecret)
		Alias = cryptoText8

		cryptoText9, _ := Encrypt(NodeUUID, MySecret)
		NodeUUID = cryptoText9

	}
	if action == "DECRYPT" {
		text, _ := Decrypt(Server, MySecret)
		Server = text
		text1, _ := Decrypt(Caroot, MySecret)
		Caroot = text1
		text2, _ := Decrypt(Clientcert, MySecret)
		Clientcert = text2
		text3, _ := Decrypt(Clientkey, MySecret)
		Clientkey = text3
		text4, _ := Decrypt(Queue, MySecret)
		Queue = text4
		text5, _ := Decrypt(Queuepassword, MySecret)
		Queuepassword = text5
		text6, _ := Decrypt(UserID, MySecret)
		UserID = text6
		text7, _ := Decrypt(UserPassword, MySecret)
		UserPassword = text7
		text8, _ := Decrypt(Alias, MySecret)
		Alias = text8
		text9, _ := Decrypt(NodeUUID, MySecret)
		NodeUUID = text9

	}
}

/*
 *	FUNCTION		: Encrypt
 *	DESCRIPTION		: Encrypt a field
 *	PARAMETERS	    : text to encrypt, secret
 *
 *	RETURNS			: enrypted text

 */
func Encrypt(text string, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, MyBytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

/*
 *	FUNCTION		: Decrypt
 *	DESCRIPTION		: Decrypt a field
 *	PARAMETERS	    : text to decrypt, secret
 *
 *	RETURNS			: derypted text

 */
func Decrypt(text string, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, MyBytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
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
		//log.Println("MyHash  CREATE ", DataStore("config.hash"))
		err := storage.Delete(DataStore("config.hash"))
		if err != nil {
			log.Println("MyHash Error Deleting", DataStore("config.hash"))
		}
		wrt, errwrite := storage.Writer(DataStore("config.hash"))
		_, err2 := wrt.Write([]byte(Passwordhash))
		if errwrite != nil || err2 != nil {
			log.Println("MyHash Error Writing", DataStore("config.hash"))
		}

	}
	if action == "LOAD" {
		//log.Println("MyHash  LOAD", DataStore("config.hash"))
		ph, errf := os.ReadFile(DataStore("config.hash").Path())
		Passwordhash = string(ph)

		if errf != nil {
			log.Println("MyHash LOAD Hash Error file", errf, " ", Passwordhash)
		}

	}
	if action == "SAVE" {
		//log.Println("MyHash Error save ", DataStore("config.hash"))
		errf := storage.Delete(DataStore("config.hash"))

		if errf != nil {
			log.Println("MyHash SAVE Hash Error file", errf)
		}
		wrt, errwrite := storage.Writer(DataStore("config.hash"))
		_, err2 := wrt.Write([]byte(Passwordhash))
		if errwrite != nil || err2 != nil {
			log.Println("MyHash SAVE Error Writing", DataStore("config.hash"))
		}

	}
}

/*
 *	FUNCTION		: NATSErase
 *	DESCRIPTION		:
 *		This function erases a messages in queue
 *
 *	PARAMETERS		:
 *
 *	RETURNS		!	:
 */
func NATSErase() {
	log.Println("Erasing  ")

	//nc, se1 := nats.Connect(Server, nats.RootCAs(DataStore("ca-root.pem").Path()), nats.ClientCert(DataStore("client-cert.pem").Path(), DataStore("client-key.pem").Path()))
	nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
	if err != nil {
		log.Println("NatsErase Connection ", err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("NatsErase JetStream ", err)
	}

	NatsMessages = nil
	// Delete Stream
	err1 := js.DeleteStream(Queue)
	if err != nil {
		log.Println("NatsErase DeleteStream ", err1)
	}

	var duration time.Duration = 604800000000

	js1, err1 := js.AddStream(&nats.StreamConfig{
		Name:     Queue,
		Subjects: []string{(Queue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   duration,

		Discard:           nats.DiscardOld,
		MaxMsgsPerSubject: 1000,
		AllowRollup:       true,
	})

	if err1 != nil {
		log.Println("NatsErase AddStream ", err1)
	}
	fmt.Printf("js1: %v\n", js1)
	ac, err1 := js.AddConsumer(Queue, &nats.ConsumerConfig{
		//Durable:   Alias,
		Durable: MyDurable,

		InactiveThreshold: duration,
		DeliverPolicy:     nats.DeliverAllPolicy,
		//		ReplayPolicy: nats.ReplayInstantPolicy,
	})
	if err1 != nil {
		log.Println("NatsErase AddConsumer ", err1, " ", ac)
	}

	nc.Close()

}

/*
 *	FUNCTION		: FormatMessage
 *	DESCRIPTION		:
 *		This function formats a message for sending
 *
 *	PARAMETERS		:
 *
 *	RETURNS		!	:  MessageStore

 */
func FormatMessage(m string) string {
	EncMessage := MessageStore{}

	//ID , err := exec.Command("uuidgen").Output()

	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "No Host Name"
		//strings.Replace(EncMessage, "#HOSTNAME", "No Host Name", -1)

	} else {
		EncMessage.MShostname = name
		//strings.Replace(EncMessage, "#HOSTNAME", name, -1)
	}
	ifas, err := net.Interfaces()
	if err == nil {

		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		for i, s := range as {
			EncMessage.MShostname += " mac: " + strconv.Itoa(i) + " - " + s + ","
		}
		addrs, _ := net.InterfaceAddrs()

		for _, addr := range addrs {
			EncMessage.MShostname += "\n addrs: " + addr.String() + ","
		}

	}
	EncMessage.MSnodeuuid = NodeUUID

	//	addrs, err := net.LookupHost(name)
	//	var addresstring = ""
	//	if err == nil {
	//		for _, a := range addrs {
	//			addresstring += a
	//			addresstring += ","
	//		}
	//		EncMessage.MSipadrs = "ip " + addresstring
	//
	//	} else {
	//		EncMessage.MSipadrs = "No IP"
	//	}
	EncMessage.MSalias = Alias

	EncMessage.MSiduuid = "Node: " + uuid.New().String()

	EncMessage.MSmessage = m
	//EncMessage += m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	ejson, _ := Encrypt(string(jsonmsg), MySecret)
	return string(ejson)

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

	if action == "cvtbool" {
		if value == "True" {
			return true
		}
		if value == "False" {
			return false
		}

	}

	if action == "URL" {

		valid := strings.Contains(value, "nats://")
		if valid == false {
			valid1 := strings.Contains(value, "NATS://")
			if valid1 == false {
				return false
			}
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
		if valid4 == false {
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
	if action == "KEY" {
		valid := strings.Contains(value, "-----BEGIN RSA PRIVATE KEY-----")
		if valid == false {
			return false
		}
		valid2 := strings.Contains(value, "-----END RSA PRIVATE KEY-----")
		if valid2 == false {
			return false
		}
	}
	if action == "TRUEFALSE" {
		valid := strings.Contains(value, "True")
		if valid == false {
			valid2 := strings.Contains(value, "False")
			if valid2 == false {
				return false
			}
		}
	}
	return true

}
