package panes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
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
var MyBytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05} // must be 16 bytes
const MySecret string = "abd&1*~#^2^#s0^=)^^7%c34"                                   // must be 24 characters
const MyDurable string = "snatsdurable"
const PasswordDefault = "123456" // default password shipped with app
var Caroot = ""
var Clientcert = ""
var Clientkey = ""

var MyMap = make(map[string]int)
var MyApp fyne.App
var MyAppDup fyne.App

// version
const Version = "snats-beta.1"

// messages from nats
var NatsMessages []MessageStore

var LoggedOn bool = false
var PasswordValid bool = false

var ErrorMessage = "None"

var Queue string         // server message queue
var Queuepassword string // server message queue password

var Msgmaxage string        // msg age in hours to keep
var PreferedLanguage string // language string
var Password string         // encrypt file password
var Passwordhash string     // hash value of password

var PasswordMinimumSize string        // set minimum password size
var PasswordMustContainNumber string  // password must contain number
var PasswordMustContainLetter string  // password must contain letter
var PasswordMustContainSpecial string // password must contain special character

// Server tab
var Server string // server url

var IdUUID string   // unique message id
var Alias string    // name the queue user
var NodeUUID string // nodeuuid created on logon

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
	MSdate     string
}
// eng esp cmn
var MyLangs = map[string]string{
	"eng-mn-intro-1":     "Encrypted Communications Using NATS ",
	"eng-mn-intro-2":     "for Additional Info.",
	"esp-mn-intro-1":     "Comunicaciones Encriptadas Usando NATS ",
	"esp-mn-intro-2":     "para Información Adicional.",
	"eng-mn-mt":          "Encryption",
	"esp-mn-mt":          "Cifrada",
	"eng-mn-err1":        "Missing panel: ",
	"esp-mn-err1":        "Falta el panel: ",
	"eng-mn-dark":        "Dark",
	"esp-mn-dark":        "Oscura",
	"eng-mn-light":       "Light",
	"esp-mn-light":       "Ligera",
	"eng-ps-title":       "Pass Reset",
	"esp-ps-title":       "Pasar Restablecer",
	"eng-ps-password":    "Enter Original Password",
	"esp-ps-passwordueue) ":    "Ingrese la Contraseña Original",
	"eng-ps-passwordc1":  "Enter New Password",
	"esp-ps-passwordc1":  "Ingrese Nueva Clave",
	"eng-ps-passwordc2":  "Enter New Password Again",
	"esp-ps-passwordc2":  "Ingrese la Nueva Contraseña Nuevamente",
	"eng-ps-trypassword": "Try Password",
	"esp-ps-trypassword": "Probar Contraseña",
	"eng-ps-err1":        "Error Creating Password Hash ",
	"esp-ps-err1":        "Error al Crear Hash de Contraseña ",
	"eng-ps-err2":        "Error Creating Password Hash",
	"esp-ps-err2":        "Error al Crear Hash de Contraseña",
	"eng-ps-err3":        "Error Reading Password Hash",
	"esp-ps-err3":        "Error al Leer el Hash de la Contraseña",
	"eng-ps-err4":        "Error Passwords Do Not Match",
	"esp-ps-err4":        "Las Contraseñas de Error no Coinciden",
	"eng-ps-err5":        "Password Accepted",
	"esp-ps-err5":        "Contraseña Aceptada",
	"eng-ps-chgpassword": "Change Password",
	"esp-ps-chgpassword": "Cambiar la Contraseña",
	"eng-ps-err6":        "Error Pasword 1 Invalid",
	"esp-ps-err6":        "Error Contraseña 1 Inválida",
	"eng-ps-err7":        "Error Pasword Does Not Meet Requirements",
	"esp-ps-err7":        "La Contraseña de Error no Cumple con los Requisitos",
	"eng-ps-err8":        "Error Password 1 Does Not Match Password 2",
	"esp-ps-err8":        "Error La Contraseña 1 no Coincide con la Contraseña 2",
	"eng-ps-err9":        "Error Saving Password Hash",
	"esp-ps-err9":        "Error al Guardar el Hash de la Contraseña",
	"eng-ps-err10":       "Error Reading Password Hash",
	"esp-ps-err10":       "Error al Leer el Hash de la Contraseña",
	"eng-ps-err11":       "Error Invalid Password",
	"esp-ps-err11":       "Error Contraseña Inválida",
	"eng-ps-title1":      "Local File Encryption",
	"esp-ps-title1":      "Cifrado de Archivos Locales",
	"eng-ps-title2":      "Enter Password To Reset",
	"esp-ps-title2":      "Ingrese la Contraseña para Restablecer",
	"eng-ps-title3":      "Enter New Password",
	"esp-ps-title3":      "Ingrese Nueva Clave",
	"eng-ss-title":       "Settings",
	"esp-ss-title":       "Ajustes",
	"eng-ss-la":          "Prefered Language",
	"esp-ss-la":          "Idioma Preferido",
	"eng-ss-pl":           "Minimum Password Length",
	"esp-ss-pl":          "Longitud Mínima de la Contraseña",
	"eng-ss-ma":          "Message Max Age In Hours",
	"esp-ss-ma":          "Edad Máxima del Mensaje en Horas",
	"eng-ss-mcletter":    "Password Must Contain Letter",
	"esp-ss-mcletter":    "La Contraseña Debe Contener una Letra",
	"eng-ss-mcnumber":    "Password Must Contain Number",
	"esp-ss-mcnumber":    "La Contraseña Debe Contener un Número",
	"eng-ss-mcspecial":   "Password Must Contain Special",
	"esp-ss-mcspeciaueue" :     "Change Settings",
	"esp-ss-heading":     "Cambiar Ajustes",
	"eng-cs-title":       "Certificates",
	"esp-cs-title":       "Certificados",
	"eng-cs-ca":          "CAROOT Certificate",
	"esp-cs-ca":          "Certificado CAROOT",
	"eng-cs-cc":          "CLIENT Certificate",
	"esp-cs-cc":          "Certificado CLIENTE",
	"eng-cs-ck":          "CLIENT Key",
	"esp-cs-ck":          "Clave CLIENTE",
	"eng-cs-ss":          "Save Certificates",
	"esp-cs-ss":          "Guardar Certificados",
	"eng-cs-err1":        "Error CAROOT is Invalid",
	"esp-cs-err1":        "Error CAROOT no es Válido",
	"eng-cs-err2":        "Error CLIENT CERTIFICATE is invalid ",
	"esp-cs-err2":        "Error CERTIFICADO DE CLIENTE no es Válido",
	"eng-cs-err3":        "Error CLIENT KEY is Invalid",
	"esp-cs-err3":        "Error CLAVE DE CLIENTE no es Cálida",
	"eng-cs-heading":     "Certificate Management",
	"esp-cs-heading":     "Gestión de Certificados",
	"eng-cs-lf":          "Logon First",
	"esp-cs-lf":          "Iniciar Sesión Primero",
	"eng-ls-title":       "Logon",
	"esp-ls-title":       "Iniciar sesión",
	"eng-ls-password":    "Password For Local Encryption",
	"esp-ls-password":    "Contraseña para el Cifrado Local",
	"eng-ls-alias":       "Alias",
	"esp-ls-alias":       "Alias",
	"eng-ls-queue":       "Queue",
	"esp-ls-queue":       "Cola",
	"eng-ls-queuepass":   "Queue Password",
	"esp-ls-queuepass":   "Contraseña de Cola",
	"eng-ls-trypass":     "Try Password",
	"esp-ls-trypass":     "Probar Contraseña",
	"eng-ls-err1":        "logon.go Error Creating Password Hash 24",
	"esp-ls-err1":        "logon.go Error al Crear el Hash de la Contraseña 24",
	"eng-ls-err2":        "logon.go Error Loading Password Hash 24",
	"esp-ls-err2":        "logon.go Error al Cargar el Hash de la Contraseña 24",
	"eng-ls-err3":        "logon.go Error Invalid Password",
	"esp-ls-err3":        "logon.go Error Contraseña no Válida",
	"eng-ls-err4":        "logon.go Error URL Incorrect Format",
	"esp-ls-err4":        "logon.go URL de Error Formato Incorrecto",
	"eng-ls-err5":        "Error Invalid Queue Password 24",
	"esp-ls-err5":        "Error Contraseña de Cola no Válida 24",
	"eng-ls-err6-1":      "Error Queue Password Length is ",
	"esp-ls-err6-1":      "La Longitud de la Contraseña de la Cola de Errores es ",
	"eng-ls-err6-2":      " should be length of 24",
	"esp-ls-err6-2":      " Debe Tener una Longitud de 24",
	"eng-ls-err7":        "No NATS connection",
	"esp-ls-err7":        "Sin Conexión NATS",
	"eng-ls-erase":       "Security Erase",
	"esp-ls-erase":       "Borrado de seguridad",
	"eng-ls-clogon":      "Communications Logon",
	"esp-ls-clogon":      "Inicio de Sesión de Comunicaciones",
	"eng-ls-err8":        "No JETSTREAM Connection",
	"esp-ls-err8":        "Sin Conexión JETSTREAM ",
	"eng-ms-title":       "Messages",
	"esp-ms-title":       "Mensajes",
	"eng-ms-mm":          "Enter Message For Encryption",
	"esp-ms-mm":          "Ingrese el Mensaje Para el Cifrado",
	"eng-ms-header1":     "Select An Item From The List",
	"esp-ms-header1":     "Seleccione un Elemento de la Lista",
	"eng-ms-err1":        "NATS No Connection",
	"esp-ms-err1":        "NATS sin Conexión",
	"eng-ms-sm":          "Send Message",
	"esp-ms-sm":          "Enviar Mensaje",
	"eng-ms-header2":     "NATS Messaging",
	"esp-ms-header2":     "Mensajería NATS",
	"eng-ms-rm":          "Receive Message",
	"esp-ms-rm":          "Recibir Mensaje",
	"eng-ms-err2":        "NATS No Connection",
	"esp-ms-err2":        "NATS sin Conexión",
	"eng-ms-err3":        "Could Not Add Consumer",
	"esp-ms-err3":        "No se Pudo Agregar el Consumidor",
	"eng-ms-err4":        "Error Pull Subscribe ",
	"esp-ms-err4":        "Error Extraer Suscribirse ",
	"eng-ms-err5":        "Error Fetch ",
	"esp-ms-err5":        "Recuperación de Errores ",
	"eng-ms-err6-1":      "Recieved ",
	"esp-ms-err6-1":      "Recibida ",
	"eng-ms-err6-2":      " Messages",
	"esp-ms-err6-2":      " Mensajes",
	"eng-ms-err7":        "Please Logon First",
	"esp-ms-err7":        "Por Favor Ingresa Primero",
	"eng-es-title":       "Enc/Dec",
	"esp-es-title":       "Codificar/Descodificar",
	"eng-es-pw":          "Enc/Dec",
	"esp-es-pw":          "Codificar/Descodificar",
	"eng-es-pass":        "Enter Password to Use For Encryption 24",
	"esp-es-pass":        "Ingrese la contraseña para usar para el cifrado 24",
	"eng-es-mv":          "Enter Value",
	"esp-es-mv":          "Introducir Valor",
	"eng-es-mo":          "Output Shows Up Here",
	"esp-es-mo":          "La Salida Aparece Aquí",
	"eng-es-em":          "Encrypt Message",
	"esp-es-em":          "Cifrar Mensaje",
	"eng-es-err1":        "Error Invalid Password",
	"esp-es-err1":        "Error Contraseña Inválida",
	"eng-es-err2-1":      "Error Password Length is ",
	"esp-es-err2-1":      "La Longitud de la Contraseña de Error es ",
	"eng-es-err2-2":      " Should be Length of 24",
	"esp-es-err2-2":      " Debe Tener una Longitud de 24",
	"eng-es-err3":        "Error Input Text",
	"esp-es-err3":        "Texto de Entrada de Error",
	"eng-es-err4":        "Cannot Encrypt Input Text",
	"esp-es-err4":        "No se Puede Cifrar el Texto de Entrada",
	"eng-es-err5":        "Cannot Decrypt Input Text",
	"esp-es-err5":        "No se Puede Descifrar el Texto de Entrada",
	"eng-es-head1":       "Input",
	"esp-es-head1":       "Aporte",
	"eng-es-head2":       "Output",
	"esp-es-head2":       "Producción",
	"eng-log-nc":         "Path to nats Configuration File \nFor TLS Certificate Paths",
}
var MyPanes = map[string]MyPane{}
var MyPanesIndex = map[string][]string{}

