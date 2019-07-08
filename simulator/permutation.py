import yaml
import simulatorInput
import time


def main():

	for index in range(0,2): # change the second number to the number of input file + 1
		file_name = 'input_data/pod_space_equal_to_110/input_'+str(index)+'.yaml'  # the file name of input file, such as input1.yaml, input2.yaml
		fp = open(file_name, 'r')
		inputs = yaml.load(fp, Loader=yaml.FullLoader)
		starttime = time.time()
		# print(inputs)

		# for p, pod in zip(pods, inputs['Pods']):
		# 	p['cpu'] = pod['cpu']
		# 	p['mem'] = pod['mem']
		# 
		# print(inputs['Pods'])
		permutation(inputs['Pods'])
		# print ('final:', perm_pod)
		AL = simulatorInput.simulator(perm_pod, inputs['Nodes'])
		print(AL)
		fp.close()

		fp = open(file_name, 'a')
		endtime = time.time()
		normal = endtime - starttime
		result = {'AL': AL, 'Time': normal}  # allocation result will be stored in the input file as a dict, the key is AL, the value is the allocation list.
		yaml.dump(result, fp)
		fp.close()


perm_pod = []
def permutation(pods):
	n = []
	permu(pods, n)


def permu(list,stack):
	if not list:
		# print (stack)
		tmp_stack = []
		for i in range(len(stack)):
			tmp_stack.append(stack[i])
		perm_pod.append(tmp_stack)
		# print (perm_pod)
	else:
		for i in range(len(list)):
			stack.append(list[i])
			del list[i]
			permu(list,stack)
			list.insert(i,stack.pop())

if __name__ == '__main__':
	main()