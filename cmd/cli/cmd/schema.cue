// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 



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
