from pulp import *

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

def schedule_solve(podList, nodeList):
    podNum = len(podList)
    nodeNum = len(nodeList)

    nRow = [i for i in range(nodeNum)]
    pCol = [i for i in range(podNum)]

    # matrix stands for pod selection
    choices = LpVariable.matrix("Choice", (nRow, pCol),0,1,LpInteger)
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
    print(prob)
    for v in prob.variables():
        print(v.name, "=", v.varValue)  

schedule_solve(podList, nodeList)


