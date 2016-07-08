package SimpleActivityRecognitionNeuralNetwork

import (
	"encoding/json"
)

type NeuralNetworkTraining struct {
	layerOut       Neuron
	layerOuter     []Neuron
	layerMid       []Neuron
	layerBot       []Neuron
	bottomInputs   []chan float64
	output         chan float64
	curOutputs     []chan float64
	curOutputSaved []float64
}

func getMax(a int, b int, c int) int {
	d := a
	if b > d {
		d = b
	}
	if c > d {
		d = c
	}
	return d
}

func NewTraining(botNum int, midNum int, outerNum int) *NeuralNetworkTraining {
	myNeuralNetwork := NeuralNetworkTraining{}
	myNeuralNetwork.layerBot = make([]Neuron, botNum)
	myNeuralNetwork.layerMid = make([]Neuron, midNum)
	myNeuralNetwork.layerOuter = make([]Neuron, outerNum)

	myNeuralNetwork.curOutputs = make([]chan float64, getMax(botNum, midNum, outerNum))
	for i := 0; i < len(myNeuralNetwork.curOutputs); i++ {
		myNeuralNetwork.curOutputs[i] = make(chan float64)
	}

	// Create bottom inputs
	for i := 0; i < botNum; i++ {
		myNeuralNetwork.layerBot[i].outputReceivers = 1
		myNeuralNetwork.layerBot[i].weights = make([]float64, 1)
		myNeuralNetwork.layerBot[i].inputs = make([]chan float64, 1)
		for j := 0; j < len(myNeuralNetwork.layerBot[i].inputs); j++ {
			myNeuralNetwork.layerBot[i].inputs[j] = make(chan float64)
		}
		myNeuralNetwork.layerBot[i].output = make(chan float64)
	}

	// Create mid layer
	for i := 0; i < midNum; i++ {
		myNeuralNetwork.layerMid[i].outputReceivers = 1
		myNeuralNetwork.layerMid[i].weights = make([]float64, botNum)
		myNeuralNetwork.layerMid[i].inputs = make([]chan float64, botNum)
		for j := 0; j < len(myNeuralNetwork.layerMid[i].inputs); j++ {
			myNeuralNetwork.layerMid[i].inputs[j] = make(chan float64)
		}
		myNeuralNetwork.layerMid[i].output = make(chan float64)
	}

	// Create outer layer
	for i := 0; i < outerNum; i++ {
		myNeuralNetwork.layerOuter[i].outputReceivers = 1
		myNeuralNetwork.layerOuter[i].weights = make([]float64, midNum)
		myNeuralNetwork.layerOuter[i].inputs = make([]chan float64, midNum)
		for j := 0; j < len(myNeuralNetwork.layerOuter[i].inputs); j++ {
			myNeuralNetwork.layerOuter[i].inputs[j] = make(chan float64)
		}
		myNeuralNetwork.layerOuter[i].output = make(chan float64)
	}

	// Wire outer to out
	myNeuralNetwork.layerOut.outputReceivers = 1
	myNeuralNetwork.layerOut.weights = make([]float64, outerNum)
	myNeuralNetwork.layerOut.inputs = make([]chan float64, outerNum)
	for i := 0; i < len(myNeuralNetwork.layerOut.inputs); i++ {
		myNeuralNetwork.layerOut.inputs[i] = make(chan float64)
	}

	myNeuralNetwork.layerOut.output = make(chan float64)

	return &myNeuralNetwork
}

func (myNeuralNetwork *NeuralNetworkTraining) InitRandom() {
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		myNeuralNetwork.layerBot[i].Randomize()
	}

	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		myNeuralNetwork.layerMid[i].Randomize()
	}

	for i := 0; i < len(myNeuralNetwork.layerOuter); i++ {
		myNeuralNetwork.layerOuter[i].Randomize()
	}

	myNeuralNetwork.layerOut.Randomize()

}

func (myNeuralNetwork *NeuralNetworkTraining) InitFromJson(data []byte) error {
	myNeuralNetworkModel := NeuralNetworkTrainedModel{}
	err := json.Unmarshal(data, &myNeuralNetworkModel)
	if err != nil {
		return err
	}

	for i := 0; i < len(myNeuralNetworkModel.LayerBot); i++ {
		myNeuralNetwork.layerBot[i].weights = myNeuralNetworkModel.LayerBot[i].Weights
		myNeuralNetwork.layerBot[i].bias = myNeuralNetworkModel.LayerBot[i].Bias
	}

	for i := 0; i < len(myNeuralNetworkModel.LayerMid); i++ {
		myNeuralNetwork.layerMid[i].weights = myNeuralNetworkModel.LayerMid[i].Weights
		myNeuralNetwork.layerMid[i].bias = myNeuralNetworkModel.LayerMid[i].Bias
	}

	for i := 0; i < len(myNeuralNetworkModel.LayerOuter); i++ {
		myNeuralNetwork.layerOuter[i].weights = myNeuralNetworkModel.LayerOuter[i].Weights
		myNeuralNetwork.layerOuter[i].bias = myNeuralNetworkModel.LayerOuter[i].Bias
	}

	myNeuralNetwork.layerOut.weights = myNeuralNetworkModel.LayerOut.Weights
	myNeuralNetwork.layerOut.bias = myNeuralNetworkModel.LayerOut.Bias

	return nil
}

