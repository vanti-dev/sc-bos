export class Zone {
  /**
   * @param {Service.AsObject} zoneService
   */
  constructor(zoneService) {
    if (!zoneService) throw new Error('zoneService must be provided');
    this._id = zoneService.id;
    this._config = JSON.parse(zoneService.configRaw);
  }

  get name() {
    return this._config.name || this._id;
  }

  get deviceIds() {
    console.log(this._config);
    return [
      ...this._config.lights ?? [],
      ...this._config.hvac ?? [],
      ...this._config.occupancySensors ?? []
    ];
  }
}
