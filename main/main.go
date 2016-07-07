package main

import (
	"bytes"
	"fmt"
	ml "github.com/cube2222/SimpleActivityRecognitionNeuralNetwork"
	"io"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	myNeuralNetwork := ml.New(2000, 1000, 10)
	myNeuralNetwork.InitRandom()
	//file, _ := os.Open("C:/Development/Projects/Go/src/github.com/cube2222/SimpleActivityRecognitionNeuralNetwork/main/test.json")
	//data, _ := ioutil.ReadAll(file)
	//myNeuralNetwork.InitFromJson(data)
	//rand.Seed(time.Now().Unix())
	file, _ := os.Create("C:/Development/Projects/Go/src/github.com/cube2222/SimpleActivityRecognitionNeuralNetwork/main/test.json")
	data, _ := myNeuralNetwork.SaveToJson()
	io.Copy(file, bytes.NewBuffer(data))
	file.Close()
	tab := make([]float64, 2000)
	for i := 0; i < 2000; i++ {
		tab[i] = rand.Float64()
	}
	fmt.Println(myNeuralNetwork.Process(tab))
}
