from pulp import *
import time as time
import numpy as np

a = 10

NODE_CPU_INDEX = 0
NODE_MEMORY_INDEX = 1
NODE_POD_SPACE_INDEX = 2

POD_CPU_INDEX = 0
POD_MEMORY_INDEX = 1

podList = [
    [10, 3],
    [10, 1],
    [10, 3]
]
nodeList = [
    [30, 5, 9],
    [40, 3, 7]
]

def schedule_solve(podList, nodeList, VERBOSE = False):
    podNum = len(podList)
    nodeNum = len(nodeList)

    print(podList)
    print(nodeList)

    nRow = [i for i in range(nodeNum)]
    pCol = [i for i in range(podNum)]

    # matrix stands for pod selection
    choices = LpVariable.matrix("choice", (nRow, pCol),0,1,LpInteger)
    # node usage
    nodeOccupation = LpVariable.matrix("node", nRow,0,1,LpInteger)

    prob = LpProblem("lp", LpMinimize)
    # minimaze the node usage
    prob += lpSum(nodeOccupation), "objective function"

    # one pod can only assign to one node
    for c in range(len(pCol)):
        prob += lpSum([choices[r][c] for r in range(len(nRow))]) == 1, ""


    for r in range(len(nRow)):
        # if there is pod in this node
        for c in range(len(pCol)):
            prob += choices[r][c] <= nodeOccupation[r], ""
        # satisfy node cpu capacity 
        prob += lpSum([podList[c][POD_CPU_INDEX] * choices[r][c]  for c in range(len(pCol))]) <= nodeList[r][NODE_CPU_INDEX], ""
        # satisfy node memory capacity 
        prob += lpSum([podList[c][POD_MEMORY_INDEX] * choices[r][c]  for c in range(len(pCol))]) <= nodeList[r][NODE_MEMORY_INDEX], ""
        # satisfy node pod number capacity 
        prob += lpSum([choices[r][c]  for c in range(len(pCol))]) <= nodeList[r][NODE_POD_SPACE_INDEX]
        
    prob.solve()
    print(LpStatus[prob.status])

    print("objective:",value(prob.objective))
    if VERBOSE: 
        print(prob)
    
    result = [[0 for col in range(podNum)] for row in range(nodeNum)]
    for v in prob.variables():
        t = v.name.split('_')
        if t[0] == 'choice':
            result[int(t[1])][int(t[2])] = v.varValue
        # print(v.name, "=", v.varValue)  
    
    print(result)
    return result


# start = time.time()
# podList = np.random.rand(2, 2)
# nodeList = np.random.rand(1, 3)
# schedule_solve(podList, nodeList)
# print(time.time() - start)


