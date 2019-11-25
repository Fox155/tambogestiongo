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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	urlTambo string = "http://api.tambogestion.ga"
)

func main() {
	file, err := os.Open("test.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	lineasBytes := make([][]byte, len(txtlines))

	for i, eachline := range txtlines {
		lineasBytes[i] = []byte(eachline)
		fmt.Println(eachline)
		res1 := strings.Split(eachline, "/")
		fmt.Println("\tSesion de Ordeño: ", res1[0])
		fmt.Println("\tRFID: ", res1[1])
		fmt.Println("\tProduccion: ", res1[2])
		fmt.Println("\tMedidor: ", res1[3])
		fmt.Println("\tFecha Inicio: ", res1[4])
		fmt.Println("\tFecha Fin: ", res1[5])
	}

	privada, errP := fileToPrivateKey()
	if errP != nil {
		log.Panic(errP)
	}
	fmt.Println("Privada:\t", privada)

	fmt.Println("Linea:\t", lineasBytes[0])
	resultado := EncryptWithPublicKey(lineasBytes[0], &privada.PublicKey)
	fmt.Println("Resultado Bytes:\t", resultado)
	fmt.Println("Resultado Bytes String:\t", string(resultado))

	url := "http://localhost:8080/producciones"

	requestBody, errJ := json.Marshal(map[string]interface{}{
		"Tambo":   "Fox",
		"Mensaje": resultado,
	})
	if errJ != nil {
		log.Panic(errJ)
	}

	res, errH := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if errH != nil {
		log.Panic(errH)
	}

	// res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
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

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha1.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Panicln(err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha1.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Panicln(err)
	}
	return plaintext
}
