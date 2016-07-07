package SimpleActivityRecognitionNeuralNetwork

import (
	"math/rand"
)

type Neuron struct {
	inputs          []chan float64
	output          chan float64
	outputReceivers int
	weights         []float64
	bias            float64
}

func (n *Neuron) heaviside(f float64) int32 {
	if f < 0 {
		return 0
	}
	return 1
}

func (myNeuron *Neuron) Randomize() {
	var i int
	w := make([]float64, len(myNeuron.inputs))
	for i = 0; i < len(w); i++ {
		w[i] = rand.Float64()*2 - 1
	}
	myNeuron.weights = w
	myNeuron.bias = rand.Float64()*2 - 1
}

func (n *Neuron) Process() {
	sum := n.bias
	for i, input := range n.inputs {
		sum += n.weights[i] * (<-input)
	}
	answer := float64(n.heaviside(sum))
	for i := 0; i < n.outputReceivers; i++ {
		n.output <- answer
	}
}

func (n *Neuron) Adjust(inputs []float64, delta int32, learningRate float64) {
	for i, input := range inputs {
		n.weights[i] += input * float64(delta) * learningRate
	}
	n.bias += float64(delta) * learningRate
}
