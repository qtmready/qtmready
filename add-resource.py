# Copyright © 2023, Breu, Inc. <info@breu.io>
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

headers = {
    "Content-Type": "application/json",
    "Authorization": "Token " + access_token
}

#Creating Resource
print("Creating Resource")
url = "http://localhost:8000/core/resources"
config_resources = {
    "name": "api-quantm",
    "location": "europe-west3",
    "launch_stage": "BETA",
    "ingress": "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER",
    "labels": {
        "app": "cargoflo",
        "component": "api",
        "resource": "cloudrun"
    },
    "template": {
        "execution_environment": "EXECUTION_ENVIRONMENT_GEN2",
        "scaling": {
            "min_instance_count": "0",
            "max_instance_count": "3"
        },
        "containers": {
            "image": "europe-west3-docker.pkg.dev/cargoflo-dev-400720/cloud-run-source-deploy/cargoflo/api",
            "ports": {
                "container_port": "8000"
            },
            "resources": {
                "cpu_idle": "true"
            },
            "env": [
                {
                    "name": "CARGOFLO_DEBUG",
                    "value": "false"
                },
                {
                    "name": "CARGOFLO_TEMPORAL_HOST", 
                    "value": "10.10.0.3"
                },
                {
                    "name": "CARGOFLO_TEMPORAL_NAMESPACE", 
                    "value": "cargoflo-dev"
                },
                {
                    "name": "CARGOFLO_DB_HOST", 
                    "value": "110.69.49.8"
                },
                {
                    "name": "CARGOFLO_DB_NAME", 
                    "value": "cargoflo"
                },
                {
                    "name": "CARGOFLO_DB_USER", 
                    "value": "cargoflo"
                },
                {
                    "name": "CARGOFLO_DB_PASS", 
                    "value": "cargoflo"
                },
                {
                    "name": "CARGOFLO_DB_MAX_OPEN_CONNECTIONS", 
                    "value": "25"
                }
            ],
            "volume_mounts" : {
                "name": "cloudsql",
                "mount_path": "/cloudsql"
            }
        },
        "vpc_access": {
            "egress": "PRIVATE_RANGES_ONLY",
            "network_interfaces": {
                "network": "cargoflo-dev-8abebbf2",
                "subnetwork" : "europe-west3-cargoflo-dev-8abebbf2"
            },
        },
        "volumes" : {
            "name": "cloudsql",
            "cloud_sql_instance": {
                "instances" : [
                    {
                        "name": "cargoflo-dev-400720:europe-west3:default-europe-west3-cargoflo-dev-8abebbf2"
                    }
                ]
            }
        }
    }
}
payloadx = {
    "Name": "CloudRun_CargoFlo",
    "provider": "GCP",
    "driver": "cloudrun",
    "stack_id": "36dc5aec-2f7b-402f-9b09-0e02e9bbe539",
    "Config": json.dumps(config_resources),
    "immutable": True
}
rsid = requests.request("POST", url, json=payloadx, headers=headers).json()
print(rsid)
