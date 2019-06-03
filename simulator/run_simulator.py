import simulator
import os
import json
import subprocess
import yaml

def getFilePath(i):
    return "result_data/result_"+str(i)+".yaml"

# DIR = 'result_data' #要统计的文件夹
# file_num = len([name for name in os.listdir(DIR) if os.path.isfile(os.path.join(DIR, name))])
# print (file_num )


def extract_output(stdout):
    out = bytes.decode(stdout)
    out = out.split('\n')
    p = (out[0].split(':'))[1].strip()
    p = p.replace("'", "\"")
    p = json.loads(p)

    t = (out[1].split(':'))[1]
    t = t.strip()
    t = float(t)
    return p, t

# results = []

# for i in range(file_num):
#     print(i)
#     d = dict()
#     file_name = getFilePath(i)

#     fp = open(file_name, 'r')
#     inputs = yaml.load(fp)
#     fp.close()

#     pods = inputs['Pods']
#     nodes = inputs['Nodes']
#     if len(pods) > 13:
#         continue

#     d['pods'] = pods
#     print('pod len:', len(pods))
#     d['nodes'] = nodes
#     print('node len:', len(nodes))

#     process = subprocess.run(
#         ['python3', 'simulator.py', file_name],
#         stdin=subprocess.PIPE, stdout=subprocess.PIPE
#     )
#     p, t = extract_output(process.stdout)

#     d['result'] = p
#     d['time'] = t
#     d['index'] = i
#     results.append(d)

# import pprint
# # pprint.pprint(results)
# print(len(results))


# with open("simulator_results.json", 'w') as f:
#     json.dump(results, f)

# file_name = 'result_data/result_4.yaml'