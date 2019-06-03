from run_simulator import extract_output
import yaml
import subprocess
# from simulator.run_simulator import
import os

def getFilePath(i, f = True):
    if f:
        return "input_data/pod_space_less_than_110/input_"+str(i)+".yaml"
    else:
        return "result_data/pod_space_equal_to_110/input_"+str(i)+".yaml"

DIR = 'input_data/pod_space_less_than_110' #要统计的文件夹
file_num = len([name for name in os.listdir(DIR) if os.path.isfile(os.path.join(DIR, name))])
print (file_num )

results = []

for i in range(2):
    print(i)
    d = dict()
    file_name = getFilePath(i)

    fp = open(file_name, 'r')
    inputs = yaml.load(fp)
    fp.close()

    pods = inputs['Pods']
    nodes = inputs['Nodes']
    # if len(pods) > 13:
    #     continue

    d['pods'] = pods
    print('pod len:', len(pods))
    d['nodes'] = nodes
    print('node len:', len(nodes))

    process = subprocess.run(
        ['python3', 'simulator.py', file_name],
        stdin=subprocess.PIPE, stdout=subprocess.PIPE
    )
    p, t = extract_output(process.stdout)
    print(p)

    d['result'] = p
    d['time'] = t
    d['index'] = i
    results.append(d)

import pprint
# pprint.pprint(results)
print(results)