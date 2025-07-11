# Health

For the purposes of this document, "health" is a measure of a device or system's operational state, indicating whether
it is functioning correctly or not.

## What we have now

There are four health concepts in Smart Core right now:

1. [Status trait](#status-trait)
2. [Alerts API](#alerts-api)
3. [SC API Health](#sc-api-health)
4. [Cohort Health](#cohort-health)
5. [Emergency Light Reports](#emergency-light-reports)

### Status trait

The status trait is a simple way for devices to report their current operational state.
It can be used to indicate whether a device is functioning correctly, has encountered an error, or is offline.
It is modelled around a `Level` enum with the following values:

```protobuf
enum Level {
  // The device is working as intended.
  // The device will respond to commands and is working within normal parameters.
  NOMINAL = 1;
  // The device is functioning but if left unattended may fail in the future.
  // For example the battery level may be low, or it is close to running out of paper.
  NOTICE = 2;
  // Some functions of the device are unavailable or in fault, but not all.
  // The intended function of the device is still available but at a reduced efficiency or capacity.
  // For example if only 3/4 lights in a corridor are working, you can still see but there is a fault.
  REDUCED_FUNCTION = 3;
  // The device is not performing its desired function.
  // Some of the device features may still be working but enough is faulty, or critical pieces are faulty such that
  // the reason for the device to exist is not being fulfilled.
  // For example an access control unit may have a working card reader but if the magnetic lock is broken then either
  // everyone (or nobody) can open the door negating the function of the device.
  NON_FUNCTIONAL = 4;
  // No communication with the device was possible
  OFFLINE = 127;
}
```

The trait supports the concept of _problems_, which is a list of additional levels for sub-aspects of the device.
For example, one such problem might be a request timeout, or a missing field, or a value outside of range.

Devices summarise all the problems they have recorded in a top level status for ease of consumption.

A zone is able to combine the problems from all devices within it to provide a zone level status using the same
problems-to-summary logic the devices use.

The status of a device is shown on the devices page of the Ops UI.
This is the live status only and doesn't include any alert acknowledgement or manual alerts that might be associated
with the device.

The status trait is used by automations to produce alert digests, sent via email, to email groups for multiple sites.
We have no reliable way to filter these digests to show only alerts that need attention.

### Alerts API

The alerts API provides access to historical issues, combining them with metadata for search, and a simple workflow
centred around acknowledging and resolving issues.
We have automations that create, update, and resolve alerts based on the status trait of devices.
A user has the ability to manually add records to the alerts list, we call these "manual alerts".

The alerts API allows for querying aggregate information, getting total alert counts for some filters, which powers
app features like alert count badges on nav items or toolbars.

The inclusion of device status in the alerts API requires automations to link the two systems together,
this is not an out-of-the-box feature of the alerts API.

Manual alerts associated with a device do not appear in the devices page.

We have automations that email a digest of alerts, from the last month, to email groups for multiple sites.

### SC API Health

Our first attempt at modelling health was in SC API, before SC BOS was a thing.
This info api split health into connection and communication health, each effectively an enum with values like
`CONNECTED` and `COMM_SUCCESS`.

There are no implementations of this API in any Smart Core (BOS or not) deployments.

### Cohort Health

There is no single API for this, but it is displayed prominently in the Ops UI on the app toolbar.
The health of the cohort is calculated by querying many APIs including the HubApi and EnrollmentApi.
Aside from connection errors or `UNIMPLEMENTED` responses, the health of the cohort also uses the
HubApi.TestHubNode to check the status of the cohort connections.

### Emergency Light Reports

The status of tests against emergency lights is accessed via the lighting test API.
The API allows for triggering both functional and duration tests, and provides a list of the results of the last tests.
A CSV of results can also be downloaded via the API.

The test results are _only_ available via this API and _only_ displayed on the dedicated emergency light reports page
of the Ops UI.
Viewing an emergency light in the devices page or on dashboards does not show the test results.

## Our wishlist

Including our existing functionality of:
seeing a devices current health,
seeing historical issues for a device,
seeing all issues in one place with filtering/search, and
aggregate health for zones.

1. We would like our health system to be more joined up.
   If a device has an issue, either detected or reported manually, we should be able to see it in all places.
   If an alert has been acknowledged, we should be able to see that in the devices page.

2. We should be able to distinguish between real issues and issues querying for issues.
   A separation between comm errors and device issues would be useful.
   An offline sensor on a critical device shouldn't (necessarily) be treated as a critical issue,
   more a reduction in the confidence or reliability of our monitoring of the device.

3. We need to be able to separate issues by criticality or impact.
   Many times we've been asked to act differently for 'life safety' issues, over other issues.
   Sometimes the device itself is a 'life safety' device, sometimes the monitored value is critical to life safety.

4. We'd like to be able to show a list of _devices that have issues_ to users.
   This should compliment the existing _list of issues_ view we have.
   This should be filterable by the usual floor, zone, etc criteria.

5. We'd like to be able to highlight when device points/values are outside expected ranges.
   Reporting what those ranges are and what the value is would be great.
   Reporting the value in a way it can be visualised in a dashboard would be even better.

6. We want to be able to show subsystem health, ideally in a way that matches the graphics/designs we put together for
   shows, the building-by-floors graphic.

7. It would be nice to see if a device issue is caused by a parent device.
   Our access to devices is often dependent on our access to head ends or gateways,
   distinguishing between a device being offline because the head end is offline, 
   or because the device itself is offline would be useful.
   Being able to collapse these alerts to avoid noise would be a good UX.

8. We want to be able to see when a system isn't available yet.
   This might be the same as offline?

9. We'd like to be able to have a "Check Now" button that forces a health check of a device.

10. We want to see how long a device has been in a bad or good state.

11. It would be nice to see when the last check was made, and when the next check will be made.

12. Out-of-range checks should support both device native and external checks.
    BACnet device already perform out-of-range checks, but we want to be able to do this for other devices too,
    or when the BACnet device isn't commissioned to perform its own checks.

13. We want to be able to configure out-of-range checks via the Ops UI.
    We'd like to be able to create a page in the Ops UI that configures these checks for one or more devices at a time.
    I don't think the driver edit pages are the correct place for this.

14. It would be nice to be able to see flaky devices.
    Devices that are frequently going offline or reporting issues, but not currently faulty.

## Implications of our wishlist

Combining point 1, 4, 12, and 13 implies that we need a managed health system,
similar to how we have a managed metadata and announcement system.
Our current approach of having each driver implement its own status trait will be difficult to get working with
health-based queries for devices, and with third-party checks updating the health of the device.

Point 8 also plays into this somewhat, though not directly.
If we added a building model that knew about subsystems, then that could contribute health status before a driver has
been commissioned.

Points 5, 9, 12, and 13 imply that we need a concept similar to a health check as a first class concept.
That's to say, something that reads the state of a device and decides if it is healthy or not.
This is different from our current approach where this is implied by Status.Problem, but not really applicable to this
workflow.
The system that does this could be part of the driver, but I suspect it'd be cleaner/better if it was a separate system.

Points 2, 3, 4, 6, and 7 all imply that we need to adjust the model we use to represent health.
We must be able to unambiguously be able to distinguish between reported faults and failures to read faults.
We must be able to search and filter by these criteria both from a _list of issues_ and a _list of devices with issues_
point of view.

Point 14 implies that current health and historical health are more closely related.
Maybe our new model of health should include a _recently_ concept, that complements the current health?
Other health metrics are calculated over time, anything that is _average_ like latency or load.

Points 10, 11, and 14 imply that some kind of persistent storage should exist even for current health.
This assumes we want to maintain _faulty since T_ across restarts.
Point 8 might also be related, especially if it is expanded to include _used to exist_.

Point 7 might be solvable via other mechanisms attached to devices.
Maybe a better model of parent-child, or a more consistent use of our existing parent-child model would help here.
Logic like: device is faulty, does device have a parent, is parent faulty, if so then device is _likely_ faulty because
of the parents issues.
While this could work, it's likely that the driver is already aware of the parent issues and can update fields in the
child health report much more efficiently.

It could also be the case that parent-child isn't the only (or correct) relationship we want to model.
Maybe a _related alerts_ or _health group_ concept is more useful.
This could apply to things like DALI buses, or BACnet gateways, each managed by a single driver.

Alternatively the parent-child relationship could be modelled as part of the health system.
A health report could have a parent health report.
This would imply that health reports are associated with devices, not owned by devices:
it is no longer _device has health report_, but more _health report for device_.

