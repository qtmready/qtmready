import requests
import json

stack_name = "quantm mutli repo poc"

#Getting Access Token
print("Logging in")
url = "https://api.breu.ngrok.io/auth/login"
payload = {
    "email": "umer@breu.io",
    "password": "pass123"
}
headers = {"Content-Type": "application/json"}
res = requests.request("POST", url, json=payload, headers=headers)
print(res)
access_token = res.json()["access_token"]

#Creating Stack
print("Creating Stack")
url = "https://api.breu.ngrok.io/core/stacks"
payload = {"name": stack_name}
headers = {
    "Content-Type": "application/json",
    "Authorization": "Token " + access_token
}
stack_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#Creating Repo
print("Creating Repo")
url = "https://api.breu.ngrok.io/core/repos"
payload = {
    "stack_id": stack_id,
    "provider": "github",
    "provider_id": "746408294",
    "default_branch": "main",
    "name": "Quantm-testing-2",
    "is_monorepo": False
}
repo_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#Creating Repo
print("Creating Repo")
url = "https://api.breu.ngrok.io/core/repos"
payload = {
    "stack_id": stack_id,
    "provider": "github",
    "provider_id": "746408201",
    "default_branch": "main",
    "name": "Quantm-testing-1",
    "is_monorepo": False
}
repo_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

# #Creating Resource
# print("Creating Resource")
# url = "https://api.breu.ngrok.io/core/resources"
# payloadx = {
#     "Name": "CloudRun_CargoFlo",
#     "provider": "GCP",
#     "driver": "cloudrun",
#     "stack_id": stack_id,
#     "Config": '{"properties":{"generation":"second-generation","cpu":"2000m","memory":"1024Mi"},"output":{"env":[{"url":"CloudRun_CargoFlo_URL"}]}}',
#     "immutable": True
# }
# rsid = requests.request("POST", url, json=payloadx, headers=headers).json()["id"]

# #creating Workload
# print("Creating Workload")
# url = "https://api.breu.ngrok.io/core/workloads"
# payload = {
#     "Name": "api-quantm",
#     "kind": "worker",
#     "repo_id": repo_id,
#     "repo_path": "https://github.com/breuHQ/cargoflo",
#     "resource_id": rsid,
#     "stack_id": stack_id,
#     "builder": "",
#     "container": '{"image": "europe-west3-docker.pkg.dev/cargoflo-dev-400720/cloud-run-source-deploy/cargoflo/api"}'
# }
# workload_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

# #creating BluePrint
# print("Creating BluePrint")
# url = "https://api.breu.ngrok.io/core/blueprints"
# payload = {

#     "name" : "CargoFlo blueprint",
#     "stack_id" : stack_id,
#     "rollout_budget" : "300",
#     "regions" : { "gcp": ["europe-west3"], "aws": [], "azure": [], "default": [] },
#     "provider_config" : '{"project": "cargoflo-dev-400720"}'
# }
# blueprint_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]
