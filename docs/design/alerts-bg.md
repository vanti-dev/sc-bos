# Alerts v2 Background Doc

Alerts in sc-bos are currently represented via these mechanisms:

1. The Status trait allows a device to report one or more problems with itself. We've implemented this in a number of
   devices to represent everything from "connection timeout" to "value out of range".
2. The Ops UI notification page and associated AlertApi. This is backed by a database of faults, typically filled via an
   automation that watches all devices that implement the status trait and any time they changes insert a row (or update
   a row) in that alerts table.

   We also support manual entry of rows into the table and limited workflow for alerts via acknowledgement.
3. We have an Emergency trait whose intent is to model fire alarms and other evacuation events.
   This trait is not widely used and has limited integration with the Ops UI and alerts system.
4. We have automations that watch for status changes and send emails to preconfigured addresses with both digests and
   live-ish updates.

## What people expect from alerts

We've heard many times that what people want from the alert system is to be able to see which parts of their building
needs attention. This could be a sensor that is reporting a value out of range (like an empty water tank) or a test that
failed (like an emergency light battery test) or an active alarm (like a fire alarm).

People have also asked to see alert status over time.
We believe the intent here is to identify areas of the building that need broader attention, either a redesign or
improved workflow.
For example a repeated alert that a fan is overheating might indicate a fault or under-specified device.

People come from a world of BMS or Fire systems which bring alerts front and centre.

People also expect to be able to prioritise their reaction based on the device that is raising the alert.
Devices that are tied to life safety need immediate attention.

## Things SC does that we want to keep

A single location where all alerts can be seen (and exported).
Be able to filter alerts to see only those that I'm looking for - while we have limited capability for this, we don't
want it to get worse. So we still need to filter by floor, device, severity, etc.
Be able to embed sub-views of alerts into dashboard pages, e.g. Floor 4 alerts, or Teal Room alerts.
See summary of how many alerts require attention, again we aren't good at 'needs attention' bit, but we do show how many
alerts are active and unacknowledged.

## Alerts in other systems

### BACnet

Each object can have a number of properties that together report on the status of the object.

**Status_Flags** is a read-only aggregate property that combines the status of other properties into a single value.

> This read-only property, of type BACnetStatusFlags, represents four Boolean flags that indicate the general "health"
> of an analog input. Three of the flags are associated with the values of other properties of this object. A more
> detailed status could be determined by reading the properties that are linked to these flags. The relationship between
> individual flags is not defined by the protocol. The four flags are
>
> `{IN_ALARM, FAULT, OVERRIDDEN, OUT_OF_SERVICE}`
>
> where:
>
> - IN_ALARM Logical FALSE (0) if the Event_State property has a value of NORMAL, otherwise logical TRUE (1).
> - FAULT Logical TRUE (1) if the Reliability property is present and does not have a value of NO_FAULT_DETECTED,
    otherwise logical FALSE (0).
> - OVERRIDDEN Logical TRUE (1) if the point has been overridden by some mechanism local to the BACnet device. In
    this context "overridden" is taken to mean that the Present_Value and Reliability properties are no
    longer tracking changes to the physical input. Otherwise, the value is logical FALSE (0).
> - OUT_OF_SERVICE Logical TRUE (1) if the Out_Of_Service property is set to TRUE, otherwise logical FALSE (0).

**Event_State** is a read-only property indicating the current fault state of the object. The value is updated based on
the reliability and current value for the object and any configuration that indicates normal bounds for the value.

> This read-only property, of type BACnetEventState, is included in order to provide a way to determine whether this
> object has an active event state associated with it (see Clause 13.2.2.1). If the object supports event reporting,
> then the Event_State property shall indicate the event state of the object. If the object does not support event
> reporting, then the value of this property shall be NORMAL.

Possible values for this property are:

- NORMAL - things are working correctly
- FAULT - reliability is not NO_FAULT_DETECTED
- OFFNORMAL - something isn't right, maybe the value isn't what was expected
- HIGH_LIMIT - the value is above the maximum normal value
- LOW_LIMIT - the value is below the minimum normal value
- LIFE_SAFETY_ALARM - the value is in one of the specified life safety alarm states

