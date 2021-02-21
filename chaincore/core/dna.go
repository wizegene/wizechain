package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	blake2b "github.com/minio/blake2b-simd"
	"math/rand"
	"time"
)

/*

	DNA crypto package for wizechain

*/

var mutationRate float64 = 0.00001

type DNARaw struct {
	a []byte
	o []byte
	D []byte
	t []byte
}

type Organism struct {
	DNA     []byte
	Fitness float64
	ParentA []byte
	ParentB []byte
	IOrganism
}

type IOrganism interface {
	calculateFitness(target []byte)
	mutate(mutationRate float64)
}

type IDna interface {
	setTarget() string
	createOrganism(target []byte) (organism Organism)
	createPopulation(target []byte, popSize uint) (population []Organism)
	createGenePool(population []Organism, target []byte, maxFitness float64) (pool []Organism)
	crossover(d1 Organism, d2 Organism) Organism
	naturalSelection(pool []Organism, population []Organism, target []byte) []Organism
	getBest(population []Organism) Organism
}

type DNAProcessor struct {
	IDna
}

type Ancestor struct {
	mother     []byte
	father     []byte
	genePool   []Organism
	generation int
}

var A *Ancestor
var O Organism

func NewProcessor() *DNAProcessor {
	return new(DNAProcessor)
}

func (d *DNAProcessor) setTarget() string {
	trainer := NewTrainer()
	target := trainer.Train()
	buf := bytes.NewBuffer(make([]byte, 32))
	i := 0
	if len(target) > 16 && len(target) <= 32 {
		tBytes := []byte(target)
		for i < 16 {
			buf.WriteByte(tBytes[i])

			i++
		}

		target = string(buf.Bytes())

	} else {
		d.setTarget()
	}
	return target
}

func (d *DNAProcessor) createOrganism(target []byte) (organism Organism) {
	ba := make([]byte, len(target))
	for i := 0; i < len(target); i++ {
		ba[i] = byte(rand.Intn(95) + 32)
	}

	organism = Organism{
		DNA:     ba,
		Fitness: 0,
	}

	organism.calculateFitness(target)
	return
}

func (d *DNAProcessor) createPopulation(target []byte, popSize uint) (population []Organism) {

	population = make([]Organism, popSize)
	for i := 0; i < int(popSize); i++ {
		population[i] = d.createOrganism(target)
	}
	return

}

func (d *DNAProcessor) createGenePool(population []Organism, target []byte, maxFitness float64) (pool []Organism) {
	pool = make([]Organism, 0)
	for i := 0; i < len(population); i++ {
		population[i].calculateFitness(target)
		num := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}

	//A.genePool = pool

	return
}

func (d *DNAProcessor) naturalSelection(pool []Organism, population []Organism, target []byte) []Organism {
	next := make([]Organism, len(population))
	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		pa, _ := json.Marshal(a)
		pb, _ := json.Marshal(b)
		next[i].ParentA = pa
		next[i].ParentB = pb

		child := d.crossover(a, b)
		child.mutate(mutationRate)
		child.calculateFitness(target)
		next[i] = child

	}
	return next
}

func (d *DNAProcessor) crossover(d1 Organism, d2 Organism) Organism {
	child := Organism{
		DNA:     make([]byte, len(d1.DNA)),
		Fitness: 0,
	}

	mid := rand.Intn(len(d1.DNA))
	for i := 0; i < len(d1.DNA); i++ {
		if i > mid {
			child.DNA[i] = d1.DNA[i]
		} else {
			child.DNA[i] = d2.DNA[i]
		}
	}
	return child
}

func (d *DNAProcessor) getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}

func (o *Organism) mutate(mutationRate float64) {
	for i := 0; i < len(o.DNA); i++ {
		if rand.Float64() < mutationRate {
			o.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

func (o *Organism) calculateFitness(target []byte) {

	score := 0
	for i := 0; i < len(o.DNA); i++ {
		if o.DNA[i] == target[i] {
			score++
		}
	}
	o.Fitness = float64(score) / float64(len(o.DNA))
	return

}

func GetDNA(popSize uint) (dnaProof *DNARaw) {

	d := NewProcessor()
	target := []byte(d.setTarget())
	// @dev todo will replace with block time
	startTime := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := d.createPopulation(target, popSize)
	found := false
	// @dev will replace generation with block count / height?
	generation := 0

	var elapsed time.Duration

	for !found {
		generation++
		bestOrganism := d.getBest(population)
		elapsed = time.Since(startTime)
		fmt.Printf("\r generation: %d | %s | fitness: %2f | elapsed: %2f", generation, string(bestOrganism.DNA), bestOrganism.Fitness, elapsed)
		dur, _ := time.ParseDuration("3s")

		if bytes.Compare(bestOrganism.DNA, target) == 0 || elapsed.Seconds() <= dur.Seconds() {
			found = true
			oBytes, _ := json.Marshal(bestOrganism)
			pop, _ := json.Marshal(population)
			h := blake2b.New256()
			h.Reset()
			h.Sum(pop)
			h.Sum(oBytes)

			Dh := blake2b.New512()
			Dh.Write(bestOrganism.DNA)

			dnaProof = &DNARaw{
				a: h.Sum(nil),
				o: oBytes,
				D: Dh.Sum(nil),
				t: target,
			}
		} else {
			maxFitness := bestOrganism.Fitness
			pool := d.createGenePool(population, target, maxFitness)
			population = d.naturalSelection(pool, population, target)
		}
	}

	//elapsed = time.Since(startTime)
	fmt.Printf("\nTime taken: %s\n", elapsed)

	spew.Dump(dnaProof)

	return

}
