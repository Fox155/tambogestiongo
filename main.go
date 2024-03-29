package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/gob"
	"encoding/json"
	"encoding/pem"
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

const urlTambo string = "http://api.tambogestion.ga"

var (
	keyFile = flag.String("key", "id_rsa", "Path to RSA private key")
)

func main() {
	file, err := os.Open("./test.txt")

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

	produc := []Producciones{}

	for i, eachline := range txtlines {
		lineasBytes[i] = []byte(eachline)
		res1 := strings.Split(eachline, "/")
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
		produc = append(produc, prod)
	}

	log.Println("Producciones:\t\t", produc)

	privada, errP := fileToPrivateKey()
	if errP != nil {
		log.Panic(errP)
	}
	// fmt.Println("Privada:\t", privada)

	fmt.Println("Linea:\t", produc[0])

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	// Encode (send) the value.
	errE := enc.Encode(produc[0])
	if errE != nil {
		log.Fatal("encode error:", errE)
	}
	// HERE ARE YOUR BYTES!!!!
	fmt.Println("Encode:\t", network.Bytes())

	var q Producciones
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Println("Decode:\t", q)

	resultado := encryptWithPublicKey(network.Bytes(), &privada.PublicKey)
	fmt.Println("Resultado Bytes:\t", resultado)
	fmt.Println("Resultado Bytes String:\t", string(resultado))

	url := "http://localhost:5000/producciones"

	requestBody, errJ := json.Marshal(map[string]interface{}{
		"Tambo":     "Fox",
		"Contenido": resultado,
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

	// fmt.Printf("Current Unix Time: %v\n", time.Now())

	// time.Sleep(time.Minute)

	// fmt.Printf("Current Unix Time: %v\n", time.Now())
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
	NroLactancia   int
	Produccion     float64
	FechaInicio    time.Time
	FechaFin       time.Time
	Medidor        map[string]string
}
