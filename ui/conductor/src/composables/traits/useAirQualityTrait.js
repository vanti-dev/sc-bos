import Vue from 'vue';
import {computed, onMounted, onUnmounted, reactive, ref, watch, watchEffect} from 'vue';
import useDevices from '@/composables/useDevices';
import useTimePeriod from '@/routes/ops/components/useTimePeriod';
import {useNow} from '@/components/now';
import {listAirQualitySensorHistory} from '@/api/sc/traits/air-quality-sensor';
import {timestampToDate} from '@/api/convpb';

/**
 *
 * @param {{
 *   span: import('vue').GetterOrRef<number>,
 *   timeFrame: import('vue').GetterOrRef<number>,
 *   pollDelay: import('vue').GetterOrRef<number>
 * }} props
 * @return {{
 *   airQualitySensorHistoryValues: {},
 *   deviceOptions: import('vue').ComputedRef<{label: string, value: string}[]>,
 *   downloadCSV: (function(): void),
 *   isMounted: import('vue').Ref<boolean>,
 *   mappedDeviceNames: import('vue').ComputedRef<string[]>,
 *   periodEnd: import('vue').ComputedRef<Date>,
 *   periodStart: import('vue').ComputedRef<Date>,
 *   queryEnd: import('vue').ComputedRef<Date>,
 *   queryStart: import('vue').ComputedRef<Date>,
 *   setUpRequest: (function(string, number, string): void)
 * }}
 */
