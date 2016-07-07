package SimpleActivityRecognitionNeuralNetwork

import (
	"encoding/json"
)

type NeuralNetwork struct {
	layerOut     Neuron
	layerOuter   []Neuron
	layerMid     []Neuron
	layerBot     []Neuron
	bottomInputs []chan float64
	output       chan float64
}

func New(botNum int, midNum int, outerNum int) *NeuralNetwork {
	myNeuralNetwork := NeuralNetwork{}
	myNeuralNetwork.layerBot = make([]Neuron, botNum)
	myNeuralNetwork.layerMid = make([]Neuron, midNum)
	myNeuralNetwork.layerOuter = make([]Neuron, outerNum)
	myNeuralNetwork.bottomInputs = make([]chan float64, botNum)
	for i := 0; i < botNum; i++ { // Create Neural Network inputs
		myNeuralNetwork.bottomInputs[i] = make(chan float64)
	}

	// Create bottom inputs
	for i := 0; i < botNum; i++ {
		myNeuralNetwork.layerBot[i].outputReceivers = midNum
		myNeuralNetwork.layerBot[i].weights = make([]float64, 1)
		myNeuralNetwork.layerBot[i].inputs = make([]chan float64, 1)
		myNeuralNetwork.layerBot[i].inputs[0] = myNeuralNetwork.bottomInputs[i]
		myNeuralNetwork.layerBot[i].output = make(chan float64)
	}

	// Create mid layer
	for i := 0; i < midNum; i++ {
		myNeuralNetwork.layerMid[i].outputReceivers = outerNum
		myNeuralNetwork.layerMid[i].weights = make([]float64, botNum)
		myNeuralNetwork.layerMid[i].inputs = make([]chan float64, botNum)
		myNeuralNetwork.layerMid[i].output = make(chan float64)
	}
	for i := 0; i < midNum; i++ {
		for j := 0; j < botNum; j++ {
			myNeuralNetwork.layerMid[i].inputs[j] = myNeuralNetwork.layerBot[j].output
		}
	}

	// Create outer layer
	for i := 0; i < outerNum; i++ {
		myNeuralNetwork.layerOuter[i].outputReceivers = 1
		myNeuralNetwork.layerOuter[i].weights = make([]float64, midNum)
		myNeuralNetwork.layerOuter[i].inputs = make([]chan float64, midNum)
		myNeuralNetwork.layerOuter[i].output = make(chan float64)
	}
	for i := 0; i < outerNum; i++ {
		for j := 0; j < midNum; j++ {
			myNeuralNetwork.layerOuter[i].inputs[j] = myNeuralNetwork.layerMid[j].output
		}
	}

	// Wire outer to out
	myNeuralNetwork.layerOut.outputReceivers = 1
	myNeuralNetwork.layerOut.weights = make([]float64, outerNum)
	myNeuralNetwork.layerOut.inputs = make([]chan float64, outerNum)
	for i := 0; i < outerNum; i++ {
		myNeuralNetwork.layerOut.inputs[i] = myNeuralNetwork.layerOuter[i].output
	}

	myNeuralNetwork.layerOut.output = make(chan float64)

	myNeuralNetwork.output = myNeuralNetwork.layerOut.output
	return &myNeuralNetwork
}

func (myNeuralNetwork *NeuralNetwork) InitRandom() {
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

func (myNeuralNetwork *NeuralNetwork) InitFromJson(data []byte) error {
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

func (myNeuralNetwork *NeuralNetwork) SaveToJson() ([]byte, error) {
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

func (myNeuralNetwork *NeuralNetwork) Process(inputs []float64) float64 {
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		go myNeuralNetwork.layerBot[i].Process()
	}
	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		go myNeuralNetwork.layerMid[i].Process()
	}
	for i := 0; i < len(myNeuralNetwork.layerOuter); i++ {
		go myNeuralNetwork.layerOuter[i].Process()
	}
	go myNeuralNetwork.layerOut.Process()
	for i, item := range myNeuralNetwork.bottomInputs {
		item <- inputs[i]
	}

	return <-myNeuralNetwork.output
}

func (myNeuralNetwork *NeuralNetwork) Train(inputs []float64, expectedOutput float64, learningRate float64) {

}
