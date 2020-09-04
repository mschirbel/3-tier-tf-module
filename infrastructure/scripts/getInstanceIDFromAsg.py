import sys json

def read_in():
    return {x.strip() for x in sys.stdin}

def get_instances(asg_name, region):
    asg_response = asg_client.describe_auto_scaling_groups(AutoScalingGroupNames=[asg_name], Region=region)
    instance_ids = []
    for i in asg_response['AutoScalingGroups']:
        for k in i['Instances']:
            instance_ids.append(k['InstanceId'])
    return instance_ids

def main():
    lines = read_in()
    for line in lines:
        if line:
            jsondata = json.loads(line)
    ids = get_instances(jsondata[asg_name], jsondata[region])
    sys.stdout.write(json.dumps(ids))

if __name__ == '__main__':
    main()