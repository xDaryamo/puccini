TOSCA BPMN Profile Examples
===========================

Integration with [BPMN](https://www.omg.org/spec/BPMN/) middleware.

Features
--------

Two features (introduced in TOSCA 1.1) are supported:

### Workflows

TOSCA declarative workflows are translated into BPMN processes. Because a TOSCA workflow is
essentially a graph of steps with sequential and parallel sections, in BPMN we must represent
the graph using parallel gateways, diverging or converging as the case may be, as well as
conditional gateways to represent step success or failure. The JavaScript analyzes the graph and
inserts the appropriate gateways between the steps.

Each step in TOSCA comprises zero or more activities that should happen in sequence. In BPMN, all
the activities in the step become a single `scriptTask` entity. For now, we create a script made of
pseudo-code that calls these activities. A complete solution would require a BPM orchestration
environment and real code that would actually call node instances deployed in a cloud.   

Once the BPMN process is imported into BPM software, you may include this process as a sub-process
within other BPM processes. The workflow may or may not hand control back to another sub-process,
so that it may or may not be a continuation of a control loop.

### Policies and Policy Triggers

Policy triggers are also translated into BPMN processes. Because the trigger must be executed from
within an orchestrator on node instances deployed in a cloud, within configurable time intervals
or schedules, this BPM process essentially hands over control of the loop to the orchestrator. By
launching a new sub-process when triggered, control is handed back to the business process: an open
loop.

A single `scriptTask` entity is created for each target node of the policy, and all are executed in
parallel using diverging/converging parallel gateways. A conditional gateway at the convergence is
used to launch a new sub-process if any of the tasks succeed. Again, the script is made of
pseudo-code that would call these operations within a BPM orchestration environment.

Open Loop Example
-----------------

Included is an example TOSCA service template, [`open-loop.yaml`](open-loop.yaml). This example
demonstrates an open loop policy, `notify_on_high_load`, which has a trigger that runs an
operation to get the CPU load on Compute nodes. If this operation returns true then a BPM process
named `NotifyUser` would be launched.

Also included is a TOSCA workflow named `backup`, which comprises a step graph that calls an
operation on an interface while making sure to set node states, notify on failure, etc. This
generated BPMN process can be executed on its own, or included as a sub-process within larger
business processes. Because it does not hand control back to any other process when done, it
represents an end event within a control loop. 

To generate the BPMN:

    puccini-tosca compile open-loop.yaml | puccini-clout scriptlet exec bpmn.generate -o open-loop.bpmn2

Also included is [`open-loop-design.bpmn2`](open-loop-with-diagram.bpmn2), which is the same file
with added diagram information so that it would appear more nicely in a BPMN GUI. We used the
[Eclipse BPMN2 modeler](https://www.eclipse.org/bpmn2-modeler/) to edit the diagram.

You can import either file into your BPM software. Tested with [jBPM](https://www.jbpm.org/) 7.8.0.
