# Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
#Version: string | *"0.1"

#ServicePort: {
  port: number,
  protocol: string | int,
}

#Resource: {
  drvier: string | *"aws" | "gcp" | "azure",
  name: string,
  exports: {
    [string]: string,
  },
}

#Service: {
  name: string,
  image: {
    repository: string,
    name: string,
  }
  ports: [...#ServicePort],
}


version: #Version
resources: [...#Resource]
services: [...#Service]
