import sys
import math
import yaml
import operator
import copy
import time
from scipy.special import perm

# Input is a dict, P is a list, N is also a list
def main(argv):
	file_name = argv[1]
	fp = open(file_name, 'r')
	inputs = yaml.load(fp)
	# inputs = yaml.load(fp, Loader=yaml.Loader)
	fp.close()

	starttime = time.time()
	# print(inputs['Pods'])
	permutation(inputs['Pods'])
	
	# print ('permutation of pods:', perm_pod)

	all_ALs = processing_permutation_of_pods(perm_pod, inputs['Nodes'])
	# print ('all_ALs:', all_ALs)

	AL = find_AL_with_smallest_na(all_ALs)
	print ('AL:', AL)

	endtime = time.time()

	normal = endtime - starttime
	print ('seconds:', normal)
	return AL

perm_pod = []
# the permutation of pods in P, pods is a list.
def permutation(pods):
	n = []
	permu(pods, n)


def permu(list,stack):
	# print('permu')
	# print(len(list))
	# print(len(stack))

	if not list:
		stack = copy.deepcopy(stack)
		perm_pod.append(stack)
	else:
		for i in range(len(list)):
			stack.append(list[i])
			del list[i]
			permu(list,stack)
			list.insert(i,stack.pop())




# perm_p is a list of permutation of pods, nodes is a list.
def processing_permutation_of_pods(perm_p, nodes):
	AL_dict = {}
	for queue in perm_p:
		AL = []
		for pod in queue:
			enodes = filter_eligible_nodes(pod['cpu'], pod['mem'], nodes) # eligible nodes
			if len(enodes) == 1:
				AL.append([pod['name'], enodes[0]['name']])
			else:
				for n in enodes:
					count = 0
					n_compare = nodes_for_compare(n, enodes)
					for nc in n_compare:
						if score(nc, n, pod) > 0:
							count += 1  # count the number of nodes with higher score than n
					# n_priority[count] = n means there are 'count' nodes with higher score than n, store n at key/index 'count'
					if count == 0:
						en = n
						break
				# AL.append([pod['name'], n_priority]) store the node list for pod['name']
				AL.append([pod['name'], en['name']])  # <pod, node>
		AL_dict[perm_p.index(queue)] = AL # add the AL list to the AL_dict, corresponding to the queue of the permutation of pods
	return AL_dict


def filter_eligible_nodes(pcpu, pmem, nodes):
	enodes = []
	for node in nodes:
		if (int(node['cpu']) > int(pcpu) & (int(node['mem']) > int(pmem))):
			enodes.append(node)
	return enodes


def score(nj, nk, p):
	x = float(nj['cpu'])/nj['mem']
	y = float(nk['cpu'])/nk['mem']
	z = float(p['cpu'])/p['mem']
	D0 = 110.0

	a = math.ceil(1-nj['pnum']/D0) - math.ceil(1-nk['pnum']/D0)
	s = (1 + a) * pow((y/z - 1), 2) - (1 + a) * pow((x/z - 1), 2)
	return s


def nodes_for_compare(n, nodes):
	n_new = []
	for node in nodes:
		if not operator.eq(n, node):
			n_new.append(node)
	return n_new




# ALs is a dict of all AL corresponding to each queue in perm_p.
def find_AL_with_smallest_na(ALs):
	n = []
	for i in range(len(ALs)):
		n.append(count_the_number_of_nodes(ALs[i]))
	index = n.index(min(n))
	return ALs[index]

# AL is a list.
def count_the_number_of_nodes(AL):
	active_node = []
	for pair in AL:
		if pair[1] in active_node:
			continue
		else:
			active_node.append(pair[1])
	return len(active_node)




# file_name = 'result_data/result_3.yaml'
# al = solve(file_name)
# print(al)
if __name__ == '__main__':
	main(sys.argv)