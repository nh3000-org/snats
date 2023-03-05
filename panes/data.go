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
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

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

var UseJetstream bool   // if set to true uses jetstream protocol otherwise regular pub/sub
var UseTLS bool         // use TLS to Authenticate  else use userid /password
var Caroot string       // for UseTLS = true CAROOT certificate for server authentication
var UserID string       // for UseTLS = false
var UserPassword string // for UseTLS = false
var Alias string        // name the user
var NodeUUID string     // nuuid created on logon

/*
 *	  These constants are set to establish a password schema for Local File Encryption and Queue password
 */

const Cipherkey = "asuperstrong32bitpasswordgohere!" // 32 byte string  for hash value of cipher key to decrypt json fields modify this field for your ntwork
const PasswordDefault = "123456"                     // default password shipped with app
//const MessageFormat = "HostName: = #HOSTNAME IPs : #IPS\n Message: #MESSAGE\n Date/Time #DATETIME\n" // default message for posting
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
	Jusejetstream               bool   `json:"usejetstream"`               // if set to true uses jetstream protocol otherwise regular pub/sub
	Jusetls                     bool   `json:"usetls"`                     // use TLS to Authenticate  else use userid /password
	Jcaroot                     string `json:"caroot"`                     // for UseTLS = true CAROOT certificate for server authentication
	Juserid                     string `json:"userid"`                     // for UseTLS = false
	Juserpassword               string `json:"userpassword"`               // for UseTLS = false
	Jalias                      string `json:"alias"`                      // user alias
	Jnodeuuid                   string `json:"nodeuuid"`                   // node id created on logon
}

// Pane defines the data structure
type MyPane struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

type MessageStore struct {
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
		log.Println("MyJson Create ", Password)
		var c = Confignats{}
		c.Jserver = string("nats://192.168.0.164:4222")
		var cert = string("-----BEGIN CERTIFICATE-----\nMIIDKjCCAhKgAwIBAgIUXswOBqsYOrtEKr+5QZZ0/Pwtg/cwDQYJKoZIhvcNAQEL\nBQAwLTEPMA0GA1UEChMGT1JBQ0xFMRowGAYDVQQDExFOZXcgSG9yaXpvbnMgMzAw\nMDAeFw0yMjA4MjUyMTUxMDBaFw0yNzA4MjQyMTUxMDBaMC0xDzANBgNVBAoTBk9S\nQUNMRTEaMBgGA1UEAxMRTmV3IEhvcml6b25zIDMwMDAwggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQC1DYy62ptxKOik3r76SJEgVKDKlVKBL1/RKB3KxiTs\n7ym7Mf73FttmSUtl8lljJVJKaksSk0xxaLwHq5EwZyPqkOcIuqzgrHieL3P17qWE\nqSfOMDVVVEOVXmCOjEqsYDjb2YeV+zvCzq7o9kS97+/muslczWQkT+NldLNXSfqi\njwmd6T2vUVEUtd7kdzr1Z/vFwfGUsTOLnD7chljtvfY1NkQVsxwDKoHaWUWrRUbi\ns2Tzkdi8R2pWhXg5eQHhNe3dRqWYHdoV5Att/IGI6IPsSCjNk78Kj3H40Qkgigdr\n9255hv+kG1ASO2tKxN1Lx+1AGXLNypgtcRb0qOYzRshlAgMBAAGjQjBAMA4GA1Ud\nDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSRqr5ZBI8/t/Dd\nvN9ebSGDCA9+1jANBgkqhkiG9w0BAQsFAAOCAQEAJTNfLISu5JdRvnGqLRiAb9rI\n1ygf+XCDwVNSrZs63ryEcrXNpfNXN8d9Q3qylGlhw+4XYQxjHQhOHO6miMcbGXC+\nCBmxVewquOHbyku6DwAb+v73rEzXyBvyFVnk3QWFDKm4a0smJH+C9Vvx8TCFmkul\nqiV6yMkIpKZ2UbSnGVjkOv3b103VhnD/Jg++jARie8P4vunA7bP+ybN2gEOfudE+\nne+90HcqJCT8EHtPvB/mjtqHEHrWNJSRu/0TgIdMTScFE8Vi18CWWUre3w6wkDrE\nBRsAAouEThO1Jew01rV1CE04s7rp67XUjggejhaTlm3sLmXiajQAj3hZnd2bhw==\n-----END CERTIFICATE-----")
		cert = strings.ReplaceAll(cert, "\n", "<>")
		c.Jcaroot = cert
		c.Jqueue = string("Announcements")
		c.Jalias = string("Myalias")
		c.Jqueuepassword = string("123456")
		c.Jpasswordminimumsize = string("6")
		c.Jpasswordmustcontainnumber = bool(false)
		c.Jpasswordmustcontainletter = bool(false)
		c.Jpasswordmustcontainspecial = bool(false)
		c.Jusejetstream = bool(true)
		c.Jusetls = bool(false)
		c.Juserid = string("None")
		c.Juserpassword = string("None")
		EncMessage = FormatMessage("Connected")

		log.Println("MyJson CREATE")

		wrt, errwrite := storage.Writer(DataStore("config.json"))
		if errwrite != nil {
			log.Println("SaveConfig Error Writing", DataStore("config.json"))
		}
		MyCrypt("ENCRYPT")

		enc := json.NewEncoder(wrt)
		err2 := enc.Encode(c)

		if err2 != nil {
			log.Println("SaveConfig Error Saving", DataStore("config.json"))
		}
		MyCrypt("DECRYPT")
		wrt.Close()
	}

	if action == "LOAD" {
		log.Println("MyJson LOAD")
		jsonfile, errf := os.ReadFile(DataStore("config.json").Path())
		log.Println("MyJson LOAD file", jsonfile)
		if errf != nil {
			log.Println("MyJson LOAD Error file", errf)
		}
		myc := Confignats{}
		err := json.Unmarshal([]byte(jsonfile), &myc)
		if err != nil {
			log.Println("MyJson LOAD Unmarshall: ", err)
		}
		Server = myc.Jserver
		myc.Jcaroot = strings.ReplaceAll(myc.Jcaroot, "<>", "\n")
		Caroot = myc.Jcaroot
		Queue = myc.Jqueue
		Queuepassword = myc.Jqueuepassword
		PasswordMinimumSize = myc.Jpasswordminimumsize
		PasswordMustContainLetter = myc.Jpasswordmustcontainletter
		PasswordMustContainNumber = myc.Jpasswordmustcontainnumber
		PasswordMustContainSpecial = myc.Jpasswordmustcontainspecial
		UseJetstream = myc.Jusejetstream
		UseTLS = myc.Jusetls
		UserID = myc.Juserid
		UserPassword = myc.Juserpassword
		Alias = myc.Jalias
		NodeUUID = myc.Jnodeuuid

		MyCrypt("DECRYPT")

		SaveCarootToFS()

	}
	if action == "SAVE" {
		log.Println("MyJson SAVE")
		//dumpglobals("from myjson save ")

		myc := Confignats{}

		myc.Jserver = Server
		myc.Jcaroot = strings.ReplaceAll(Caroot, "\n", "<>")
		myc.Jqueue = Queue
		myc.Jqueuepassword = Queuepassword
		myc.Jpasswordminimumsize = PasswordMinimumSize
		myc.Jpasswordmustcontainletter = PasswordMustContainLetter
		myc.Jpasswordmustcontainnumber = PasswordMustContainNumber
		myc.Jpasswordmustcontainspecial = PasswordMustContainSpecial
		myc.Jusejetstream = UseJetstream
		myc.Jusetls = UseTLS
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
		MyCrypt("ENCRYPT")

		enc := json.NewEncoder(wrt)
		err2 := enc.Encode(myc)
		if err2 != nil {
			log.Println("SaveConfig Error Saving", DataStore("config.json"))
		}
		MyCrypt("DECRYPT")

	}
}