**Reliability** indicates whether the BACnet representation of the object is tied to the physical object. Basically how
reliable is the information being read from BACnet. See 12.1.8 for an overview.

> This property, of type BACnetReliability, provides an indication of whether the Present_Value or the operation of the
> physical input in question is "reliable" as far as the BACnet device or operator can determine and, if not, why.
>
> If a fault algorithm is applied, then this property shall be the pCurrentReliability parameter for the object's fault
> algorithm.

Possible values for this property are:

- no-fault-detected (0),
- no-sensor (1),
- over-range (2),
- under-range (3),
- open-loop (4),
- shorted-loop (5),
- no-output (6),
- unreliable-other (7),
- process-error (8),
- multi-state-fault (9),
- configuration-error (10),
- -- enumeration value 11 is reserved for a future addendum
- communication-failure (12),
- member-fault (13),
- monitored-object-fault (14),
- tripped (15),
- lamp-failure (16),
- activation-failure (17),
- renew-dhcp-failure (18),
- renew-fd-registration-failure (19),
- restart-auto-negotiation-failure (20),
- restart-failure (21),
- proprietary-command-failure (22),
- faults-listed (23),
- referenced-object-fault (24),
- multi-state-out-of-range (25),

**Out_Of_Service** indicates whether the object is currently in service (whatever that means). The docs suggest that
when the device is out of service the BACnet values are no longer tied to the physical properties of the object.
Additionally the BACnet object becomes writable (when it used to be read only) allowing testing to be performed on
connected systems.

> This property, of type BOOLEAN, is an indication whether (TRUE) or not (FALSE) the physical input that the object
> represents is not in service.
>
> When Out_Of_Service is TRUE:
> 1. the Present_Value property is decoupled from the physical input and will not track changes to the physical input;
> 2. the Reliability property, if present, and the corresponding state of the FAULT flag of the Status_Flags property
     shall be decoupled from the physical input;
> 3. the Present_Value property and the Reliability property, if present and capable of taking on values other than
     NO_FAULT_DETECTED, shall be writable to allow simulating specific conditions or for testing purposes;
> 4. other functions that depend on the state of the Present_Value or Reliability properties shall respond to changes
     made to these properties as if those changes had occurred in the physical input.

### DALI Emergency Lighting

