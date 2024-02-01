# Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
#
# This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
# found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
# THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
#
# The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
# portions of the software.
#
# Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
# SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
# SOFTWARE.
#
# Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
# CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
# ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

import requests
import json

stack_name = "quantm poc"

#Getting Access Token
print("Logging in")
url = "http://localhost:8000/auth/login"
payload = {
    "email": "mahad@breu.io",
    "password": "pass123"
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
    "provider_id": "684576037",
    "default_branch": "main",
    "name": "cargoflo",
    "is_monorepo": True
}
repo_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#Creating Resource
print("Creating Resource")
url = "http://localhost:8000/core/resources"
payloadx = {
    "Name": "CloudRun_CargoFlo",
    "provider": "GCP",
    "driver": "cloudrun",
    "stack_id": stack_id,
    "Config": '{"properties":{"generation":"second-generation","cpu":"2000m","memory":"1024Mi"},"output":{"env":[{"url":"CloudRun_CargoFlo_URL"}]}}',
    "immutable": True
}
rsid = requests.request("POST", url, json=payloadx, headers=headers).json()["id"]

#creating Workload
print("Creating Workload")
url = "http://localhost:8000/core/workloads"
payload = {
    "Name": "api-quantm",
    "kind": "worker",
    "repo_id": repo_id,
    "repo_path": "https://github.com/breuHQ/cargoflo",
    "resource_id": rsid,
    "stack_id": stack_id,
    "builder": "",
    "container": '{"image": "europe-west3-docker.pkg.dev/cargoflo-dev-400720/cloud-run-source-deploy/cargoflo/api"}'
}
workload_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]

#creating BluePrint
print("Creating BluePrint")
url = "http://localhost:8000/core/blueprints"
payload = {
    
    "name" : "CargoFlo blueprint",
    "stack_id" : stack_id,
    "rollout_budget" : "300",
    "regions" : { "gcp": ["europe-west3"], "aws": [], "azure": [], "default": [] },
    "provider_config" : '{"project": "cargoflo-dev-400720"}'
}
blueprint_id = requests.request("POST", url, json=payload, headers=headers).json()["id"]