/*
 *	FUNCTION		: SaveCarootToFS Public function thats save Caroot certificate to fs
 *	DESCRIPTION		:
 *		This function handles caroot certificate usage
 *
 *	PARAMETERS		        :

 *
 *	RETURNS			:
 *		         	: None
 */
func SaveCarootToFS() {

	err := storage.Delete(DataStore("ca-nats.pem"))
	if err != nil {
		//log.Println("SaveCarootToFS Error Deleting", DataStore("ca-nats.pem"))
	}

	//log.Println("SaveCarootToFS caroot " + Caroot)
	wrt, errwrite := storage.Writer(DataStore("ca-nats.pem"))
	_, err2 := wrt.Write([]byte(Caroot))
	if errwrite != nil || err2 != nil {
		log.Println("SaveCarootToFS Error Writing", DataStore("ca-nats.pem"))
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
	if action == "ENCRYPTnow" {
		var newvalue, _ = encrypt([]byte(Cipherkey), Server)
		Server = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Caroot)
		Caroot = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Queue)
		Queue = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Queuepassword)
		Queuepassword = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), UserID)
		UserID = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), UserPassword)
		UserPassword = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), Alias)
		Alias = newvalue
		newvalue, _ = encrypt([]byte(Cipherkey), NodeUUID)
		NodeUUID = newvalue

	}
	if action == "DECRYPTnow" {
		var newvalue, _ = decrypt([]byte(Server), Cipherkey)
		Server = newvalue
		newvalue, _ = decrypt([]byte(Caroot), Cipherkey)
		Caroot = newvalue
		newvalue, _ = decrypt([]byte(Queue), Cipherkey)
		Queue = newvalue
		newvalue, _ = decrypt([]byte(Queuepassword), Cipherkey)
		Queuepassword = newvalue
		newvalue, _ = decrypt([]byte(Cipherkey), UserID)
		UserID = newvalue
		newvalue, _ = decrypt([]byte(Cipherkey), UserPassword)
		UserPassword = newvalue
		newvalue, _ = decrypt([]byte(Cipherkey), Alias)
		Alias = newvalue
		newvalue, _ = decrypt([]byte(Cipherkey), NodeUUID)
		NodeUUID = newvalue
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
		log.Println("MyHash  CREATE ", DataStore("config.hash"))
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
		log.Println("MyHash  LOAD", DataStore("config.hash"))
		ph, errf := os.ReadFile(DataStore("config.hash").Path())
		Passwordhash = string(ph)

		if errf != nil {
			log.Println("MyHash LOAD Hash Error file", errf, " ", Passwordhash)
		}

	}
	if action == "SAVE" {
		log.Println("MyHash Error save ", DataStore("config.hash"))
		errf := storage.Delete(DataStore("config.hash"))

		if errf != nil {
			log.Println("MyHash SAVE Hash Error file", errf)
		}
		wrt, errwrite := storage.Writer(DataStore("config.hash"))
		_, err2 := wrt.Write([]byte(Passwordhash))
		if errwrite != nil || err2 != nil {
			log.Println("MyHash Error Writing", DataStore("config.hash"))
		}

	}
}