From [the DALI website](https://www.dali-alliance.org/dali/emergency.html):

> Emergency lighting, which provides light when the mains supply fails, is a critical feature that is mandated by
> various regulations. “Self-contained” means that the battery – which provides power during an emergency – is inside,
> or placed next to, the luminaire.
>
> In many countries, there is a legal requirement for periodic testing of emergency lighting. DALI enables
> self-contained emergency tests to be automated, triggered by DALI commands or by an optional timer.
>
> Emergency control gear must implement both a function test and a duration test. The function test is a quick test of
> the battery, charging circuit, driver/relay and lamp, while the duration test ensures that the battery will be able to
> operate the lamp for the full rated duration (for example 1 hour or 3 hours).

The important part here is that test failures (functional and duration) should be surfaced as alarms associated with the
device.

Additional information (without the spec) can be found on the
[Beckhoff integration page](https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc2_dali/4345721611.html&id=6221676384172518865)
or in this introduction presentation by
[The DALI Alliance](https://www.dali-alliance.org/data/downloadables/2/4/1/210318_diia-seminar-on-dali-emergency-for-zvei_18-march-2021_final.pdf).

Emergency lights have two different tests:

- Functional: lamp operation, circuit faults, battery switch over
- Duration: tests that the battery can power the light for the full rated duration

Running tests against an emergency light is not a synchronous process.
There are factors that can cause the test to be postponed, such as when the light is in emergency mode (no mains power,
battery too low, test already running).
Similarly, the test can be aborted (I'm unsure if it'll be rescheduled automatically) if running and the mains power is
cut off.
Tests can also be manually stopped.
While functional tests can be quick (under a minute), durations tests can take a long time (up to 3 hours) to complete.
Tests can be scheduled to run at specific intervals, functional tests usually every month, duration tests every year.

The lamp can also report how long the battery lasted when a duration test is run.

Knowing that a test hasn't been performed when it should have been is also an alert state.
Regular testing is part of the regulations for emergency lighting.

Possible fault states are as follows, defined as a bitmask so multiple can be active at once:

```
Failure Status (bit set, only one of bits 1-5 + bits 6 or 7):
Error in the control gear circuit.
Battery operation time fault.
Battery fault.
Emergency lamp fault.
Timeout during function test. // These usually mean a timeout waiting for suitable conditions to run the test,
Timeout during duration test. //  not that the test took too long to run. E.G. tests aren't run during an emergency.
Function test failed.
Operating period test failed.
```

Indications of whether a test is running or can be run come from these properties:

```
Emergency Status (bit set):
Inhibit - Mains is reset, a physical or command INHIBIT has been set, suspends the normal->emergency switch for a period
          Tests can only be run when in normal mode, so inhibit also delay execution of tests.
FunctionTestDone - These are set when a test completes (successfully or not), see Failure States for the results.
DurationTestDone -  Only one of these are set at a time.
BatteryFull
FunctionTestPending - These indicate that the preconditions for running the test haven't yet been met.
DurationTestPending -  Maybe the battery is too low or the mains is off.
IdentificationActive
PhysicallySelected
```

```
Emergency Mode (bit set, but only one is every active):
Rest - I think is like _everything else_ rather than _inactive_. Causes the light to turn off even when mains is off.
Normal - Mains is on, light is in default state depending on type of fixture. Tests can be started.
Emergency - Mains is off, light is on. Tests can't be run. Running tests can transition to this state.
Extended Emergency - Acts like Emergency for a period of time after mains is restored.
Function Test In Progress - The light is running a functional test. Transitions to Emergency is mains is cut.
Duration Test In Progress - The light is running a duration test. Transitions to Emergency is mains is cut.
// these aren't really modes, more like flags:
Inhibit Signal - A physical switch is preventing emergency mode from activating.
Mains Active - Mains power is on.
```

The Emergency Status / FooTestDone bits are sticky, this means they will remain set until manually cleared.
DALI commands exist to clear this flag (set it back to 0).
Only successful test results can be cleared in this way, failures must be cleared by fixing the fault and rerunning the
test.
There is no 'test last run time' property, so discovering when (or how often) a test has been run is not possible by
querying the device.
A mechanism outside the DALI device needs to exist to record this, for example polling.

So in summary:

1. There are two types of test: Functional and Duration
2. Tests can be automatically scheduled or run manually
3. Tests can only be started when in `normal mode` and other conditions like battery charge are met
4. Tests that have been scheduled but not run are identified via Emergency Status / FooTestPending
5. When tests complete they are identified via Emergency Status / FooTestDone being set
6. The results of tests can be read from Failure Status / Foo test failed
7. Specific failure reasons are bits 1-5 from Failure Status
8. Knowing if a test has been run since we last asked the device is not possible without resetting the test done bits

### Gallagher Access Control System

Alerts in access control systems follow a different principle to device and system alerts, type generally represent
moments in time rather than states, though this isn't guaranteed.

In an ACS system you will see alerts like:

- Amy was denied access to Door 1 at 12:34
- Bob was granted access to Floor 2 through Door 3 (stairwell) at 12:35
- Communication with Floor 2 Zone 1 controller is down
- Administrator (Charlie) login at 12:36
- GF - BOH Door (Next to WC Toilet) has been forced - 105-1CS-GF Controller - at 12:37

Events in the log appear to be related by a few criteria:

- Event Type: things like "Access Granted", "Access Denied", "Communication Down", "Forced Door"
    - These seem very specific: Door open too long, Card Exit Granted, and so on
- Alarm Zone: this looks to be a controller reference, like "105-1CS-GF Controller 03" or "111-1CS-L13 Controller"
- Even Source: I think this is the unit that raised the event, "Elevator HLI 3" or "GF - BOH Door (Next to WC Toilet)"
  or "Level 13 - Gas Suppression Corridor Door"
- Cardholder: the person the system thinks is in ownership of the card that interacted with the system to produce the
  alert, "Contractor 05" or

## Idea grab bag

### Critical Devices or Critical Alerts

We've been asked to distinguish between critical (life safety) and non-critical alerts.
I think it would also make sense to highlight alerts/devices that contribute to conformance and regulations too.
We need to decide if alerts are critical or devices are critical.

At the very least we need to be able to filter alerts by criticality.

### Check Again Button

An ability to "check again" would be useful.
We'd need to be able to rate limit the checks, and not all checks can be done on demand.
Consider whether this would work for emergency lighting tests.

The above kind of implies that checks are something the alert system knows about.
We're getting into the realm of general purpose network/device health monitoring here.
At the very least a model of why an alert exists would be useful:

- The device is reporting a fault
- 'Foo Property' is out of range [1, 10]
- The functional test failed
- RPC requests to the device are failing

### Error Codes

What if kinds of fault are fully specified (or at least specified to be unique).
Rather than relying on text descriptions that only humans understand, we could model a way for devices/drivers to
identify alerts - like an error code, but one that works across different device types.
If you want to see all "Emergency Light Battery Faults" you can, or "temperature out of range".
We could specify a bunch of these based on what we know, but provide a namespaced way to add new ones with descriptions
and so on.

I feel like different types of alert would also have different supporting properties:

- Battery duration fault: wants 3h, lasted 1h 24m
- Temperature out of range: wants [1, 10], is 12

### Underlying Problem

A number of devices and points are accessed via a gateway or head end.
When the gateway has an issue (like it's offline), all the devices accessed via it also go down with the same issue.
This can cause a lot of noise in our alerts log.

It would be great if we could identify/model this case in some way so we can see the underlying cause of a device issue,
and/or see those underlying issues as top level alerts.

For example:

- I want to see when a gateway/controller/head end goes down and what devices are affected
- I want to see when a device is offline because the gateway is offline and distinguish this from an issue with the
  device itself

### Global Event Log

I'm wondering if we can model alerts as a sub-view of a central event log system.
If the log contained all events (ACS, device, system, UI, access logs, etc) and was quick at filtering and sorting, we
would be able to treat alerts like we treat devices on dashboards.
If you can put together a query for the alerts you want then you can have an alerts page, or an access control log,
or a system log, or a user journey log, etc.

My worry here is that the distinction between 'something just happened' and 'something is wrong' is too great for it to
be worth combining these concepts together.
You can always convert between the two, 'something is wrong' can be converted to an event via the leading edge: 'we just
noticed something is wrong'

Would we want to put our text logs into this system too? How would we then be different from Elastic or something like
that?
Would we end up with a _programming by configuration_ system, something we've tried to avoid as it just pushes the
data modelling down to the app developer. One of the reasons we use gRPC over MQTT.

### Reliability Trait

Should reliability be a different concept from status? Should it be a trait?
BACnet has a good separation for 'trouble getting the data' and 'the data says something is wrong'.
It calls the first 'reliability' and the second 'state'.

Reliability concerns all interactions with a device, it covers how accurate the telemetry we have _and_ how likely
control will be to succeed.
Maybe we need a new trait to cover reliability and we should refocus the status trait on issues reported by the device,
or discovered by the telemetry of the device.
If we did this we might be able to pull reliability into a common middleware based on RPC responses and error codes.

I think in Smart Core we should be able to model network failures, or any failure in the comms chain as reliability.
This would include network timeouts, COM issues, problems talking to slices or serial ports as well as any errors
returned by devices that also check for those things, aka errors further down the chain.

In general though people don't care as much about reliability, or more specifically the details of why the
_reliability score_ is low, than they do about the state of the system.

Splitting out reliability would allow us to annotate interactions with devices with a little badge saying "data may not
be accurate" but would also potentially make other use cases harder, like seeing the full status of a device.

### Triaging and understanding for alerts

TODO: Fill this out

## Design

### Status vs Alerts

I think it is useful to maintain the distinction between live state and historical logs of state.
However, I think we might need to define the relationship between the two more clearly, specifically the conversion
between them.
To satisfy requirements like "re-test this alert to check it's still active" the alert system can't just be a simple
database populated by autos reading the status trait.

Similarly, systems like BACnet have acknowledgement as a first class concept in their alert system.
It would be nice to be able to use the SC alert system to feed back these acks into the underlying system.
