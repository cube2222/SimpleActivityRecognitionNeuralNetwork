package main

import (
	"fmt"
	ml "github.com/cube2222/SimpleActivityRecognitionNeuralNetwork"
	"io/ioutil"
	"math/rand"
	"os"
)

func main() {
	myNeuralNetwork := ml.New(10, 2000)
	file, _ := os.Open("C:/Development/Projects/Go/src/github.com/cube2222/SimpleActivityRecognitionNeuralNetwork/main/test.json")
	data, _ := ioutil.ReadAll(file)
	myNeuralNetwork.InitFromJson(data)
	//rand.Seed(time.Now().Unix())
	tab := make([]float64, 2000)
	for i := 0; i < 2000; i++ {
		tab[i] = rand.Float64()
	}
	fmt.Println(myNeuralNetwork.Process(tab))
}
