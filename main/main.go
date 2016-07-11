package main

import (
	"fmt"
	ml "github.com/cube2222/SimpleActivityRecognitionNeuralNetwork"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	testCount, err := strconv.Atoi(os.Args[1])
	trainingCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	rand.Seed(time.Now().Unix())
	myNeuralNetworkTraining := ml.NewTraining(2000, 1000, 10)
	myNeuralNetworkTraining.InitRandom()
	data, _ := myNeuralNetworkTraining.SaveToJson()
	myNeuralNetwork := ml.New(2000, 1000, 10)
	myNeuralNetwork.InitFromJson(data)
	inputs := make([][]float64, testCount)
	outputs := make([]float64, testCount)
	for i := 0; i < testCount; i++ {
		inputs[i] = make([]float64, 2000)
	}
	for i := 0; i < testCount; i++ {
		for j := 0; j < 2000; j++ {
			inputs[i][j] = rand.Float64()
		}
		outputs[i] = float64(rand.Int() % 2)
	}
	goodOnes := 0
	for i := 0; i < testCount; i++ {
		if myNeuralNetwork.Process(inputs[i]) == outputs[i] {
			goodOnes++
		}
	}
	fmt.Println("Got", goodOnes, "out of", testCount)

	for i := 0; i < testCount; i++ {
		for j := 0; j < trainingCount; j++ {
			myNeuralNetworkTraining.Train(inputs[i], outputs[i], 0.1)
		}
		fmt.Println("Trained", i+1)
	}
	data, _ = myNeuralNetworkTraining.SaveToJson()

	myNeuralNetwork.InitFromJson(data)
	goodOnes = 0
	for i := 0; i < testCount; i++ {
		if myNeuralNetwork.Process(inputs[i]) == outputs[i] {
			goodOnes++
		}
	}
	fmt.Println("Got", goodOnes, "out of", testCount)
}
