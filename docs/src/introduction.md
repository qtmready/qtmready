# ctrlplane.ai

Because an infraless movement is finally possible.

> Note: This is alpha-quality software.The surface API is subject to change without warning.

## Raison d'être

Getting the deployment mechanism for a distributed system right is one of the most challenging aspects of software control. In addition to many nuanced failure modes, there are table stakes that are extermely hard to cordinate. e.g.

- [Velocity Controls](https://exp-platform.com/Documents/2019%20TongXiaSumitBhardwajPavelDmitrievAleksanderFabijan_Safe-Velocity-ICSE-SEI.pdf) (Ability to do a controlled rollout).
- Reliable Rollbacks (Restoring to last known good configuration).

Then we have more nuanced factors like

- Traffic Drainning (for how long does it wait for traffic to drain before it deploys).
- Make before Break (verifying that a newly deployed host is healthy before it moves on to the next one).
- Safety Invariants (ensuring deployments leave the service in a consistent state, even if some of them fail).

On top of it all, the system must have the ability to provide a manual override incase something is needed. Typically, these details are tuned through years and years of operational experience. A blue/green deployment avoids most of these concerns by not deploying on a live system. However, blue/green deployments are not always feasible (e.g. stateful services) and require specialized infrastructure.

To standarize a controlled software delivery mechanism, we have to identify the heads of the [cerberus](https://en.wikipedia.org/wiki/Cerberus) and the techniques available to us e.g.

- Continous Delivery
- Versioned Infrastructure
- Monitoring

### Continous Delivery

Continous delivery has the most mature ecosystem of the above three. From container specific tools e.g. [argo](https://argo-cd.readthedocs.io/en/stable/), [dynatrace's keptn](https://keptn.sh), [tekton](https://tekton.dev), to generic platforms like [AWS Code Pipeline](https://aws.amazon.com/codepipeline/), [Google Cloud Build](https://cloud.google.com/build), [CicleCI](https://circleci.com), [Github Actions](https://github.com/features/actions) etc. You get the idea! I can go on forever.

> The bottom line is, it is a vast problem space and still being attacked with some novel approaches e.g. [Dagger](https://dagger.io) or [Acorn](https://acorn.dev).

As far as ctrlplane.ai is concerned, we are not concerned with the _building_ part. We leverage GitOps and we start after a new artifact is built to an artifact registry. **_We version it_**, and take control of the rest of the delivery process. Our version acts as control flow for the entire rollout process. In our case, the artifact for time being is OCI compatible image. In future, we might look at WASI or MicroVMs/UniKernels but for now, we are narrowing our scope to OCI.

### Versioned Infrastructure

With rapid adoption of cloud, critical infrastructure is what you can call "software defined infra". The success of tools like [terraform](https://terraform.io), [pulumi](https://pulumi.com) or [crossplane](https://crossplane.io) or [Amazon CDK](https://aws.amazon.com/cdk) have attacked the problem of automating the infrastructure and have provided a solid foundation for versioned infra. But these tools need stitching with deployment pipelines & feedback controls need to be built on top of it. Of the market data we have collected, the most forward looking shops are using pulumi or amazon CDK but terraform has the most mind share. _crossplane_, although a novel idea to use kubernetes CRDs are cloud resources, hasn't really taken off so far.

> One of the most revolutionary approach IMHO is by [Wing](https://winglang.io) i.e take the ifrastructure problem and turn it on its head and by making it a compiler problem. But it is still in its design phase, and the market will will decide the ultimate vote.

Our bet however is immutable versioned infra for each artifact from the continous delivery pipeline. For the time being, we are only handling stateless applications. Versioning a database migration is a seperate problem and needs more bandwidth.

### Monitoring

This is where things really get interesting and almost all the eggs lie in this basket. [Historically](https://failingfast.io/opentelemetry-observability/), for observability, each vendor had its own format i.e. [DataDog](https://datadoghq.com), [New Relic](https://newrelic.com), [Sentry](https://sentry.io) but recently, thanks to the merger between [OpenTracing](https://opentracing.org) & [OpenCencus](https://opencensus.org), we now have a [vendor neutral](https://failingfast.io/opentelemetry/) observability standard and saner heads in the industry have prevailed. All the aforementioned, thankfully, are jumping onto the [OpenTelemetry](https://opentelemetry.org) train. Major cloud providers either support opentelemetry out of the box or are in the process of adding support for it.

## Thesis

Our thesis is that, if the stack's performance output is a standard spec, one can design an analytics engine that will work anywhere. This is regardless of the cloud provider and application stack you choose. Using distributed traces and traffic paths, it is possible to identify inter-service dependencies in a distributed application stack. With these identified service dependencies, we may be able to create smart version sets containing only the dependent services to minimize the blast radius. By combining immutable infrastructure with this design, we will potentially be able to build a cloud-native, vendor neutral, release rollout engine where developers won't have to worry about their applications being down ever again when releasing new code.
