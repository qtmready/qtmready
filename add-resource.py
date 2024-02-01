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
