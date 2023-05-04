# BACnet Network Discovery Tool

This tool scans a local BACnet network for devices, outputting their network information to a file.

```shell
bacnet-whois --iface 10.1.103.50
```

Produces a file like this:

```csv
BACnet Device ID,IP:Port,Network,Address,Max APDU,Segmentation,Vendor
819452,10.211.55.3:49665,,,1464,0,260
819453,10.211.55.3:58356,101,56,1464,0,260
```

The `--iface` arg is optional, and will default to the external IP of the local machine. You may specify either a name
or IP to choose the interface the BACnet client binds to.

```shell
bacnet-whois --iface eth0        # Binds to interface 'eth0'
bacnet-whois --iface 10.1.103.50 # Binds to the interface associated with '10.1.104.50'
```

See `bacnet-whois --help` for more configuration arguments.
