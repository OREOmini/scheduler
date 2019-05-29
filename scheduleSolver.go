package main

import (
	"fmt"
	"github.com/sbinet/go-python"
)

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


func main() {
	res := make(chan *python.PyObject)
	ok := make(chan bool)

	go callPySolve(res, ok)

	//<- ok

	pyMatrix := <- res
	//pyMatrix = python.PyList_GetItem(pyMatrix, 0)

	size := python.PyList_Size(pyMatrix)
	println("---",size)
	podAllocation := make([][]int, size)
	for i := 0; i < size; i++ {
		pyListTemp := python.PyList_GET_ITEM(pyMatrix, i)
		temp := toGoList(pyListTemp)
		fmt.Println(temp)
		podAllocation[i] = temp
	}
}


