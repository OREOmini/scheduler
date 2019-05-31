package main

import (
	"fmt"
	"github.com/sbinet/go-python"
	"strconv"
	"strings"
)

type PodToSolve struct {
	cpu float64
	memory float64
}
type NodeToSolve struct {
	cpu float64
	memory float64
	pod_space float64
}

// for simplicity mi = m, ki = k
func toMb(s string) float64 {
	s = strings.ToLower(s)
	if strings.HasSuffix(s, "m") || strings.HasSuffix(s, "mi") || strings.HasSuffix(s, "mb") {
		num := strings.Split(s, "m")[0]
		if s, err := strconv.ParseFloat(num, 1); err == nil {
			return s
		}
	} else if strings.HasSuffix(s, "k") || strings.HasSuffix(s, "ki") || strings.HasSuffix(s, "kb") {
		num := strings.Split(s, "k")[0]
		if s, err := strconv.ParseFloat(num, 1); err == nil {
			return s / 1000
		}
	} else if strings.HasSuffix(s, "g") || strings.HasSuffix(s, "gi") || strings.HasSuffix(s, "gb") {
		num := strings.Split(s, "g")[0]
		if s, err := strconv.ParseFloat(num, 1); err == nil {
			return s * 1000
		}
	} else {
		if s, err := strconv.ParseFloat(s, 1); err == nil {
			return s * 1000
		}
	}
	return 0
}

func getCpuFromString(cpuStr string) float64{
	cpuStr = strings.TrimSpace(cpuStr)
	if cpuStr == "" {
		return 0
	} else {
		return toMb(cpuStr)
	}
}
func getMemoryFromString(memoryStr string) float64{
	memoryStr = strings.TrimSpace(memoryStr)
	memoryStr = strings.ToLower(memoryStr)
	if memoryStr == "" {
		return 0
	} else {
		return toMb(memoryStr)
	}
}

// [cpu, memory. pod space]
func (node *Node) getSolveParamList() []float64 {
	res := make([]float64, 3)
	for k, v := range node.Status.Allocatable {
		if k == "cpu" {
			res[0] = getCpuFromString(v)
		} else if k == "memory" {
			res[1] = getMemoryFromString(v)
		} else if k == "pods" {
			i, err := strconv.Atoi(v)
			if err == nil {
				res[2] = float64(i)
			}
		}
	}
	return res
}

// [cpu, memory]
func (pod *Pod) getSolveParamList() []float64 {
	totalCpu := 0.0
	totalMemory := 0.0
	for _, con := range pod.Spec.Containers {
		for k, v := range con.Resources.Requests {
			if k == "cpu" {
				totalCpu += getCpuFromString(v)
			} else if k == "memory" {
				totalMemory += getMemoryFromString(v)
			}
		}
	}
	return []float64{totalCpu, totalMemory}
}

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

func getPyList(list []float64) *python.PyObject {
	size := len(list)
	l := python.PyList_New(size)
	for i := 0; i < size; i++ {
		tempNum := python.PyInt_FromLong(int(list[i]))
		python.PyList_SET_ITEM(l, i, tempNum)
	}
	return l
}

func toGoList(pyList *python.PyObject) []int {
	size := python.PyList_Size(pyList)
	res := make([]int, size)
	for i := 0; i < size; i++ {
		temp := python.PyList_GET_ITEM(pyList, i)
		//println(python.PyInt_AS_LONG(temp))
		num := python.PyInt_AS_LONG(temp)
		if num == 0 {
			res[i] = 0
		} else {
			res[i] = 1
		}
	}
	return res
}

