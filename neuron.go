package SimpleActivityRecognitionNeuralNetwork

import "math/rand"

type Neuron struct {
	inputs  []chan float64
	output  chan float64
	weights []float64
	bias    float64
}

func (p *Neuron) heaviside(f float64) int32 {
	if f < 0 {
		return 0
	}
	return 1
}

func (myNeuron *Neuron) Randomize() {
	var i int32
	w := make([]float64, len(myNeuron.inputs))
	for i = 0; i < len(w); i++ {
		w[i] = rand.Float64()*2 - 1
	}
	return &Neuron{
		weights: w,
		bias:    rand.Float64()*2 - 1,
	}
}

func (p *Neuron) Process(inputs []float64) int32 {
	sum := p.bias
	for i, input := range inputs {
		sum += input * p.weights[i]
	}
	return p.heaviside(sum)
}

func (p *Neuron) Adjust(inputs []float64, delta int32, learningRate float64) {
	for i, input := range inputs {
		p.weights[i] += input * float64(delta) * learningRate
	}
	p.bias += float64(delta) * learningRate
}
