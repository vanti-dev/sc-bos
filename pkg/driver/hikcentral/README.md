# HikCentral - CCTV

This driver integrates with HikCentral, a CCTV system provided by Hikvision.
Specifically, it talks to the HikVision OpenAPI, which is a separate service from the HikCentral platform used to monitor the devices.  
HikVision OpenAPI provides a REST API that allows for integration with the HikCentral platform.

There is a small example config file in `pkg/driver/hikcentral/config/example.json`. 
The config generation itself is project specific but for all integrations, the HikCentral server will need to be configured and an APi user created
details of this can be found here https://www.hikvisioneurope.com/eu//portal/portal/Technical%20Materials/24%20How%20To/HikCentral%20Professional/HCP%20Platform%20OpenAPI%20Deployment%20%26%20Online%20Debug.pdf

Once the API user is created, this driver requires the APP key and APP secret which needs to be grabbed from the OpenAPI (Artemis) web interface.