package SimpleActivityRecognitionNeuralNetwork

import (
	"encoding/json"
)

type NeuralNetwork struct {
	layerOut     Neuron
	layerMid     []Neuron
	layerBot     []Neuron
	bottomInputs []chan float64
	output       chan float64
}

func New(midNum int, botNum int) *NeuralNetwork {
	myNeuralNetwork := NeuralNetwork{}
	myNeuralNetwork.layerBot = make([]Neuron, botNum)
	myNeuralNetwork.layerMid = make([]Neuron, midNum)
	myNeuralNetwork.bottomInputs = make([]chan float64, botNum)
	for i := 0; i < botNum; i++ { //Create Neural Network inputs
		myNeuralNetwork.bottomInputs[i] = make(chan float64)
	}

	for i := 0; i < midNum; i++ {
		myNeuralNetwork.layerMid[i].inputs = make([]chan float64, botNum/midNum)
	}
	for i := 0; i < midNum; i++ {
		for j := 0; j < botNum/midNum; j++ {
			myNeuralNetwork.layerMid[i].inputs[j] = make(chan float64)
			myNeuralNetwork.layerBot[i*(botNum/midNum)+j].output = myNeuralNetwork.layerMid[i].inputs[j]
		}
	}

	for i := 0; i < botNum; i++ {
		myNeuralNetwork.layerBot[i].inputs = append(myNeuralNetwork.layerBot[i].inputs, myNeuralNetwork.bottomInputs[i])
	}

	for i := 0; i < midNum; i++ {
		myNeuralNetwork.layerMid[i].output = make(chan float64)
		myNeuralNetwork.layerOut.inputs = append(myNeuralNetwork.layerOut.inputs, myNeuralNetwork.layerMid[i].output)
	}

	myNeuralNetwork.output = make(chan float64)

	myNeuralNetwork.layerOut.output = myNeuralNetwork.output
	return &myNeuralNetwork
}

func (myNeuralNetwork *NeuralNetwork) InitRandom() {
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		myNeuralNetwork.layerBot[i].Randomize()
	}

	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		myNeuralNetwork.layerMid[i].Randomize()
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
	go myNeuralNetwork.layerOut.Process()
	for i, item := range myNeuralNetwork.bottomInputs {
		item <- inputs[i]
	}

	return <-myNeuralNetwork.output
}
