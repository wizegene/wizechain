package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/mb-14/gomarkov"
	"io/ioutil"
	"strings"
	"sync"
)

type DnaTrainingSet []DnaTrainingData

type DnaTrainingData struct {
	blockText   string
	contentText string
}

type DNATrainer interface {
	GetTrainingSet() map[int]string
	calculateDifficultyOfTarget(target string) (score float64)
	Train() string
	buildMarkovModel() (*gomarkov.Chain, error)
	saveMarkovModel(chain *gomarkov.Chain)
	loadModel() (*gomarkov.Chain, error)
	generateProof(chain *gomarkov.Chain) string
}

type Trainer struct {
	DNATrainer
}

func NewTrainer() *Trainer {
	return &Trainer{}
}

func (t *Trainer) GetTrainingSet() map[int]string {
	// todo change quotes.json to real data source (content and blocks)
	data, err := ioutil.ReadFile("./dict/quotes.json")
	if err != nil {
		panic(err)
	}

	mapper := make(map[int]string, len(data))
	splitted := strings.Split(string(data), "\n")

	for i := 0; i < len(splitted); i++ {

		mapper[i] = string(splitted[i])

	}
	return mapper

}

func (t *Trainer) calculateDifficultyOfTarget(target string) (score float64) {
	return 0
}

func (t *Trainer) Train() string {

	chain, err := t.buildMarkovModel()
	if err != nil {
		panic(err)
	}
	t.saveMarkovModel(chain)
	chain, err = t.loadModel()
	return t.generateProof(chain)

}

func (t *Trainer) buildMarkovModel() (chain *gomarkov.Chain, err error) {

	trainingData := t.GetTrainingSet()
	chain = gomarkov.NewChain(1)
	var wg sync.WaitGroup
	wg.Add(len(trainingData))

	for _, td := range trainingData {
		td := td
		go func() {
			defer wg.Done()
			chain.Add(strings.Split(td, " "))
		}()
	}
	wg.Wait()

	return chain, nil

}

func (t *Trainer) saveMarkovModel(chain *gomarkov.Chain) {
	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (t *Trainer) loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func (t *Trainer) generateProof(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	return fmt.Sprintf(strings.Join(tokens[1:len(tokens)-1], " "))
}
