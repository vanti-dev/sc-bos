# Automation

## Features we want from automation stack

1. Configure via config files (or strings or bytes)
2. Update config for running automations without tearing down and building up again - allow the automation to choose
3. See the status of an automation - num runs, up time, switches toggled, etc.
4. Be able to stop/start automations via API calls
5. Add / remove automations at runtime
6. Support hand written and ITTT style automations
7. Support time based automations: i.e. turn off at 10, working hours mode, etc
8. Can contribute device traits to the nodes api

## Automation examples

- If a floor has been unoccupied for 30 minutes turn off all HVAC and lighting on the floor
- When the sensor level reaches 80% send an alert to Alice
- Adjust the ramp time for HVAC based on outdoor temperature/conditions
- Automatically return points to defaults after a period of inactivity after user setting
- Screens in common areas to be connected to timer to automatically turned off outside working hours
