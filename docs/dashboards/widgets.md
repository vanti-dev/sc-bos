# Dashboard Widgets

## Energy Widgets

### Energy Consumption Widget

![Energy Consumption Widget](../assets/dashboard/widgets/energy-consumption.png)

Shows the total energy consumed, as measured by meters, over a specified time period.
The time period is something like day, week, month.
Can also show other metered consumption, like water or gas.

### Energy History Widget

![Energy History Widget](../assets/dashboard/widgets/energy-history.png)

Shows the energy consumed over a time period divided into time buckets and stacked by source.
Can show consumption and generation, above and below the line respectively.

### Energy Power Compare Widget

![Energy Power Compare Widget](../assets/dashboard/widgets/energy-power-compare.png)

Compares the current power consumption of one or more sources, showing a total.

### Energy Demand History Widget

![Energy Demand History Widget](../assets/dashboard/widgets/energy-demand-history.png)

Shows comparative or stacked historical electrical demand for one or more sources over time.
Has similar controls to the [Energy History Widget](#energy-history-widget), but averages the kW usage over each time
period, rather than summing the kWh used.
Can be configured to show different power metrics such as apparent power (kVA), reactive power (kVAR), or current (A).

### Energy Power Widget

![Energy Power Widget](../assets/dashboard/widgets/energy-power-history.png)

An alternative to the [Energy History Widget](#energy-history-widget) that also shows a total summary.

## Environmental Widgets

### Environmental Air Quality Widget

![Environmental Air Quality Widget](../assets/dashboard/widgets/environmental-air-quality.png)

Shows the current air quality for a single device.

### Environmental Air Quality History Widget

![Environmental Air Quality History Widget](../assets/dashboard/widgets/environmental-air-quality-history.png)

Shows historical readings for a single air quality metric for a single device over time.

### Environmental Card

![Environmental Card](../assets/dashboard/widgets/environmental-card.png)

Shows temperature and humidity readings for one or two devices.

## General Widgets

### General Data and Time Widget

![General Date and Time Widget](../assets/dashboard/widgets/general-date-and-time.png)

Shows the current date and time.

### General Smart Core Status Widget

![General Smart Core Status Widget](../assets/dashboard/widgets/general-sc-status.png)

Shows the status of the Smart Core cohort, just like on the app toolbar.

### General OpenWeather Map Widget

![General OpenWeather Map Widget](../assets/dashboard/widgets/general-open-weather-map.png)

Shows the current weather for a location using the OpenWeather API.
Needs an API key to work.

## Graphic Widgets

The graphic widget is described by [a dedicated page](../feats/opsui-graphics.md).
This widget can show any SVG graphic and bind it to data sources and make it interactive.

<p>
<img src="../assets/dashboard/widgets/graphic-fcu.png" alt="FCU Graphic Widget" height="150">
<img src="../assets/dashboard/widgets/graphic-floor-plan.png" alt="Floorplan Graphic Widget" height="150">
</p>

## Notification Widgets

![Notification Table Widget](../assets/dashboard/widgets/notifications-table.png)

Shows notifications, like the Notifications page.
Supports filtering notifications using any query normally supported by the Notifications page.

## Occupancy Widgets

### Occupancy People Count Widget

![Occupancy People Count Widget](../assets/dashboard/widgets/occupancy-people-count.png)

Shows the occupancy os a single device, with an optional maximum occupancy limit.

### Occupancy People Count History Widget

![Occupancy People Count History Widget](../assets/dashboard/widgets/occupancy-people-count-history.png)

Show the people count for a single device over time.

### Occupancy History Widget

![Occupancy History Widget](../assets/dashboard/widgets/occupancy-history.png)

An older version of the [Occupancy People Count History Widget](#occupancy-people-count-history-widget)
that shows historical occupancy data and a summary.

## Security Widgets

### Security Events Widget

![Security Events Widget](../assets/dashboard/widgets/security-events.png)

Shows access control events for a single device, such as door access or security system events.