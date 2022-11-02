// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

// A sentinel runs on customer premise. It acts as tunnel to ctrlplane.ai main server.
// It is responsible for:
// - cloning the repo
// - identifying the builder (dockerfile, packer, etc)
// - figuring out the required infra components
// - provisioning and versioning the infra components
// - deploying the infra
// - deploying the stack
// - setting up the stack
//
// The only communication channel between sentinel and the mothership is via
// temporal server. Each sentinel has its own namespae on the server.
package main