/*
 *	FUNCTION		: NATSPublish
 *	DESCRIPTION		:
 *		This function publishes to the select queue
 *
 *	PARAMETERS		:
 *
 *	RETURNS		!	:
 */
func NATSPublish(mm MessageStore) {
	//log.Println("publishing  ")
	if UseJetstream == false {
		//nc, err := nats.Connect(Server, nats.RootCAs("./ca-nats.pem"))

		//if err == nil {
		//nc.Publish(Queue+".*", []byte(FormatMessage("Client Connected")))
		//nc.Subscribe(Queue+".*", func(msg *nats.Msg) {
		//	NatsMessages = append(NatsMessages, string(msg.Data))

		//})

		//}
	}
	if UseJetstream == true {

		NC, err := nats.Connect(Server, nats.RootCAs(DataStore("ca-nats.pem").Path()))
		c, _ := nats.NewEncodedConn(NC, nats.JSON_ENCODER)
		if err == nil {
			//log.Println("publishing  ", mm)
			c.Publish(Queue, mm)
		} else {
			log.Println("publish ", err)
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

	if UseJetstream == false {
		//nc, err := nats.Connect(Server, nats.RootCAs("./ca-nats.pem"))

		//if err == nil {
		//nc.Publish(Queue+".*", []byte(FormatMessage("Client Connected")))
		//nc.Subscribe(Queue+".*", func(msg *nats.Msg) {
		//	NatsMessages = append(NatsMessages, string(msg.Data))

		//})

		//}
	}
	if UseJetstream == true {
		log.Println("NatsErase ")
		NC, se1 := nats.Connect(Server, nats.RootCAs(DataStore("ca-nats.pem").Path()))
		if se1 != nil {
			log.Println("NatsErase se1 ", se1.Error())
		}
		js, se2 := NC.JetStream()
		if se2 != nil {
			log.Println("NatsErase se2 ", se2)
		}

		// Delete Consumer
		se3 := js.DeleteConsumer(Queue, "MONITOR")
		if se3 != nil {
			log.Println("natsconnect se3 ", se3)
		}
		// delete memory store
		NatsMessages = nil
		// Delete Stream
		se4 := js.DeleteStream(Queue)
		if se4 != nil {
			log.Println("natsconnect se4 ", se4)
		}
		// configure default stream
		// keep messages 1 week
		var duration time.Duration = 604800000000
		cfg := nats.StreamConfig{
			Name:     Queue,
			Subjects: []string{Queue},
			Storage:  nats.FileStorage,
			MaxAge:   duration,
		}
		log.Println("natsconnect se5 ", cfg)
		// add the stream
		js, se5 := NC.JetStream()
		if se5 != nil {
			log.Println("natsconnect se5 ", se5)
		}

		//stream, se6 := js.StreamInfo(Queue)

		//if se6 != nil {
		//	log.Println("natsconnect se6 ", stream)
		//}

		// Create a Consumer
		_, se7 := js.AddConsumer(Queue, &nats.ConsumerConfig{
			Durable:      Queue,
			ReplayPolicy: nats.ReplayInstantPolicy,
		})
		if se7 != nil {
			log.Println("natsconnect se7 ", se7)
		}

	}
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
func FormatMessage(m string) MessageStore {
	EncMessage := MessageStore{}
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
			EncMessage.MShostname += " mac " + strconv.Itoa(i) + " - " + s + ","
		}

	}
	EncMessage.MSnodeuuid = NodeUUID

	addrs, err := net.LookupHost(name)
	var addresstring = ""
	if err == nil {
		for _, a := range addrs {
			addresstring += a
			addresstring += ","
		}
		EncMessage.MSipadrs = addresstring

	} else {
		EncMessage.MSipadrs = "No IP"
	}
	EncMessage.MSalias = Alias
	EncMessage.MSmessage = m
	//EncMessage += m
	return EncMessage

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
