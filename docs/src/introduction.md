# ctrlplane.ai

Because an infraless movement is finally possible.

> Note: This is alpha-quality software.The surface API is subject to change without warning.

## Raison d'être

Getting the deployment mechanism for a distributed system right is one of the most challenging aspects of software control. In addition to many nuanced failure modes, there are table stakes that are hard to extermely hard to cordinate. e.g.

- [Velocity Controls](https://exp-platform.com/Documents/2019%20TongXiaSumitBhardwajPavelDmitrievAleksanderFabijan_Safe-Velocity-ICSE-SEI.pdf) (Ability to do a controlled rollout)
- Reliable Rollbacks (Restoring to last known good configuration)

Then we have more nuacned factors like

- Traffic Drainning (for how long does it wait for traffic to drain before it deploys)
- Make before Break (verifying that a newly deployed host is healthy before it moves on to the next one)
- Safety Invariants (ensuring deployments leave the service in a consistent state, even if some of them fail)

On top of it all, the system must have the ability to provide a manual override incase something is needed. Typically, these details are tuned through years and years of operational experience. A blue/green deployment avoids most of these concerns by not deploying on a live system. However, blue/green deployments are not always feasible (e.g. stateful services) and require specialized infrastructure.

To standarize a controlled software delivery mechniasm, we have to identify mark the heads of the [cerberus](https://en.wikipedia.org/wiki/Cerberus) and the techniques available to us e.g.

- Continous Delivery
- Versioned Infrastructure
- Monitoring

### Continous Delivery

Continous delivery has the most mature solutions, from container specific tools e.g. [argo](https://argo-cd.readthedocs.io/en/stable/), [dynatrace's keptn](https://keptn.sh) or [tekton](https://tekton.dev) to generic solutions like [AWS Code Pipeline](https://aws.amazon.com/codepipeline/), [Google Cloud Build](https://cloud.google.com/build), [CicleCI](https://circleci.com) or [Github Actions](https://github.com/features/actions). Cotainer specific tools have tried to solve Green/Blue Deployments & Keptn even have quality gates but at the time of writing this, their scope is limited to kubernetes platform. If we look at the automation part of the continous delivery, there are some novel approaches e.g. [Dagger](https://dagger.io).

As far as controlplane.ai is concerned, we are not aiming to solve this piece. We aim to leverage this. We might support dagger as an SDK in the future, but for now, we aim to provide a container as the last step of the deployment pipeline.

### Versioned Infrastructure

With rapid adoption of cloud, critical infrastructure is what you can call "software defined infra". The success of tools like [terraform](https://terraform.io), [pulumi](https://pulumi.com) or [crossplane](https://crossplane.io) have attacked the problem of "versioned infra" and are somehwat successful. But these tools need stitching and deployment pipelines & feedback controls need to be built on top of it. Of the market data we have collected, the most forward looking shops are using pulumi, but terraform.io has the most mind share. crossplane, although a novel idea to use kubernetes CRDs are cloud resources, hasn't really taken off so far. As far as the terraform is concerned, most shops haven't really plugged the infra into CI/CD pipelines as yet.

### Monitoring

This is where things really get interesting.

## Contributing

How to work as team is summarized at [CONTRIBUTING](./CONTRIBUTING.md).

## License

The code is made available under _Breu Community License Agreement_.

See [License](./LICENSE) for details.
