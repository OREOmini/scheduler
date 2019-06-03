import yaml
import json
import numpy as np


def podDict(l, i):
    d = dict()
    d['name'] = "pod"+str(i)
    d['cpu'] = l[0]
    d['mem'] = l[1]
    return d

def nodeDict(l, i):
    d = dict()
    d['name'] = "node"+str(i)
    d['cpu'] = l[0]
    d['mem'] = l[1]
    d['pnum'] = l[2]
    return d


def get_random_pod_matrix(pod_num):
    return np.random.randint(2, 10, (pod_num, 2))
    
def get_random_node_matrix(node_num):
    nodeList = np.random.randint(2000, 4000, (node_num, 2))
    # t = np.random.randint(80, 110, (node_num, 1))
    t = np.full((node_num, 1), 110)
    return np.c_[nodeList, t]

res = []
node_num = 3
node_max = 20
c = 0
i = 0

while node_num <= node_max:
    print(i)
    d = dict()
    c = c+1
    if c == 300:
        c = 0
        node_num = node_num+1
    p = get_random_pod_matrix(20)
    n = get_random_node_matrix(node_num)
    d['Pods'] = p.tolist()
    d['Nodes'] = n.tolist()
    res.append(d)
    i = i+1

print(len(res))

# # pprint.pprint(result_dict)
with open("temp1s.json", 'w') as f:
    json.dump(res, f)
    