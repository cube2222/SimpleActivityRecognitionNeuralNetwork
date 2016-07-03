package SimpleActivityRecognitionNeuralNetwork

type NeuronTrained struct {
	Weights []float64
	Bias    float64
}

type NeuralNetworkTrainedModel struct {
	LayerOut Neuron   `json:"layerOut"`
	LayerMid []Neuron `json:"layerMid"`
	LayerBot []Neuron `json:"layerBot"`
}