func GetLangs(mystring string) string {

	value, err := MyLangs[PreferedLanguage+"-"+mystring]
	//log.Println("GetLangs ", PreferedLanguage+"-"+mystring)
	if err == false {
		return "err"
	}
	return value
}

func Init() {
	MyPanes = map[string]MyPane{
		"password":     {GetLangs("ps-title"), "", passwordScreen, true},
		"settings":     {GetLangs("ss-title"), "", settingsScreen, true},
		"certificates": {GetLangs("cs-title"), "", certificatesScreen, true},
		"logon":        {GetLangs("ls-title"), "", logonScreen, true},
		"messages":     {GetLangs("ms-title"), "", messagesScreen, true},
		"encdec":       {GetLangs("es-title"), "", encdecScreen, true},
	}

	// PanesIndex  defines how our panes should be laid out in the index tree
	MyPanesIndex = map[string][]string{
		"": {"password", "logon", "settings", "certificates", "messages", "encdec"},
	}
}

func DataStore(myfile string) fyne.URI {
	DataLocation, dlerr := storage.Child(fyne.CurrentApp().Storage().RootURI(), myfile)
	if dlerr != nil {
		log.Println("DataStore error ", dlerr)
	}
	return DataLocation
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func SetMyApp(a fyne.App) {
	MyAppDup = a
}
func GetMyApp() fyne.App {
	return MyAppDup
}
func MyJson(action string) {
	if GetMyApp() == nil {
		MyApp := app.NewWithID("org.nh3000.snats")
		SetMyApp(MyApp)
	}
	MyApp = GetMyApp()
	if action == "LOAD" {
		// prepare fallback or just load

		PreferedLanguage = MyApp.Preferences().StringWithFallback("PreferedLanguage", "eng")

		xServer, _ := Encrypt("nats://nats.newhorizons3000.org:4222", MySecret)
		Server = MyApp.Preferences().StringWithFallback("Server", xServer)
		xQueue, _ := Encrypt("MESSAGES", MySecret)
		Queue = MyApp.Preferences().StringWithFallback("Queue", xQueue)
		xAlias, _ := Encrypt("MyAlias", MySecret)
		Alias = MyApp.Preferences().StringWithFallback("Alias", xAlias)
		xQueuepassword, _ := Encrypt("123456789012345678901234", MySecret)
		Queuepassword = MyApp.Preferences().StringWithFallback("Queuepasword", xQueuepassword)

		var xCaroot = strings.ReplaceAll("-----BEGIN CERTIFICATE-----\nMIICFDCCAbugAwIBAgIUDkHxHO1DwrlkTzUimG5PoiswB6swCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAgFw0yMzAzMzExNzI5MDBaGA8yMDUzMDMyMzE3MjkwMFowZjELMAkG\nA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UEChMDU0VD\nMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9y\nZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHXwMUfMXiJix3tuzFymcA+3RkeY\nZE7urUzVgaqkv/Oef3jhqhtf1XzK/qVYGxWWmpvADGB252PG1Mp7Z5wmzqyjRTBD\nMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMB0GA1UdDgQWBBQm\nFA5caanuqxGFOf9DtZkVYv5dCzAKBggqhkjOPQQDAgNHADBEAiB3BheNP4XdBZ27\nxVBQ7ztMJqK7wDi1V3LuMy5jmXr7rQIgHCse0oaiAwcl4VwF00aSshlV+T/da0Tx\n1ANkaM+rie4=\n-----END CERTIFICATE-----\n", "\n", "<>")
		ycaroot, _ := Encrypt(xCaroot, MySecret)
		Caroot = MyApp.Preferences().StringWithFallback("Caroot", ycaroot)

		var xClientcert = strings.ReplaceAll("-----BEGIN CERTIFICATE-----\nMIIDUzCCAvigAwIBAgIUUyhlJt8mp1XApRbSkdrUS55LGV8wCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAeFw0yMzAzMzExNzI5MDBaFw0yODAzMjkxNzI5MDBaMHIxCzAJBgNV\nBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlDcmVzdHZpZXcxGjAY\nBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMSEwHwYDVQQLExhuYXRzLm5ld2hvcml6\nb25zMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDFttVH\nQ131JYwazAQMm0XAQvRvTjTjOY3aei1++mmQ+NQ9mrOFk6HlZFoKqsy6+HPXsB9x\nQbWlYvUOuqBgb9xFQZoL8jiKskLLrXoIxUAlIBTlyf76r4SV+ZpxJYoGzXNTedaU\n0EMTyAiUQ6nBbFMXiehN5q8VzxtTESk7QguGdAUYXYsCmYBvQtBXoFYO5CHyhPqu\nOZh7PxRAruYypEWVFBA+29+pwVeaRHzpfd/gKLY4j2paInFn7RidYUTqRH97BjdR\nSZpOJH6fD7bI4L09pnFtII5pAARSX1DntS0nWIWhYYI9use9Hi/B2DRQLcDSy1G4\n0t1z4cdyjXxbFENTAgMBAAGjgawwgakwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQM\nMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFAzgPVB2/sfT7R0U\ne3iXRSvUkfoQMB8GA1UdIwQYMBaAFCYUDlxpqe6rEYU5/0O1mRVi/l0LMDQGA1Ud\nEQQtMCuCGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZ4IJMTI3LDAsMCwxhwTAqABn\nMAoGCCqGSM49BAMCA0kAMEYCIQCDlUH2j69mJ4MeKvI8noOmvLHfvP4qMy5nFW2F\nPT5UxgIhAL6pHFyEbANtSkcVJqxTyKE4GTXcHc4DB43Z1F7VxSJj\n-----END CERTIFICATE-----\n", "\n", "<>")
		yclientcert, _ := Encrypt(xClientcert, MySecret)
		Clientcert = MyApp.Preferences().StringWithFallback("Clientcert", yclientcert)

		var xClientkey = strings.ReplaceAll("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAxbbVR0Nd9SWMGswEDJtFwEL0b0404zmN2notfvppkPjUPZqz\nhZOh5WRaCqrMuvhz17AfcUG1pWL1DrqgYG/cRUGaC/I4irJCy616CMVAJSAU5cn+\n+q+ElfmacSWKBs1zU3nWlNBDE8gIlEOpwWxTF4noTeavFc8bUxEpO0ILhnQFGF2L\nApmAb0LQV6BWDuQh8oT6rjmYez8UQK7mMqRFlRQQPtvfqcFXmkR86X3f4Ci2OI9q\nWiJxZ+0YnWFE6kR/ewY3UUmaTiR+nw+2yOC9PaZxbSCOaQAEUl9Q57UtJ1iFoWGC\nPbrHvR4vwdg0UC3A0stRuNLdc+HHco18WxRDUwIDAQABAoIBACe0XMZP4Al//c/P\n0qxZbjt69q13jiVnhHYwfPx3+0UywySP8adMi4GOkop73Ftb05+n7diHspvA8KeB\nkP1s2VZLI01s2i/4NnPCpbQnMIeEFs5Cr2LWZpDbrEk2ma5eCd/kotQFssLBM//a\nSrfeMh2TA0TJo7WEft9Cnf4ZeEkKnycplfvwTyv286iFZCYo2dv66BfTej6kkVCo\nAi+ZVCe2zSqRYyr0u4/j/kE3b3eSkCnY2IVcqlP7epuEGVOZyxeFLwM5ljbWL816\npA6WIJgQo2EQ1N7L531neg5WjXQ/UwTQoXP1jvuuVtKtOBFqm1IshEyFk3WpsfpD\nr16OTdECgYEA6FB6NYxYtnWPaIYAOqP7GtMKoJujH8MtZy6J33LkxI7nPkMkn8Mv\nva32tvjU4Bu1FVNp9k5guC+b+8ixXK0URj25IOhDs6K57tck22W9WiTZlmnkCO01\nJOavrelWbvYt5xNWIdnPualoPfGB0iJKXsKY/bpH4eVfhWwpNPI5sMkCgYEA2d9G\nEPuWN6gUjZ+JfdS+0WHK1yGD7thXs7MPUlhGqDzBryh2dkywyo8U8+tMLuDok1RZ\njnT3PYkLQEpzoV0qBkpFFShL6ubaGmDz1UZsozl0YcIg4diZeuPHnIAeXOFrhgYf\n825163LmT3jYHCROFEMLtTYyIQP0EznE+qFT3TsCgYEApgtvbfqkJbWdDL5KR5+R\nCLky7VyQmVEtkIRI8zbxoDPrwCrJcI9X/iDrKBhuPshPA7EdGXkn1D3jJXFqo6zp\nwtK3EXgxe6Ghd766jz4Guvl/s+x3mpHA3GEtzAXtS14VrQW7GHLP8AnPggauHX14\n3oYER8XvPtxtC7YlNbyz01ECgYAe2b7SKM3ck7BVXYHaj4V1oKNYUyaba4b/qxtA\nTb+zkubaJqCfn7xo8lnFMExZVv+X3RnRUj6wN/ef4ur8rnSE739Yv5wAZy/7DD96\ns74uXrRcI2EEmechv59ESeACxuiy0as0jS+lZ1+1YSc41Os5c0T1I/d1NVoaXtPF\nqZJ2gQKBgBp/XavdULBPzC7B8tblySzmL01qJZV7MSSVo2/1vJ7gPM0nQPZdTDog\nTfA5QKSX9vFTGC9CZHSJ+fabYDDd6+3UNYUKINfr+kwu9C2cysbiPaM3H27WR5mW\n5LhStAfwuRRYBDsG2ndjraxcBrrPdtkbS0dpeQUDJxvkMIuLHnhQ\n-----END RSA PRIVATE KEY-----\n", "\n", "<>")
		yclientkey, _ := Encrypt(xClientkey, MySecret)
		Clientkey = MyApp.Preferences().StringWithFallback("Clientkey", yclientkey)

		var ymsgmaxage = []string{"12h", "24h", "161h", "8372h"}
		xmsgmaxage, _ := Encrypt(strings.Join(ymsgmaxage, ","), MySecret)
		Msgmaxage = MyApp.Preferences().StringWithFallback("Msgmaxage", xmsgmaxage)

		PasswordMinimumSize = MyApp.Preferences().StringWithFallback("PasswordMinimumSize", "6")
		PasswordMustContainNumber = MyApp.Preferences().StringWithFallback("PasswordMustContainNumber", "Yes")
		PasswordMustContainLetter = MyApp.Preferences().StringWithFallback("PasswordMustContainLetter", "Yes")
		PasswordMustContainSpecial = MyApp.Preferences().StringWithFallback("PasswordMustContainSpecial", "Yes")
		
		// prepare for operations
		yServer, _ := Decrypt(Server, MySecret)
		Server = yServer
		yMsgmaxage, _ := Decrypt(Msgmaxage, MySecret)
		Msgmaxage = yMsgmaxage
		yQueue, _ := Decrypt(Queue, MySecret)
		Queue = yQueue
		yAlias, _ := Decrypt(Alias, MySecret)
		Alias = yAlias
		yQueuepassword, _ := Decrypt(Queuepassword, MySecret)
		Queuepassword = yQueuepassword
		yCaroot, _ := Decrypt(Caroot, MySecret)
		Caroot = strings.ReplaceAll(yCaroot, "<>", "\n")
		yClientcert, _ := Decrypt(Clientcert, MySecret)
		Clientcert = strings.ReplaceAll(yClientcert, "<>", "\n")
		yClientkey, _ := Decrypt(Clientkey, MySecret)
		Clientkey = strings.ReplaceAll(yClientkey, "<>", "\n")
	}
	if action == "SAVE" {
		xCaroot, _ := Encrypt(Caroot, MySecret)
		MyApp.Preferences().SetString("Caroot", xCaroot)
		xClientcert, _ := Encrypt(Clientcert, MySecret)
		MyApp.Preferences().SetString("Clientcert", xClientcert)
		xClientkey, _ := Encrypt(Clientkey, MySecret)
		MyApp.Preferences().SetString("Clientkey", xClientkey)

		xMsgmaxage, _ := Encrypt(Msgmaxage, MySecret)
		MyApp.Preferences().SetString("Msgmaxage", xMsgmaxage)

		xServer, _ := Encrypt(Server, MySecret)
		MyApp.Preferences().SetString("Server", xServer)
		xQueue, _ := Encrypt(Queue, MySecret)
		MyApp.Preferences().SetString("Queue", xQueue)
		xAlias, _ := Encrypt(Alias, MySecret)
		MyApp.Preferences().SetString("Alias", xAlias)
		MyApp.Preferences().SetString("PreferedLanguage", PreferedLanguage)
		xQueuepassword, _ := Encrypt(Queuepassword, MySecret)
		MyApp.Preferences().SetString("Queuepassword", xQueuepassword)
		MyApp.Preferences().SetString("PasswordMinimumSize", PasswordMinimumSize)
		MyApp.Preferences().SetString("PasswordMustContainNumber", PasswordMustContainNumber)
		MyApp.Preferences().SetString("PasswordMustContainLetter", PasswordMustContainLetter)
		MyApp.Preferences().SetString("PasswordMustContainSpecial", PasswordMustContainSpecial)
	}

}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println("DataStore Decode ", s, "  ", err)
		log.Panic(err)
	}
	return data
}

