import jenkins
import sys

# 默认Jenkins项目
project_name = "dev"

# 分支名称
branch_name = ""
dev_name = ""
l = len(sys.argv)
if l == 1:
    print("输入有误")
    exit()

if l > 1:
    project_name = sys.argv[1]

if l > 2:
    branch_name = sys.argv[2]

if l > 3:
    dev_name = sys.argv[3]

print("argv is:", project_name, branch_name, dev_name)

#配置不同Jenkins
if project_name == "dev":
    # 定义远程的jenkins master server的url,以及port | 定义用户的User Id 和 API Token,获取方式同上文
    jenkins_server_url = 'http://你的Jenkins服务器IP:8080/'
    user_id = '登陆用户abc'
    api_token = '登陆用户设置设置token'
elif project_name == "test":
    jenkins_server_url = 'http://你的Jenkins服务器IP:8080/'
    user_id = '登陆用户abc'
    api_token = '登陆用户设置设置token'
else:
    print("项目配置不正确")
    exit()

# 对应job命名
job_name = branch_name


# 自定义参数
param_dict = {}
is_more = False
if job_name == "school":
    if dev_name == "dev":
        is_more = True
        param_dict = {"parameter": [{"name": "Branch", "value": "*/dev"}]}

# 实例化jenkins对象,连接远程的jenkins master server
server = jenkins.Jenkins(jenkins_server_url, username=user_id, password=api_token)

"""
POST http://你的Jenkins服务器IP:8080/view/hd/job/school/build?delay=0sec

name: Branch
value: */dev
statusCode: 303
redirectTo: .
Jenkins-Crumb: 1b8ca0e01b8af82c326a2ef7540189c7e759a00c79b12e492cf131f1c7ff6924
json: {"parameter": {"name": "Branch", "value": "*/dev"}, "statusCode": "303", "redirectTo": ".", "Jenkins-Crumb": "1b8ca0e01b8af82c326a2ef7540189c7e759a00c79b12e492cf131f1c7ff6924"}
Submit: 开始构建
"""

if is_more == True:
    # String参数化构建job名为job_name的job, 参数param_dict为字典形式,如：param_dict= {"param1"：“value1”, “param2”：“value2”}
    qNum = server.build_job(job_name, parameters=param_dict)
else:
    # 构建job名为job_name的job（不带构建参数）
    qNum = server.build_job(job_name)
# print("当前队列名称：", qNum)

# 获取job名为job_name的job的最后次构建号
build_number = server.get_job_info(job_name)['lastBuild']['number']
# 获取job名为job_name的job的某次构建的执行结果状态
print("project_name:[{}],job_name:[{}],build_number:[{}],building:[{}],result:[{}]".format(
    project_name,
    job_name,
    build_number,
    server.get_build_info(job_name, build_number)['building'],
    server.get_build_info(job_name, build_number)['result']
))