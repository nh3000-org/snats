package panes

import (
	//	"natsgui/pkg/cmd/main"

	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	//	"errors"
	"encoding/json"

	//	"io/ioutil"
	"fyne.io/fyne/v2"
	//	"fyne.io/fyne/v2/canvas"
	//	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

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

//var config ConfigNats

type Confignats struct {
	Jserver        string
	Jcaroot        string
	Jqueue         string
	Jqueuepassword string
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func dumpglobals(from string) {
	log.Println(from, Password)
	log.Println(from, Server)
	log.Println(from, Caroot)
	log.Println(from, Queue)
	log.Println(from, Queuepassword)
}

func loadjsonfromglobals() map[string]interface{} {
	log.Println("loadjson Password", Password)
	data := map[string]interface{}{
		"server":        string(Server),
		"caroot":        string(Caroot),
		"queue":         string(Queue),
		"queuepassword": string(Queuepassword),
	}
	log.Println("loadjson data", data)
	return data
}
func MyJson(action string) {

	if action == "CREATE" {
		log.Println("create Password", Password)
		Server = "None"
		Caroot = "None"
		Queue = "None"
		Queuepassword = "None"
		configfile, configfileerr := os.Create("config.json")
		if configfileerr == nil {
			enc := json.NewEncoder(configfile)

			log.Println("myjson save config", loadjsonfromglobals())
			enc.Encode(loadjsonfromglobals())
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
		log.Println("myjson load jc", jc)
		if je != nil {
			log.Println("LOAD Error read all", je)
		}
		jf.Close()

		json.Unmarshal([]byte(jc), &c)
		for k, v := range c {
			fmt.Println(k, "=>", v)
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

	}
	if action == "SAVE" {
		e := os.Remove("config.json")
		if e != nil {
			log.Fatal(e)
		}
		sc, se := os.Create("config.json")

		if se == nil {
			enc := json.NewEncoder(sc)
			log.Println("myjson save", loadjsonfromglobals())
			enc.Encode(loadjsonfromglobals())

		}

		sc.Close()
	}
}

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {

	_, configfileerr := os.Stat("config.json")
	if configfileerr != nil {

		MyJson("CREATE")
	}

	MyJson("LOAD")

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
		MyJson("LOAD")
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
		dumpglobals("myjson try passworde")
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

		MyJson("SAVE")

		dumpglobals("myjson save")
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

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
	log.Println(stringToEncrypt)
	log.Println(keyString)
	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
