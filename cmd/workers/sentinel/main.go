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

func main() {}
