package SimpleActivityRecognitionNeuralNetwork

type NeuronTrained struct {
	Weights []float64
	Bias    float64
}

type NeuralNetworkTrainedModel struct {
	LayerOut   NeuronTrained   `json:"layerOut"`
	LayerOuter []NeuronTrained `json:"layerOuter"`
	LayerMid   []NeuronTrained `json:"layerMid"`
	LayerBot   []NeuronTrained `json:"layerBot"`
}
