package main

import "fmt"
import "github.com/draffensperger/golp"

type Pod struct {
	cpu float64
	memory float64
}
type Node struct {
	cpu float64
	memory float64
	pod_space float64
}


func T(matrix [][]float64) [][]float64 {
	row := len(matrix)
	col := len(matrix[0])
	result := make([][]float64, col)

	for i := 0; i < col; i++ {
		temp := make([]float64, row)
		for j := 0; j < row; j++ {
			temp[j] = matrix[j][i]
		}
		result[i] = temp
	}
	return result
}

func getData(podList []Pod, nodeList []Node) (a []float64, b []float64, A []float64, B []float64, C []float64, E1 []float64, E2 []float64){
	podNum := len(podList)
	nodeNum := len(nodeList)

	E1 = make([]float64, podNum)
	E2 = make([]float64, podNum)
	a = make([]float64, podNum)
	b = make([]float64, podNum)
	A = make([]float64, nodeNum)
	B = make([]float64, nodeNum)
	C = make([]float64, nodeNum)


	for i := 0; i < podNum; i++ {
		E1[i] = 1
	}
	for i := 0; i < nodeNum; i++ {
		E2[i] = 1
	}

	for i := 0; i < podNum; i++ {
		a[i] = podList[i].cpu
	}
	for i := 0; i < podNum; i++ {
		b[i] = podList[i].memory
	}
	for i := 0; i < nodeNum; i++ {
		A[i] = nodeList[i].cpu
	}
	for i := 0; i < nodeNum; i++ {
		B[i] = nodeList[i].memory
	}
	for i := 0; i < nodeNum; i++ {
		C[i] = nodeList[i].pod_space
	}
	return
}


func main() {
	pod1 := Pod{20, 2}
	pod2 := Pod{20, 3}
	pod3 := Pod{20, 3}

	node1 := Node{70, 20, 10}
	node2 := Node{90, 10, 10}

	podList := []Pod{pod1, pod2, pod3}
	nodeList := []Node{node1, node2}

	podNum := len(podList)
	nodeNum := len(nodeList)

	//a, b, A, B, C, E1, E2 := getData(podList, nodeList)
	a, b, A, B, _, _, _ := getData(podList, nodeList)


	myLp := golp.NewLP(nodeNum, podNum)
	for i := 0; i < nodeNum; i++ {
		myLp.AddConstraint(a, golp.LE, A[i])
		myLp.AddConstraint(b, golp.LE, B[i])
	}


	lp := golp.NewLP(0, 2)
	lp.AddConstraint([]float64{110.0, 30.0}, golp.LE, 4000.0)
	lp.AddConstraint([]float64{1.0, 1.0}, golp.LE, 75.0)
	lp.SetObjFn([]float64{143.0, 60.0})
	lp.SetMaximize()

	//lp.Solve()
	//vars := lp.Variables()
	//fmt.Printf("Plant %.3f acres of barley\n", vars[0])
	//fmt.Printf("And  %.3f acres of wheat\n", vars[1])
	//fmt.Printf("For optimal profit of $%.2f\n", lp.Objective())


	// No need to explicitly free underlying C structure as golp.LP finalizer will
}