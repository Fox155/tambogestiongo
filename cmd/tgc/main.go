package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	urlTambo string = "http://localhost:5000/producciones" // "http://api.tambogestion.ga"
	tambo    string = "Tambo Ejemplo"
)

var (
	keyFile = flag.String("key", "id_rsa", "Path to RSA private key")
)

func main() {
	file, err := os.Open("./test.txt")
	fileTemp, errt := os.Create("./test_temp.txt")

	if err != nil || errt != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer fileTemp.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	writer := bufio.NewWriter(fileTemp)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
		log.Println(scanner.Text())
		res1 := strings.Split(scanner.Text(), "/")
		ids, _ := strconv.ParseInt(res1[0], 10, 64)
		idr, _ := strconv.ParseInt(res1[1], 10, 64)
		prof, _ := strconv.ParseFloat(res1[2], 64)
		medidor := map[string]string{
			"Nombre": res1[3],
		}
		inicio, _ := time.Parse(time.RFC3339, res1[4])
		fin, _ := time.Parse(time.RFC3339, res1[5])
		prod := Producciones{
			IDSesionOrdeño: ids,
			IDRFID:         idr,
			Produccion:     prof,
			Medidor:        medidor,
			FechaInicio:    inicio,
			FechaFin:       fin,
		}
		log.Println("Produccion:\t", prod)

		if errEP := enviarProduccion(prod); errEP != nil {
			fileTemp.Sync()
			writer.WriteString(fmt.Sprintln(prod))
			writer.Flush()
		}
	}

	defer file.Close()

	// os.Create("./humo.txt")
	// if errR := os.Remove("./humo.txt"); errR != nil {
	// 	log.Fatal(errR)
	// }
}

// enviarProduccion Encrypta y envia la una produccion
func enviarProduccion(produccion Producciones) error {
	privada, errP := fileToPrivateKey()
	if errP != nil {
		log.Panic(errP)
	}

	jsonProd, errj := json.Marshal(produccion)
	if errj != nil {
		log.Fatal("Decryp Error:", errj)
	}

	conf, err := initConfig()
	if err != nil {
		panic(err)
	}

	resultado := encryptWithPublicKey(jsonProd, &privada.PublicKey)

	requestBody, errJ := json.Marshal(map[string]interface{}{
		"Tambo":     tambo,
		"Sucursal":  conf.Sucursal.Nombre,
		"Contenido": resultado,
	})
	if errJ != nil {
		log.Panic(errJ)
	}

	res, errH := http.Post(urlTambo, "application/json", bytes.NewBuffer(requestBody))
	if errH != nil {
		log.Panic(errH)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	men := mensaje{}
	json.Unmarshal(body, &men)
	log.Println(res)

	if men.Mensaje != "OK" {
		return errors.New("Error al enviar la Produccion")
	}
	return nil
}

// fileToPrivateKey bytes to private key
func fileToPrivateKey() (*rsa.PrivateKey, error) {
	flag.Parse()
	// Read the private key
	priv, errR := ioutil.ReadFile(*keyFile)
	if errR != nil {
		return nil, errR
	}

	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// encryptWithPublicKey encrypts data with public key
func encryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha1.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Panicln(err)
	}
	return ciphertext
}

// Producciones estructura modelo de una produccion
type Producciones struct {
	IDSesionOrdeño int64
	IDRFID         int64
	Produccion     float64
	FechaInicio    time.Time
	FechaFin       time.Time
	Medidor        map[string]string
}

// Producciones estructura modelo de una produccion
type mensaje struct {
	Mensaje string `json:"Mensaje"`
}

// 1/1003/101.1/M0/2016-09-01T10:11:12Z/2016-09-01T10:15:25Z
// 1/1003/102.2/M0/2016-09-01T10:11:12Z/2016-09-01T10:15:25Z
// 1/1003/103.3/M0/2016-09-01T10:11:12Z/2016-09-01T10:15:25Z
