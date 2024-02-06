# Available Smart-Core roles and their privileges

## Table of contents

- [Available Smart-Core roles and their privileges](#available-smart-core-roles-and-their-privileges)
    - [Table of contents](#table-of-contents)
    - [Table with all available roles and their privileges](#table-with-all-available-roles-and-their-privileges)
    - [Available pages and their functionalities](#available-pages-and-their-functionalities)
        - [User and Tenant](#user-and-tenant)
        - [Devices](#devices)
        - [Operations](#operations)
        - [Workflows and Automations](#workflows-and-automations)
        - [Site Configuration](#site-configuration)
        - [System](#system)
        - [Signage](#signage)
    - [Available pages depending on the role](#available-pages-depending-on-the-role)
    - [Role details](#role-details)
        - [Admin role](#admin-role)
        - [Commissioner role](#commissioner-role)
        - [Operator role](#operator-role)
        - [Signage role](#signage-role)
        - [Super-Admin role](#super-admin-role)
        - [Viewer role](#viewer-role)
    - [Notes](#notes)

> Last updated 07/09/2023

## Table with all available roles and their privileges

| Role/Privilege | Read | Write | Limited    |
|----------------|------|-------|------------|
| Admin          | True | True  | _False_    |
| Commissioner   | True | True  | **_True_** |
| Operator       | True | True  | **_True_** |
| Signage        | True | -     | **_True_** |
| Super-Admin    | True | True  | _False_    |
| Viewer         | True | -     | **_True_** |

## Available pages and their functionalities

### User and Tenant

> - Add / Edit / Delete users
> - Add / Edit / Delete tenants
> - Associate zones with tenant
> - Create / Revoke tenant API keys _(aka. tenant tokens)_

### Devices

> - List all devices (all and / or grouped by device type)
> - Filter devices by floor
> - View all and / or selected device's network status
> - View selected device's information such as general information _(location, name, zone etc.)_ and it's relevant
    trait(s) _(temperature, brightness, energyStorage etc.)_
> - Control device via it's trait such as _dimming brightness; increasing temperature etc._

### Operations

> - Building Overview:
    >
- View power related information
>   - Toggle between different value types for power
>   - View occupancy details
>   - View live information regarding the environment such as _average indoor temperature; external temperature etc._
> - Notifications:
    >
- View device related notifications (alert, error, informative, severe, warning types) (up to 10 notification per
  device)
>   - Filter notifications by Subsystem (device types such as ACU, HVAC, Lighting etc.), Floor and / or Zone
>   - Acknowledge notification
> - Emergency Lighting:
    >
- View emergency lighting status per device
>   - Download report as CSV
> - Security Overview:
    >
- Search and / or filter Access Control Units (ACU) by floors
>   - View the list of ACU's with live information such as severity type, actor information, status etc. (view 1.)
>   - View a floor plan for each floor with live ACU information (view 2.)

### Workflows and Automations

> - List all automations (all and / or grouped by automation type)
> - Start / Stop automation
> - Copy selected automation's configuration code

### Site Configuration

> Note[^FutureChange]
>
> - List building zones
> - Filter devices by floors
> - Edit zone (add or remove devices grouped into zone)

### System

> - List drivers; features; components
> - Start / Stop drivers and / or features
> - Copy selected driver's or feature's configuration code
> - Monitor components

### Signage

> - `Signage` page is a community page which represents a reduced number of page elements and device information - such
    as power consumption; occupancy data; environment data etc.

---

## Available pages depending on the role

| Role / Page                        | [User and Tenant](#user-and-tenant) | [Devices](#devices) | [Operations](#operations) | [Workflows and Automations](#workflows-and-automations) | [Site Configuration](#site-configuration) | [System](#system) | [Signage](#signage) |
|------------------------------------|-------------------------------------|---------------------|---------------------------|---------------------------------------------------------|-------------------------------------------|-------------------|---------------------|
| [Admin](#admin-role)               | Yes                                 | Yes                 | Yes                       | Yes                                                     | Yes                                       | Yes               | Yes                 |
| [Commissioner](#commissioner-role) | -                                   | Yes                 | Yes                       | Yes                                                     | Yes                                       | Yes               | Yes                 |
| [Operator](#operator-role)         | -                                   | Yes                 | Yes                       | Yes _(limited)_                                         | Yes _(limited)_                           | Yes _(limited)_   | Yes                 |
| [Signage](#signage-role)           | -                                   | -                   | -                         | -                                                       | -                                         | -                 | Yes                 |
| [Super-Admin](#super-admin-role)   | Yes                                 | Yes                 | Yes                       | Yes                                                     | Yes                                       | Yes               | Yes                 |
| [Viewer](#viewer-role)             | Yes _(limited)_                     | Yes _(limited)_     | Yes _(limited)_           | Yes _(limited)_                                         | Yes _(limited)_                           | Yes _(limited)_   | Yes                 |

---

## Role details

### Admin role

> A user with `admin` privilege can read, write and manage the whole Smart-Core system - highest level of access with
> full access to all of the routes and no blocked routes or actions.[^Options]

---

### Commissioner role

> A role with considerable access but restricted from accessing authentication-related routes.

---

### Operator role

> A role with a mixture of full and limited access to various routes, ideal for operation-related tasks - has limited
> system management capabilities.

---

### Signage role

> Primarily has access to signage related routes, with all other routes being blocked.

---

### Super-Admin role

> Vanti/Smart-Core Admin. Similar to the admin but tailored for supervisory roles, with extensive full access routes.

---

### Viewer role

> As the role says, a `Viewer` only able to view the Smart-Core system and it's options and won't be able to do any type
> of action, such as:
>
> - **can't edit or update** settings or values;
> - **can't run** tests;
> - **can't add or remove** users or tenants;
> - **can't control** device and it's functionality;
> - **can't start or stop** drivers and / or features;
>
> just to mention a few.
> However, a `Viewer` is able to:
>
> - **view** all available pages and their non-invasive functionalities;
> - **list** all available devices and their information;
> - **list** all available automations and their configuration code;
> - **view** all available drivers and their configuration code;
> - **view** all available features and their configuration code;
> - **view** all available components and their status;
> - **list** all available zones and their devices;
> - **read** all available notifications;
>
> and many others, just to mention a few.

---

## Notes

[^Options]: Available Smart-Core options; device groups; device traits and functionality may vary depending on the
Smart-Core system configuration.
[^FutureChange]: Temporary functionality. This is going to change in future update(s)