func MyCrypt(action string) {

	if action == "ENCRYPT" {

		cryptoText, _ := Encrypt(Server, MySecret)

		Server = cryptoText

		cryptoText4, _ := Encrypt(Queue, MySecret)
		Queue = cryptoText4

		cryptoText5, _ := Encrypt(Queuepassword, MySecret)
		Queuepassword = cryptoText5

		cryptoText8, _ := Encrypt(Alias, MySecret)
		Alias = cryptoText8

		cryptoText9, _ := Encrypt(NodeUUID, MySecret)
		NodeUUID = cryptoText9

	}
	if action == "DECRYPT" {
		text, _ := Decrypt(Server, MySecret)
		Server = text

		text4, _ := Decrypt(Queue, MySecret)
		Queue = text4
		text5, _ := Decrypt(Queuepassword, MySecret)
		Queuepassword = text5

		text8, _ := Decrypt(Alias, MySecret)
		Alias = text8
		text9, _ := Decrypt(NodeUUID, MySecret)
		NodeUUID = text9

	}
}

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

func MyHash(action string) bool {

	if action == "CREATE" {

		err := storage.Delete(DataStore("config.hash"))
		if err != nil {
			log.Println("MyHash Error Deleting ", DataStore("config.hash"))
		}
		wrt, errwrite := storage.Writer(DataStore("config.hash"))
		_, err2 := wrt.Write([]byte(Passwordhash))

		if errwrite != nil || err2 != nil {
			log.Println("MyHash Error Writing ", DataStore("config.hash"))
			return true
		}

	}
	if action == "LOAD" {
		ph, errf := os.ReadFile(DataStore("config.hash").Path())
		Passwordhash = string(ph)
		if errf != nil {
			log.Println("MyHash LOAD Hash Error File ", errf, " ", DataStore("config.hash"))
			return true
		}

	}
	if action == "SAVE" {

		errf := storage.Delete(DataStore("config.hash"))

		if errf != nil {
			log.Println("MyHash SAVE Hash Error file", errf)
			return true
		}
		wrt, errwrite := storage.Writer(DataStore("config.hash"))
		_, err2 := wrt.Write([]byte(Passwordhash))
		if errwrite != nil || err2 != nil {
			log.Println("MyHash SAVE Error Writing", DataStore("config.hash"))
			return true
		}

	}
	return false
}

