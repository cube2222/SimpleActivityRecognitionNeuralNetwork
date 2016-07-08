package main

import (
	"fmt"
	ml "github.com/cube2222/SimpleActivityRecognitionNeuralNetwork"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	myNeuralNetworkTraining := ml.NewTraining(2000, 1000, 10)
	myNeuralNetworkTraining.InitRandom()
	data, _ := myNeuralNetworkTraining.SaveToJson()
	myNeuralNetwork := ml.New(2000, 1000, 10)
	myNeuralNetwork.InitFromJson(data)
	inputs := make([][]float64, 10)
	outputs := make([]float64, 10)
	for i := 0; i < 10; i++ {
		inputs[i] = make([]float64, 2000)
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 2000; j++ {
			inputs[i][j] = rand.Float64()
		}
		outputs[i] = float64(rand.Int() % 2)
	}
	goodOnes := 0
	for i := 0; i < 10; i++ {
		if myNeuralNetwork.Process(inputs[i]) == outputs[i] {
			goodOnes++
		}
	}
	fmt.Println("Got", goodOnes, "out of 10.")

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			myNeuralNetworkTraining.Train(inputs[i], outputs[i], 0.1)
		}
		fmt.Println("Trained", i+1)
	}
	data, _ = myNeuralNetworkTraining.SaveToJson()

	myNeuralNetwork.InitFromJson(data)
	goodOnes = 0
	for i := 0; i < 10; i++ {
		if myNeuralNetwork.Process(inputs[i]) == outputs[i] {
			goodOnes++
		}
	}
	fmt.Println("Got", goodOnes, "out of 10.")
}
