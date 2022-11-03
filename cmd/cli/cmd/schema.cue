

#Version: string | *"0.1"

#ServicePort: {
  port: number,
  protocol: string | int,
}

#Resource: {
  provider: string | *"aws" | "gcp" | "azure" | "do",
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
