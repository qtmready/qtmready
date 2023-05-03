# Core

## Workflows

### Pull Request
.
```mermaid
sequenceDiagram
  autonumber
  actor UR as User
  participant WH as API :: Webhook RX
  participant gwfpr as github::workflow::OnPullRequest
  participant cwfpr as core::workflow::OnPullRequest
  participant mwf as core::workflow::MutexWorkflow

  activate cwfpr
  activate mwf
  activate UR
   UR ->> WH: PR generated
   WH ->> gwfpr: start
   activate gwfpr
   WH ->> UR: return

     gwfpr ->> gwfpr: get core repo(provider, provider ID)
     gwfpr ->> gwfpr: get Stack(Stack ID)
     gwfpr ->> cwfpr: signal PR event
     Note left of cwfpr: core repo id, sender workflowID
    deactivate gwfpr
    cwfpr ->> mwf: signal acquire mutex
    Note left of mwf: resource ID

    mwf ->> cwfpr: Lock acquired
    cwfpr ->> cwfpr: execute critical section
    cwfpr ->> mwf: release lock 
    
  deactivate cwfpr
  deactivate mwf
```

### Create Stack API
```mermaid
sequenceDiagram
    autonumber
    actor u as user
    participant api as API
    participant cwfpr as core::workflow::OnPullRequest
    participant mwf as core::workflow::MutexWorkflow

  activate u
    u ->>+ api: create Stack
    api ->>+ cwfpr: start
    api ->>- u: return

    cwfpr ->> mwf: start child workflow
    deactivate cwfpr
    <!-- deactivate mwf -->
```