func (myNeuralNetwork *NeuralNetworkTraining) SaveToJson() ([]byte, error) {
	myNeuralNetworkModel := NeuralNetworkTrainedModel{}

	myNeuralNetworkModel.LayerBot = make([]NeuronTrained, len(myNeuralNetwork.layerBot))
	myNeuralNetworkModel.LayerMid = make([]NeuronTrained, len(myNeuralNetwork.layerMid))

	for i := 0; i < len(myNeuralNetworkModel.LayerBot); i++ {
		myNeuralNetworkModel.LayerBot[i].Weights = myNeuralNetwork.layerBot[i].weights
		myNeuralNetworkModel.LayerBot[i].Bias = myNeuralNetwork.layerBot[i].bias
	}

	for i := 0; i < len(myNeuralNetworkModel.LayerMid); i++ {
		myNeuralNetworkModel.LayerMid[i].Weights = myNeuralNetwork.layerMid[i].weights
		myNeuralNetworkModel.LayerMid[i].Bias = myNeuralNetwork.layerMid[i].bias
	}

	for i := 0; i < len(myNeuralNetworkModel.LayerOuter); i++ {
		myNeuralNetworkModel.LayerOuter[i].Weights = myNeuralNetwork.layerOuter[i].weights
		myNeuralNetworkModel.LayerOuter[i].Bias = myNeuralNetwork.layerOuter[i].bias
	}

	myNeuralNetworkModel.LayerOut.Weights = myNeuralNetwork.layerOut.weights
	myNeuralNetworkModel.LayerOut.Bias = myNeuralNetwork.layerOut.bias

	data, err := json.Marshal(&myNeuralNetworkModel)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (myNeuralNetwork *NeuralNetworkTraining) Train(input []float64, output float64, learningRate float64) {

	myNeuralNetwork.curOutputSaved = make([]float64, len(myNeuralNetwork.layerBot))

	// Train bottom layer
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		myNeuralNetwork.layerBot[i].output = myNeuralNetwork.curOutputs[i]
		go myNeuralNetwork.layerBot[i].Process()
		myNeuralNetwork.layerBot[i].inputs[0] <- input[i]
	}
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		myNeuralNetwork.curOutputSaved[i] = <-myNeuralNetwork.curOutputs[i]
	}
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		myNeuralNetwork.layerBot[i].Adjust(input[i:i+1], output-myNeuralNetwork.curOutputSaved[i], learningRate)
	}

	// Train middle layer

	input = input[0:len(myNeuralNetwork.layerBot)]
	for i := 0; i < len(input); i++ {
		input[i] = myNeuralNetwork.curOutputSaved[i]
	}
	myNeuralNetwork.curOutputSaved = myNeuralNetwork.curOutputSaved[0:len(myNeuralNetwork.layerMid)]

	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		myNeuralNetwork.layerMid[i].output = myNeuralNetwork.curOutputs[i]
		go myNeuralNetwork.layerMid[i].Process()
		for j, curChan := range myNeuralNetwork.layerMid[i].inputs {
			curChan <- input[j]
		}
	}
	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		myNeuralNetwork.curOutputSaved[i] = <-myNeuralNetwork.curOutputs[i]
	}
	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		myNeuralNetwork.layerMid[i].Adjust(input, output-myNeuralNetwork.curOutputSaved[i], learningRate)
	}

	// Train outer layer

	input = input[0:len(myNeuralNetwork.layerMid)]
	for i := 0; i < len(input); i++ {
		input[i] = myNeuralNetwork.curOutputSaved[i]
	}
	myNeuralNetwork.curOutputSaved = myNeuralNetwork.curOutputSaved[0:len(myNeuralNetwork.layerOuter)]

	for i := 0; i < len(myNeuralNetwork.layerOuter); i++ {
		myNeuralNetwork.layerOuter[i].output = myNeuralNetwork.curOutputs[i]
		go myNeuralNetwork.layerOuter[i].Process()
		for j, curChan := range myNeuralNetwork.layerOuter[i].inputs {
			curChan <- input[j]
		}
	}
	for i := 0; i < len(myNeuralNetwork.layerOuter); i++ {
		myNeuralNetwork.curOutputSaved[i] = <-myNeuralNetwork.curOutputs[i]
	}
	for i := 0; i < len(myNeuralNetwork.layerOuter); i++ {
		myNeuralNetwork.layerOuter[i].Adjust(input, output-myNeuralNetwork.curOutputSaved[i], learningRate)
	}

	// Train out layer

	input = input[0:len(myNeuralNetwork.layerOuter)]
	for i := 0; i < len(input); i++ {
		input[i] = myNeuralNetwork.curOutputSaved[i]
	}

	myNeuralNetwork.layerOut.output = myNeuralNetwork.curOutputs[0]
	go myNeuralNetwork.layerOut.Process()
	for j, curChan := range myNeuralNetwork.layerOut.inputs {
		curChan <- input[j]
	}

	myNeuralNetwork.curOutputSaved[0] = <-myNeuralNetwork.curOutputs[0]

	myNeuralNetwork.layerOut.Adjust(input, output-myNeuralNetwork.curOutputSaved[0], learningRate)
}
