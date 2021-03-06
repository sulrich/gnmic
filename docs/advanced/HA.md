## High Availability

`gnmic` can run in high availability mode to protect against gNMI connections loss. 
This is achieved by running multiple instances of `gnmic` loaded with the same configuration. 
In order to loadshare the targets connections between the different instances, each `gnmic` instance uses ephemeral key locks in a configured KV store ( such as [`Consul`](https://www.consul.io/)) to declare ownership over a specific target.

The Locker configuration is as simple as:

```yaml
locker:
  # type of locker, only consul is supported currently
  type: consul
  # address of the locker server
  address: localhost:8500
  # session-ttl, session time-to-live after which a session is considered 
  # invalid if not renewed
  session-ttl: 10s
  # delay, a time duration (0s to 60s), in the event of  a session invalidation 
  # consul will prevent the lock from being acquired for this duration.
  # The purpose is to allow a gnmic instance to stop active subscriptions before another one takes over.
  delay: 15s
  # retry-timer, wait period between retries to acquire a lock 
  # in the event of client failure, key is already locked or lock lost.
  retry-timer: 2s
  # renew-period, session renew period, must be lower that session-ttl. 
  # if the value is greater or equal than session-ttl, is will be set to half of session-ttl
  renew-period: 5s
  # debug, enable extra logging messages
  debug: false
```

A `gnmic` instance creates gNMI subscriptions only towards targets for which it acquired locks. It is also responsible for maintaining that lock for the duration of the subscription.
In the event of connection loss, the ephemeral lock expires leaving the opportunity for another `gnmic` instance to acquire the lock and re-create the gNMI subscription.

<div class="mxgraph" style="max-width:100%;border:1px solid transparent;margin:0 auto; display:block;" data-mxgraph="{&quot;page&quot;:12,&quot;zoom&quot;:1.4,&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;nav&quot;:true,&quot;check-visible-state&quot;:true,&quot;resize&quot;:true,&quot;url&quot;:&quot;https://raw.githubusercontent.com/karimra/gnmic/diagrams/diagrams//locking.drawio&quot;}"></div>

<script type="text/javascript" src="https://cdn.jsdelivr.net/gh/hellt/drawio-js@main/embed2.js?&fetch=https%3A%2F%2Fraw.githubusercontent.com%2Fkarimra%2Fgnmic%2Fdiagrams%2F/locking.drawio" async></script>

## Scalability

Using the same above-mentioned locking mechanism, `gnmic` can horizontally scale the number of supported gNMI connections distributed across multiple `gnmic` instances.

The collected gNMI data can then be aggregated and made available through any of the running `gnmic` instances, regardless of whether that instance collected the data from the target or not.

The data aggregation is done by chaining `gnmic` [outputs](multi_outputs/output_intro.md) and [inputs](inputs/input_intro.md) to build a gNMI data pipeline.

In the diagram below, the `gnmic` instances on the left and right side of NATS server can be identical.

<div class="mxgraph" style="max-width:100%;border:1px solid transparent;margin:0 auto; display:block;" data-mxgraph="{&quot;page&quot;:12,&quot;zoom&quot;:1.4,&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;nav&quot;:true,&quot;check-visible-state&quot;:true,&quot;resize&quot;:true,&quot;url&quot;:&quot;https://raw.githubusercontent.com/karimra/gnmic/diagrams/diagrams//scalability.drawio&quot;}"></div>

<script type="text/javascript" src="https://cdn.jsdelivr.net/gh/hellt/drawio-js@main/embed2.js?&fetch=https%3A%2F%2Fraw.githubusercontent.com%2Fkarimra%2Fgnmic%2Fdiagrams%2F/scalability.drawio" async></script>