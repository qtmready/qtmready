# Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
#
# Functional Source License, Version 1.1, Apache 2.0 Future License
#
# We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
# is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
# the Software under the Apache License, Version 2.0, in which case the following will apply:
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
# the License.
#
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
# an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.

import requests
import json

stack_name = "quantm mutli repo poc"

#Getting Access Token
print("Logging in")
url = "https://qw4oxd513shx.share.zrok.io/auth/login"
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
url = "https://qw4oxd513shx.share.zrok.io/core/stacks"
payload = {"name": stack_name}
headers = {
    "Content-Type": "application/json",
    "Authorization": "Token " + access_token
}
stack_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#Creating Repo
print("Creating Repo")
url = "https://qw4oxd513shx.share.zrok.io/core/repos"
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
url = "https://qw4oxd513shx.share.zrok.io/core/repos"
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
# url = "https://qw4oxd513shx.share.zrok.io/core/resources"
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
# url = "https://qw4oxd513shx.share.zrok.io/core/workloads"
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
# url = "https://qw4oxd513shx.share.zrok.io/core/blueprints"
# payload = {

#     "name" : "CargoFlo blueprint",
#     "stack_id" : stack_id,
#     "rollout_budget" : "300",
#     "regions" : { "gcp": ["europe-west3"], "aws": [], "azure": [], "default": [] },
#     "provider_config" : '{"project": "cargoflo-dev-400720"}'
# }
# blueprint_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]
