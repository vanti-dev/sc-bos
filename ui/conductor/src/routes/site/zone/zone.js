import {configureService, ServiceNames as ServiceTypes} from '@/api/ui/services';

export class Zone {
  /**
   * @param {Service.AsObject} zoneService
   */
  constructor(zoneService) {
    if (!zoneService) throw new Error('zoneService must be provided');
    this._id = zoneService.id;
    if (zoneService.configRaw === '') {
      this._config = {};
    } else {
      this._config = JSON.parse(zoneService.configRaw);
    }
  }

  get name() {
    return this._config.name || this._id;
  }

  get deviceIds() {
    return [
      ...this._config.lights ?? [],
      ...this._config.thermostats ?? [],
      ...this._config.occupancySensors ?? []
    ];
  }

  _deviceTypes = [
    {key: 'lights', type: 'light', trait: 'smartcore.traits.Light'},
    {key: 'thermostats', type: 'thermostat', trait: 'smartcore.traits.AirTemperature'},
    {key: 'occupancySensors', type: 'occupancy sensor', trait: 'smartcore.traits.OccupancySensor'}
  ];

  get devices() {
    const d = [];
    this._deviceTypes.forEach(t => {
      if (this._config.hasOwnProperty(t.key)) {
        this._config[t.key].forEach(n => d.push({name: n, type: t.type}));
      }
    });
    return d;
  }

  set devices(deviceList) {
    this._newConfig = {
      name: this._config.name,
      type: this._config.type
    };
    // for each device in the zone
    deviceList.forEach(d => {
      // loop through the deviceTypes and add to the corresponding part of the zone config
      this._deviceTypes.forEach(t => {
        if (d.metadata.traitsList.map(t => t.name).includes(t.trait)) {
          if (!this._newConfig.hasOwnProperty(t.key)) {
            this._newConfig[t.key] = [];
          }
          this._newConfig[t.key].push(d.name);
        }
      });
    });
  }

  async save(saveTracker) {
    if (this._newConfig) {
      this._config = this._newConfig;
      const req = {
        name: ServiceTypes.Zones,
        id: this._id,
        configRaw: JSON.stringify(this._newConfig, null, 2)
      };
      await configureService(req, saveTracker);
    }
    this._newConfig = undefined;
  }
}
