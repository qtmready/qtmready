import requests
import json

stack_name = "quantum poc"

#Getting Access Token
print("Logging in")
url = "http://localhost:8000/auth/login"
payload = {
    "email": "amna@breu.io",
    "password": "password"
}
headers = {"Content-Type": "application/json"}
access_token = requests.request("POST", url, json=payload, headers=headers).json()["access_token"]

#Creating Stack
print("Creating Stack")
url = "http://localhost:8000/core/stacks"
payload = {"name": stack_name}
headers = {
    "Content-Type": "application/json",
    "Authorization": "Token " + access_token
}
stack_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#Creating Repo
print("Creating Repo")
url = "http://localhost:8000/core/repos"
payload = {
    "stack_id": stack_id,
    "provider": "github",
    "provider_id": "648084184",
    "default_branch": "main",
    "name": "HelloWorld",
    "is_monorepo": True
}
repo_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#Creating Resource
print("Creating Resource")
url = "http://localhost:8000/core/resources"
payloadx = {
    "Name": "CloudRun_HelloWorld",
    "provider": "GCP",
    "driver": "cloudrun",
    "stack_id": stack_id,
    "Config": '{"properties":{"generation":"second-generation","cpu":"2000m","memory":"1024Mi"},"output":{"env":[{"url":"CloudRun_HelloWorld_URL"}]}}',
    "is_immutable": True
}
rsid = requests.request("POST", url, json=payloadx, headers=headers).json()["id"]

#creating Workload
print("Creating Workload")
url = "http://localhost:8000/core/workloads"
payload = {
    "Name": "helloworld",
    "kind": "worker",
    "repo_id": repo_id,
    "repo_path": "https://github.com/amnabreu/HelloWorld",
    "resource_id": rsid,
    "stack_id": stack_id,
    "builder": "",
    "container": '{"image": "asia-southeast1-docker.pkg.dev/breu-dev/ctrlplane/helloworld"}'
}
workload_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#creating BluePrint
print("Creating BluePrint")
url = "http://localhost:8000/core/blueprints"
payload = {
    
    "name" : "Helloworld blueprint",
    "stack_id" : stack_id,
    "rollout_budget" : "300",
    "regions" : { "gcp": ["asia-southeast1"], "aws": [], "azure": [], "default": [] },
    "provider_config" : '{"project": "breu-dev"}'
}
blueprint_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]
