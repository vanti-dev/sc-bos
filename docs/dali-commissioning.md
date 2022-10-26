DALI Commissioning
==================

The PLC project `bsp-ew` integrates the functionality of
Beckhoff's [DALI PLC Commissioning Tool](https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc3_dali/5967742731.html&id=)
.
This serves a GUI via TwinCAT 3 HMI Web.

To access the DALI Commissioning GUI, navigate to `https://<controller-ip>/Tc3PlcHmiWeb/Port_851/Visu/webvisu.htm`
in a web browser.

Each PLC controls several DALI buses. You can only work with one bus at a time in this tool. Select the bus using
the *DALI Line* box in the upper-left. Refer to the following table for which bus number is which:

| Bus                   | *DALI Line* number |
|-----------------------|--------------------|
| Tenant 1 Bus 1        | 1                  |
| Tenant 1 Bus 2        | 2                  |
| Tenant 2 Bus 1        | 3                  |
| Tenant 2 Bus 2        | 4                  |
| Tenant 2 Bus 3        | 5                  |
| Tenant 3 Bus 1        | 6                  |
| Tenant 3 Bus 2        | 7                  |
| Tenant 4 Bus 1        | 8                  |
| Tenant 4 Bus 2        | 9                  |
| Life Safety Bus       | 10                 |
| Landlord Buses 1 to 6 | 11-16              |

## Addressing Control Gear (Luminaires inc. Emergency Lights)

1. Switch to the *Control Gears - Addressing* page using the dropdown at the top.
2. Configure the addressing process using the *Random Addressing* box on the left
    1. If setting up the bus from scratch, tick *Complete new installation* and set the start address to 0.
    2. If adding new control gear to an existing box, ensure *Complete new installation* is not ticked.
       Set the start address to be one higher than the largest address currently on the bus. If you're not sure what's
       on the bus, see the *Discovering Control Gear* section for a procedure to detect them.
3. Commence the addressing by pressing *Start*. The process can take a long time - when it is done, the
   spinner will disappear. You can now proceed to *Discovering Control Gear*

## Discovering Control Gear

1. Switch to the *Control Gears - Addressing* page.
2. Press *Scan Devices by Address*. Once complete, a list of addresses for all detected Control Gear is shown in
   the list above. Note down the lowest and highest address numbers.
3. Switch to the *Control Gears - Parameters* page.
4. Select *Short Address Range* and enter the lowest and highest address numbers (from step 2) in the adjacent boxes.
5. Press *Read*. Parameters will be read from all the Control Gear.
6. Scroll horizontally to find the *DeviceTypes* column. Device Type 1 indicates that control gear is an emergency
   light. This is a convenient time to record a list of the luminaires and emergency lights you are expecting to find.

## Locating Control Gear

1. Switch to *Control Gears - Commands*
2. Select *Short Address* and enter the address of the control gear you want to find
3. Send commands by pressing buttons in the *Control* Section
   1. For luminaires, you can press *Recall Max Level* and *Off* alternately to flash the luminaire.
   2. For emergency lights, you can press the *Identify* button. This will cause the integrated colour LED to
      flash alternating red-green for about ten seconds. After that, press the button again to keep it flashing.

NB. For testing the entire bus, you can select *Broadcast* and send commands to all luminaires that way.

## Addressing Control Devices (PIR sensors, switches)

1. Switch to the *Control Devices - Addressing* page using the dropdown at the top.
2. Configure the addressing process using the *Random Addressing* box on the left
    1. If setting up the bus from scratch, tick *Complete new installation* and set the start address to 0.
    2. If adding new control device to an existing bus, ensure *Complete new installation* is not ticked.
       Set the start address to be one higher than the largest address currently on the bus. If you're not sure what's
       on the bus, see the *Discovering Control Devices* section for a procedure to detect them.
3. Commence the addressing by pressing *Start*. The process can take a long time - when it is done, the
   spinner will disappear. You can now proceed to *Discovering Control Devices*

## Discovering Control Devices

1. Switch to the *Control Devices - Addressing* page.
2. Press *Scan Devices by Address*. Once complete, a list of addresses for all detected Control Devices is shown in
   the list above. Note down the lowest and highest address numbers.
3. Switch to the *Control Devices - Parameters* page.
4. Select *Short Address Range* and enter the lowest and highest address numbers (from step 2) in the adjacent boxes.
5. Press *Read*. Parameters will be read from all the Control Devices.
   light. This is a convenient time to record a list of the luminaires and emergency lights you are expecting to find.
6. If a control device has data in the *IT1 (PushButton)* column, then it's a DALI switch. If it has data in the
   *IT3 (Occupancy)* column then it's a PIR sensor.

## Locating Control Devices

1. Switch to *Control Devices - Device Type 0 (Generic instance)*
2. Next to *Short Address* enter the address of the control device you want to find
3. Press the *Identify* button to ask the control device to identify itself. For PIR sensors, they will blink
   the built-in green LED very rapidly for about 10 seconds. You can keep pressing the *Identify* button until you find
   it. Note that PIR sensors will blink the LED slowly when occupancy is detected - the *Identify* command causes it
   to blink faster than that.

## Troubleshooting
#### Luminaire responds to broadcast commands, but I can't find an address for it
The luminaire may not have been addressed. See the Addressing Control Gear section.
Take care not to select *Complete new installation* as that will change the addresses of luminaires that have
already been addressed.

#### Luminaire is always on, and doesn't respond to broadcast commands
The luminaire is likely not connected to the DALI bus properly. Check the wiring.

#### Luminaire is always off, and doesn't respond to broadcast commands
The luminaire is likely not receiving power. Check the wiring

#### PIR sensor appears not to have an address
Check if the PIR sensor is connected to the bus by walking under it. If you observe a slowly flashing LED, then
the PIR sensor is connected to the bus. It might not have been addressed - see the Addressing Control Devices
section.
If the PIR sensor does not flash the LED and you are certain it is wired into the DALI bus correctly, the PIR
sensor may be defective. Try swapping it out with another unit.