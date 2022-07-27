Authentication and Authorization Use Cases
==========================================

## Types of Smart Core users and their requirements

### Landlord Employees
  - They want to access the Smart Core Web Apps from their workstations in the building
  - May need access to control, manage and monitor both common areas and tenant areas
  - They have accounts in the on-prem Active Directory server
  - Individual employees have functional areas they are responsible for and need to access (roles)
  - The system may have access to some (limited) personal information (think name, email, job title etc.)
    which needs to be protected

### Smart Core node-to-node auth
  - Landlord-owned devices need to communicate to each other
  - Control of area controllers from centralised location
  - Sending data from area controllers to be aggregated
  - Must be able to start up without human intervention

### Tenant Services
  - Servers operated by tenants
  - Not under our control how these are developed, deployed etc.
    - Auth scheme for tenants has to be fairly straightforward
  - They want access to monitor and control their tenant area using custom automations
  - Untrusted

### Tenant-area touch panels
  - Touch panels must be permitted to control the part of the tenant area that they are physically located in
  - Touch panels need to be commissioned. This should be feasible to do on a touch screen.
  - Touch panels are located in the tenant areas on the tenant areas, and should therefore be considered untrusted

### Tenant users 
  - Tenant users may want to access shared building apps, to monitor and control their tenant area
  - Authentication of individual tenant users is out of scope, but we need a way to permit access from the apps they use
  - Untrusted users using untrusted devices

### External Contractor
  - Not from the landlord or tenant
  - Visitor to the building
  - Wants to install, commission, update, maintain and diagnose smart building systems
  - May require quite high privilege levels
  - Only need access for a limited time

### Developers
  - Non-production environment
    - Developers assisting on-site are represented as External Contractor (above)
  - Running services locally
  - As few requirements/dependencies as possible
  - Should be easy and fast to set up
  - Not handling real data

## Open Questions

  1. 