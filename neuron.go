package SimpleActivityRecognitionNeuralNetwork

import "math/rand"

type Neuron struct {
	weights []float32
	bias    float32
}

func (p *Neuron) heaviside(f float32) int32 {
	if f < 0 {
		return 0
	}
	return 1
}

func NewNeuron(n int32) *Neuron {
	var i int32
	w := make([]float32, n, n)
	for i = 0; i < n; i++ {
		w[i] = rand.Float32()*2 - 1
	}
	return &Neuron{
		weights: w,
		bias:    rand.Float32()*2 - 1,
	}
}

func (p *Neuron) Process(inputs []int32) int32 {
	sum := p.bias
	for i, input := range inputs {
		sum += float32(input) * p.weights[i]
	}
	return p.heaviside(sum)
}

func (p *Neuron) Adjust(inputs []int32, delta int32, learningRate float32) {
	for i, input := range inputs {
		p.weights[i] += float32(input) * float32(delta) * learningRate
	}
	p.bias += float32(delta) * learningRate
}