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
			return s
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
		fmt.Printf("%s: %s\n", k, v)
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

func callPySolve(res chan *python.PyObject, ok chan bool) {
	m := python.PyImport_ImportModule("sys")
	if m == nil {
		fmt.Println("import error")
		return
	}
	path := m.GetAttrString("path")
	if path == nil {
		fmt.Println("get path error")
		return
	}
	//加入当前目录，空串表示当前目录
	currentDir := python.PyString_FromString("")
	python.PyList_Insert(path, 0, currentDir)


	solver := python.PyImport_ImportModule("lp-solver")
	if solver == nil {
		fmt.Println("import error")
		return
	}
	a := solver.GetAttrString("a")
	fmt.Printf("[VARS] a = %#v\n", python.PyInt_AsLong(a))

	podList := python.PyList_New(3)
	python.PyList_SET_ITEM(podList, 0, getPyList([]float64{10, 3}))
	python.PyList_SET_ITEM(podList, 1, getPyList([]float64{10, 1}))
	python.PyList_SET_ITEM(podList, 2,  getPyList([]float64{10, 3}))
	//
	nodeList := python.PyList_New(3)
	python.PyList_SET_ITEM(nodeList, 0, getPyList([]float64{20, 5, 10}))
	python.PyList_SET_ITEM(nodeList,1, getPyList([]float64{10, 4, 10}))
	python.PyList_SET_ITEM(nodeList,2, getPyList([]float64{10, 4, 10}))

	args := python.PyTuple_New(2)
	python.PyTuple_SET_ITEM(args, 0, podList)
	python.PyTuple_SET_ITEM(args, 1, nodeList)

	fmt.Println(python.PyList_GET_ITEM(podList, 0))
	//
	solverFunc := solver.GetAttrString("schedule_solve")
	fmt.Printf("[FUNC] = %#v\n", solverFunc)

	done := make(chan bool)
	//res := make(chan *python.PyObject)
	go func() {
		res <- solverFunc.Call(args, python.Py_None)
		close(done)

	}()
	//res := solverFunc.Call(args, python.Py_None)
	<- done
	if res == nil {
		fmt.Println("call error")
		return
	}
	close(ok)
	//return res
}

//
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
//
