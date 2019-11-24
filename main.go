package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// file, err := os.Open("test.txt")

	// if err != nil {
	// 	log.Fatalf("failed opening file: %s", err)
	// }

	// scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanLines)
	// var txtlines []string

	// for scanner.Scan() {
	// 	txtlines = append(txtlines, scanner.Text())
	// }

	// file.Close()

	// for _, eachline := range txtlines {
	// 	fmt.Println(eachline)
	// 	res1 := strings.Split(eachline, "/")
	// 	fmt.Println("\tSesion de Ordeño: ", res1[0])
	// 	fmt.Println("\tRFID: ", res1[1])
	// 	fmt.Println("\tProduccion: ", res1[2])
	// 	fmt.Println("\tMedidor: ", res1[3])
	// 	fmt.Println("\tFecha Inicio: ", res1[4])
	// 	fmt.Println("\tFecha Fin: ", res1[5])
	// }

	url := "http://tambogestion.ga"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cookie", "advanced-backend=m87tgnu6itsj5kvlepqt9p58mq")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