func NATSErase() {
	log.Println("Erasing  ")
	//msgmaxage, _ := time.ParseDuration("148h")
	msgmaxage, _ := time.ParseDuration(Msgmaxage)
	nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
	if err != nil {
		log.Println("NatsErase Connection ", err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("NatsErase JetStream ", err)
	}

	NatsMessages = nil

	err1 := js.DeleteStream(Queue)
	if err != nil {
		log.Println("NatsErase DeleteStream ", err1)
	}

	js1, err1 := js.AddStream(&nats.StreamConfig{
		Name:     Queue,
		Subjects: []string{strings.ToLower(Queue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})

	if err1 != nil {
		log.Println("NatsErase AddStream ", err1)
	}
	fmt.Printf("js1: %v\n", js1)

	ac, err1 := js.AddConsumer(Queue, &nats.ConsumerConfig{
		Durable:       MyDurable,
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverAllPolicy,
		ReplayPolicy:  nats.ReplayInstantPolicy,
	})
	if err1 != nil {
		log.Println("NatsErase AddConsumer ", err1, " ", ac)
	}

	js.Publish(strings.ToLower(Queue)+"."+NodeUUID, []byte(FormatMessage("Security Erase")))

	nc.Close()

}

func FormatMessage(m string) []byte {
	EncMessage := MessageStore{}

	//ID , err := exec.Command("uuidgen").Output()

	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\nNo Host Name"
		//strings.Replace(EncMessage, "#HOSTNAME", "No Host Name", -1)

	} else {
		EncMessage.MShostname = "\nHost - " + name
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
		EncMessage.MShostname += "\nMac ids"
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\nAddress"
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}

	}

	EncMessage.MSalias = Alias

	EncMessage.MSnodeuuid = "\nNode Id - " + NodeUUID[0:8]
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\nMessage Id - " + iduuid[0:8]
	EncMessage.MSdate = "\nOn -" + time.Now().Format(time.UnixDate)
	//EncMessage.MSdate = "\nOn -"
	EncMessage.MSmessage = m
	//EncMessage += m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	ejson, _ := Encrypt(string(jsonmsg), Queuepassword)
	return []byte(ejson)

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
		valid := strings.Contains(strings.ToLower(value), "nats://")
		if valid == false {
			return true
		}
		valid2 := strings.Contains(value, ".")
		if valid2 == false {
			return true
		}
		valid3 := strings.Contains(value, ":")
		if valid3 == false {
			return true
		}

		return false
	}
	if action == "STRING" {

		if len(value) == 0 {
			return true
		}
		return false
	}

	if action == "PASSWORD" {
		var iserrors = false
		vlen, _ := strconv.Atoi(PasswordMinimumSize)
		if (len(value) <= vlen) == false {
			iserrors = true
		}

		if PasswordMustContainLetter == "Yes" && !iserrors {

			for _, r := range value {
				if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
					iserrors = true
					break
				}
			}
		}

		if PasswordMustContainNumber == "Yes" && !iserrors {
			iserrors = true
			for _, r := range value {
				if unicode.IsNumber(r) {
					iserrors = false
					break
				}
			}
		}
		if PasswordMustContainSpecial == "Yes" && !iserrors {
			iserrors = true
			var schars = []string{"|", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "-", "+", "=", "{", "}", "]", "[", "|", ":", ";", ",", ".", "#", "'", "\"", "\\", "%", "?", "\n", "<", "Ø", "ð", ">", "ï", "û"}
			for _, sc := range schars {
				if strings.Contains(value, sc) {
					iserrors = false
					break
				}
			}
		}
		return iserrors
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