export default function(props) {
  // Computed properties for handling the 'search' for devices
  const {devicesData} = useDevices(props);
  // Mapping the device names to an array
  const mappedDeviceNames = computed(() => devicesData.value.map(device => device.name));
  const deviceOptions = computed(() => devicesData.value.map(device => {
    return {
      label: device.name,
      value: device.name
    };
  }
  ));

  // Time period and now
  const {now} = useNow(() => props.span);
  const {periodStart, periodEnd} = useTimePeriod(now, () => props.timeFrame, () => props.span);
  const queryStart = computed(() => new Date(periodStart.value.getTime() - props.span * 2));
  const queryEnd = computed(() => periodEnd.value);

  // Flag to indicate if the component is still mounted
  const isMounted = ref(false);

  // Data and polling state
  const airQualitySensorHistoryValues = reactive({});

  // Initialize data and polling state for each device
  const initializeAirQualityData = (deviceName) => {
    Vue.set(airQualitySensorHistoryValues, deviceName, {
      data: [],
      fetching: false,
      lastFetchTime: 0,
      pollHandler: 0,
      records: [],
      request: {}
    });
  };

  // Setting up the request object
  const setUpRequest = (deviceName, pageSize, pageToken) => {
    const request = {
      name: deviceName,
      pageSize: pageSize || 1000,
      pageToken: pageToken || '',
      timePeriod: {
        startTime: periodStart.value,
        endTime: periodEnd.value
      }
    };
    // Ensuring reactivity when setting up the request
    Vue.set(airQualitySensorHistoryValues[deviceName], 'request', request);
  };

  // Watcher for changes in devicesData to initialize data and polling state
  watch(devicesData, () => {
    if (!isMounted.value) return;

    mappedDeviceNames.value.forEach(name => {
      initializeAirQualityData(name);
      setUpRequest(name);
    });
  }, {immediate: true});


  // Function to fetch data
  // This function will fetch all pages of data and store it in the records array
  // It will also set the fetching flag to true while fetching
  const fetchData = async (device) => {
    const currentTime = Date.now();
    const deviceData = airQualitySensorHistoryValues[device];

    if (deviceData.fetching) return;

    // Check if enough time has passed since the last fetch
    const lastFetchTime = deviceData.lastFetchTime;
    if (lastFetchTime && (currentTime - lastFetchTime < props.pollDelay)) {
      console.warn('Not enough time has passed since last fetch');
      return;
    }

    deviceData.fetching = true;

    try {
      const listAction = listAllPages(deviceData.request, deviceData.records);
      deviceData.records = await listAction;
      deviceData.lastFetchTime = currentTime;
    } catch (e) {
      console.error('Error fetching air quality data:', e);
    } finally {
      deviceData.fetching = false;
    }

    // Initialize polling for each device
    Vue.set(deviceData, 'pollHandler', setInterval(
        () => fetchData(device), props.pollDelay)
    );
  };

  // If the airQualitySensorHistoryValues initialized and the request object is set up,
  // loop through each device and fetch data
  watchEffect(async () => {
    if (!isMounted.value && !Object.keys(airQualitySensorHistoryValues).length) return;

    mappedDeviceNames.value.forEach(device => {
      if (airQualitySensorHistoryValues[device].request && !airQualitySensorHistoryValues[device].lastFetchTime) {
        fetchData(device);
      }
    });
  });


  // Watcher for changes in records to process records for graph
  mappedDeviceNames.value.forEach(device => {
    watch(() => airQualitySensorHistoryValues[device].records, (records) => {
      if (!records) return;

      // Process records
      Vue.set(airQualitySensorHistoryValues[device], 'data', processRecordsForGraph(records));
    }, {immediate: true, deep: true});
  });


  // ---- Lifecycle hooks ---- //
  onMounted(async () => {
    isMounted.value = true;
  });

  onUnmounted(() => {
    // Clear pollHandler
    mappedDeviceNames.value.forEach(device => {
      clearInterval(airQualitySensorHistoryValues[device].pollHandler);
      Vue.delete(airQualitySensorHistoryValues, device);
    });

    isMounted.value = false;
  });


  // ---- Helper functions ---- //
  // Function to fetch all pages of data
  /**
   *
   * @param {ListAirQualityHistoryRequest.AsObject} request
   * @param {AirQualityRecord.AsObject[]} existingRecords
   * @return {Promise<AirQualityRecord.AsObject[]>}
   */
  async function listAllPages(request, existingRecords = []) {
    const allRecords = [...existingRecords];
    let pageToken = '';

    do {
      request.pageToken = pageToken;

      try {
        const page = await listAirQualitySensorHistory(request, {});
        allRecords.push(...page.airQualityRecordsList);
        pageToken = page.nextPageToken;

        // Check if component is still mounted before continuing
        if (!page.nextPageToken) break;
      } catch (e) {
        console.error('Error in listAllPages:', e);
        break; // Stop fetching on error
      }
    } while (pageToken);

    const removeDuplicates = (records) => {
      const seen = new Set();
      return records.filter(record => {
        const duplicate = seen.has(record.recordTime);
        seen.add(record.recordTime);
        return !duplicate;
      });
    };

    // Filter records to include only those within the specified time frame
    const filteredRecords = filterRecordsByTimeFrame(allRecords);

    // Remove duplicates and set data for graph
    const uniqueRecords = removeDuplicates(filteredRecords);

    Vue.set(airQualitySensorHistoryValues[request.name], 'records', uniqueRecords);
    Vue.set(airQualitySensorHistoryValues[request.name], 'data', processRecordsForGraph(uniqueRecords));

    return uniqueRecords;
  }

  // ---- Graph data ---- //
  const processRecordsForGraph = (records) => {
    return records.map(record => {
      return {
        y: {...record.airQuality},
        x: timestampToDate(record.recordTime)
      };
    });
  };

  // ---- CSV download ---- //
  const processRecordsForCSV = (records, deviceName) => {
    if (!records || !records.length) return '';

    const flattenedRecords = records.map(record => {
      const recordDate = timestampToDate(record.recordTime);

      return {
        deviceName,
        ...record.airQuality,
        recordTime: `${recordDate.toLocaleDateString()} ${recordDate.toLocaleTimeString()}`
      };
    });

    const header = Object.keys(flattenedRecords[0]);
    const rows = flattenedRecords.map(record =>
      header.map(fieldName => record[fieldName]).join(',')
    );

    return [header.join(',')].concat(rows).join('\n');
  };

  const downloadCSV = () => {
    let csvString = '';
    let headerIncluded = false;

    mappedDeviceNames.value.forEach(device => {
      const deviceRecords = airQualitySensorHistoryValues[device].records;
      const deviceCSV = processRecordsForCSV(deviceRecords, device);

      if (deviceCSV) {
        if (!headerIncluded) {
          csvString += deviceCSV;
          headerIncluded = true;
        } else {
          // Skip the header for subsequent devices
          const dataWithoutHeader = deviceCSV.substring(deviceCSV.indexOf('\n') + 1);
          csvString += '\n' + dataWithoutHeader;
        }
      }
    });

    const blob = new Blob([csvString], {type: 'text/csv;charset=utf-8;'});
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);

    link.href = url;
    link.setAttribute('download', 'air_quality_data.csv');
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  // ---- Cleanup ---- //

  // Function to remove records outside the time frame
  const filterRecordsByTimeFrame = (records) => {
    const timeFrameStart = new Date(now.value.getTime() - props.timeFrame);
    return records.filter(record => {
      const recordTime = timestampToDate(record.recordTime);
      return recordTime >= timeFrameStart && recordTime <= now.value;
    });
  };

  // ---- Return ---- //
  return {
    airQualitySensorHistoryValues,
    deviceOptions,
    downloadCSV,
    isMounted,
    mappedDeviceNames,
    periodEnd,
    periodStart,
    queryEnd,
    queryStart,
    setUpRequest
  };
}