func callPySolve(n [][]float64, p [][]float64, ok chan bool) *python.PyObject{
	m := python.PyImport_ImportModule("sys")
	if m == nil {
		fmt.Println("import error")
		return nil
	}
	path := m.GetAttrString("path")
	if path == nil {
		fmt.Println("get path error")
		return nil
	}
	//add to current path
	currentDir := python.PyString_FromString("")
	python.PyList_Insert(path, 0, currentDir)


	solver := python.PyImport_ImportModule("lp-solver")
	if solver == nil {
		fmt.Println("import error")
		return nil
	}
	a := solver.GetAttrString("a")
	fmt.Printf("[VARS] a = %#v\n", python.PyInt_AsLong(a))

	podList := python.PyList_New(len(p))
	for i := 0; i < len(p); i++ {
		python.PyList_SET_ITEM(podList, i, getPyList(p[i]))
	}

	nodeList := python.PyList_New(len(n))
	for i := 0; i < len(n); i++ {
		python.PyList_SET_ITEM(nodeList, i, getPyList(n[i]))
	}

	args := python.PyTuple_New(2)
	python.PyTuple_SET_ITEM(args, 0, podList)
	python.PyTuple_SET_ITEM(args, 1, nodeList)

	solverFunc := solver.GetAttrString("schedule_solve")
	res := solverFunc.Call(args, python.Py_None)
	if res == nil {
		fmt.Println("call error")
		return nil
	}
	close(ok)
	return res
}

func solve(nodeList []Node, podList []*Pod) [][]int{
	podMatrix := make([][]float64, len(podList))
	for i := 0;i < len(podList); i++ {
		podMatrix[i] = podList[i].getSolveParamList()
	}
	fmt.Println("pod param list:", podMatrix)

	nodeMatrix := make([][]float64, len(nodeList))
	for i := 0;i < len(nodeList); i++ {
		nodeMatrix[i] = nodeList[i].getSolveParamList()
	}
	fmt.Println("node param list:", nodeMatrix)


	//res := make(*python.PyObject)
	ok := make(chan bool)

	pyMatrix := callPySolve(nodeMatrix, podMatrix, ok)
	<- ok
	//pyMatrix := <- res
	//pyMatrix = python.PyList_GetItem(pyMatrix, 0)

	size := python.PyList_Size(pyMatrix)
	podAllocation := make([][]int, size)
	for i := 0; i < size; i++ {
		pyListTemp := python.PyList_GET_ITEM(pyMatrix, i)
		temp := toGoList(pyListTemp)
		podAllocation[i] = temp
	}
	return podAllocation
}

func schedulePodUsingSolver(podList []*Pod) error{
	fmt.Println("get pod list len:", len(podList))
	nodeList, err := getNodes()
	if err != nil {
		return err
	}

	podAllocation := solve(nodeList.Items, podList)
	fmt.Println(podAllocation)

	//ok := make(chan bool)
	err = schedulePodOnResult(podAllocation, podList, nodeList.Items)

	if err != nil {
		return err
	}

	//<- ok
	//return fmt.Errorf("Unable to schedule pod (%s) failed to fit in any node", pod.Metadata.Name)
	return nil
}

func schedulePodOnResult(podAllocation [][]int,podList []*Pod, nodeList []Node) error {
	if len(podAllocation) != len(nodeList) || len(podAllocation[0]) != len(podList) {
		return fmt.Errorf("Matrix Incompability.")
	}
	for i := 0; i < len(podAllocation); i++ {
		for j := 0; j < len(podAllocation); j++ {
			if podAllocation[i][j] == 1 {
				n := nodeList[i]
				p := podList[j]
				fmt.Println("bind", p.Metadata.Name, "to", n.Metadata.Name)
				//mu := sync.sMutex{}
				//mu.Lock()
				e := bind(p, n)
				if e != nil {
					return e
				}
				//mu.Unlock()
			}
		}
	}
	//close(ok)
	return nil
}


//func main() {
//	res := make(chan *python.PyObject)
//	ok := make(chan bool)
//
//	go callPySolve(res, ok)
//
//	//<- ok
//
//	pyMatrix := <- res
//	//pyMatrix = python.PyList_GetItem(pyMatrix, 0)
//
//	size := python.PyList_Size(pyMatrix)
//	println("---",size)
//	podAllocation := make([][]int, size)
//	for i := 0; i < size; i++ {
//		pyListTemp := python.PyList_GET_ITEM(pyMatrix, i)
//		temp := toGoList(pyListTemp)
//		fmt.Println(temp)
//		podAllocation[i] = temp
//	}
//}
//

