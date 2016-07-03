package SimpleActivityRecognitionNeuralNetwork

import "encoding/json"

type NeuralNetwork struct {
	layerOut     Neuron
	layerMid     []Neuron
	layerBot     []Neuron
	bottomInputs []chan float64
	output chan float64
}

func New(midNum int64, botNum int64) *NeuralNetwork {
	myNeuralNetwork := NeuralNetwork{}
	myNeuralNetwork.layerBot = make([]Neuron, botNum)
	myNeuralNetwork.layerMid = make([]Neuron, midNum)
	myNeuralNetwork.bottomInputs = make([]chan float64, botNum)

	for i := 0; i < midNum; i++ {
		myNeuralNetwork.layerMid[i].inputs = make([]chan float64, botNum/midNum)
	}
	for i := 0; i < midNum; i++ {
		for j := 0; j < botNum/midNum; i++ {
			myNeuralNetwork.layerBot[i*(botNum/midNum)+j].output = myNeuralNetwork.layerMid[i].inputs[j]
		}
	}

	for i := 0; i < botNum; i++ {
		myNeuralNetwork.layerBot[i].inputs = append(myNeuralNetwork.layerBot[i].inputs, myNeuralNetwork.bottomInputs[i])
	}

	for i := 0; i < midNum; i++ {
		myNeuralNetwork.layerOut.inputs = append(myNeuralNetwork.layerOut.inputs, myNeuralNetwork.layerMid[i].output)
	}

	myNeuralNetwork.output = make(chan float64)

	myNeuralNetwork.layerOut.output = myNeuralNetwork.output
}

func (myNeuralNetwork *NeuralNetwork) initRandom() {
	for i := 0; i < len(myNeuralNetwork.layerBot); i++ {
		myNeuralNetwork.layerBot[i].Randomize()
	}

	for i := 0; i < len(myNeuralNetwork.layerMid); i++ {
		myNeuralNetwork.layerMid[i].Randomize()
	}

	myNeuralNetwork.layerOut.Randomize()
}

func (myNeuralNetwork *NeuralNetwork) initFromJson(data []byte) {
	myNeuralNetworkModel := NeuralNetworkTrainedModel{}
	err := json.Unmarshal(data, &myNeuralNetworkModel)
	if err != nil {
		return err
	}

}

func (myNeuralNetwork *NeuralNetwork) saveToJson() ([]byte, error) {
	myNeuralNetworkModel := NeuralNetworkTrainedModel{}

	myNeuralNetworkModel.LayerBot = make([]NeuronTrained, len(myNeuralNetwork.layerBot))
	myNeuralNetworkModel.LayerMid = make([]NeuronTrained, len(myNeuralNetwork.layerMid))

	for i := 0; i < len(myNeuralNetworkModel.LayerBot); i++ {
		myNeuralNetworkModel.LayerBot[i].weights = myNeuralNetwork.layerBot[i].weights
		myNeuralNetworkModel.LayerBot[i].bias = myNeuralNetwork.layerBot[i].bias
	}

	for i := 0; i < len(myNeuralNetworkModel.LayerMid); i++ {
		myNeuralNetworkModel.LayerMid[i].weights = myNeuralNetwork.layerMid[i].weights
		myNeuralNetworkModel.LayerMid[i].bias = myNeuralNetwork.layerMid[i].bias
	}

	myNeuralNetworkModel.LayerOut.weights = myNeuralNetwork.layerOut.weights
	myNeuralNetworkModel.LayerOut.bias = myNeuralNetwork.layerOut.bias

	data, err := json.Marshal(&myNeuralNetworkModel)
	if err != nil {
		return err
	}
	return data
}