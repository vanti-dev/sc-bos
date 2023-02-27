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
      ...this._config.hvac ?? [],
      ...this._config.occupancySensors ?? []
    ];
  }

  _deviceTypes = [
    {key: 'lights', type: 'light'},
    {key: 'hvac', type: 'hvac'},
    {key: 'occupancySensors', type: 'occupancy sensor'}
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
    console.log('zone:set devices', deviceList);
  }
}
