package main

import (
	"fmt"
	"github.com/irr/bayesian"
)

const (
    Good bayesian.Class = "good"
    Bad bayesian.Class = "bad"
)

func main() {
	c := bayesian.NewClassifier(Good, Bad)
    
    c.Learn([]string{"tall", "handsome", "rich"}, Good)
    c.Learn([]string{"bald", "poor", "ugly", "bitch", "none"}, Bad)
    
    doc := []string{"tall", "poor", "rich", "dummy", "nothing"}

    var scores []float64
    
    scores, _, _ = c.LogScores(doc)
    fmt.Println("Log", scores)

    scores, _, _ = c.ProbScores(doc)
    fmt.Println("Prob", scores)

    scores, _, _ = c.SafeProbScores(doc)
    fmt.Println("SafeProb", scores)

    fmt.Println("Learned", c.Learned())
    fmt.Println("Seen", c.Seen())
    fmt.Println("WordCount", c.WordCount())
}

