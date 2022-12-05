# Alerts System

The alerts system manages a collection of alerts, these are records that indicate a change to the controller that can be
acknowledged by an operator.

**Example alerts:**

- Toilet paper is low in cubical 2
- Unable to communicate with BMS subsystem
- Power has been restored to floor 6

The alert system as implemented is backed by a database holding the alerts themselves, one row per alert. The alert
system will setup the relevant tables in the database when started.

## Alert metadata

Alerts can contain some metadata, typically to describe where the alert originated. This metadata is fairly unstructured
including the fields `floor`, `zone`, and `source` which are intended to aid with filtering in user interfaces. For
example to show only the alerts on floor 3, or that involve Meeting Room 2.
