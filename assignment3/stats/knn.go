package stats

import (
	"math"
	"sort"
)

// KNNClassifier performs k-nearest neighbors classification.
type KNNClassifier struct {
	K           int
	TrainX      [][]float64 // training features (normalized)
	TrainY      []float64   // training labels
	// Metric controls the distance function: "euclidean" (default) or "cosine"
	Metric      string
}

// euclideanDistance computes distance between two feature vectors.
func euclideanDistance(a, b []float64) float64 {
	var sum float64
	for i := 0; i < len(a); i++ {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

// cosineDistance returns 1 - cosine similarity so that lower is "closer".
func cosineDistance(a, b []float64) float64 {
	var dot, na, nb float64
	for i := 0; i < len(a); i++ {
		ai := a[i]
		bi := b[i]
		dot += ai * bi
		na += ai * ai
		nb += bi * bi
	}
	denom := math.Sqrt(na) * math.Sqrt(nb)
	if denom == 0 {
		if na == 0 && nb == 0 {
			return 0
		}
		return 1
	}
	sim := dot / denom
	return 1 - sim
}

func (knn *KNNClassifier) distance(a, b []float64) float64 {
	switch knn.Metric {
	case "", "euclidean":
		return euclideanDistance(a, b)
	case "cosine":
		return cosineDistance(a, b)
	default:
		return euclideanDistance(a, b)
	}
}

// Predict returns predicted label for a single test instance.
func (knn *KNNClassifier) Predict(x []float64) float64 {
	if len(knn.TrainX) == 0 {
		return 0
	}
	
	// Compute distances to all training points
	type distIdx struct {
		dist float64
		idx  int
	}
	distances := make([]distIdx, len(knn.TrainX))
	for i, trainX := range knn.TrainX {
		distances[i] = distIdx{
			dist: knn.distance(x, trainX),
			idx:  i,
		}
	}
	
	// Sort by distance
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist < distances[j].dist
	})
	
	// Take k nearest neighbors
	k := knn.K
	if k > len(distances) {
		k = len(distances)
	}
	
	// Majority vote
	voteSum := 0.0
	for i := 0; i < k; i++ {
		voteSum += knn.TrainY[distances[i].idx]
	}
	
	// Return sign of average (for binary -1/1 classification)
	if voteSum >= 0 {
		return 1
	}
	return -1
}

// ErrorRate computes classification error rate on a test set.
// Returns fraction of incorrect predictions.
func (knn *KNNClassifier) ErrorRate(testX [][]float64, testY []float64) float64 {
	if len(testX) == 0 {
		return 0
	}
	errors := 0
	for i := 0; i < len(testX); i++ {
		pred := knn.Predict(testX[i])
		if pred != testY[i] {
			errors++
		}
	}
	return float64(errors) / float64(len(testX))
}
