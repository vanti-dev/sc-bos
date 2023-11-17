import {camelToSentence} from '@/util/string';
import {AirQuality} from '@smart-core-os/sc-api-grpc-web/traits/air_quality_sensor_pb';
import Vue from 'vue';
import {computed, onMounted, onUnmounted, reactive, ref, watch, watchEffect} from 'vue';
import useDevices from '@/composables/useDevices';
import useTimePeriod from '@/routes/ops/components/useTimePeriod';
import {useNow} from '@/components/now';
import {listAirQualitySensorHistory} from '@/api/sc/traits/air-quality-sensor';
import {timestampToDate} from '@/api/convpb';
import {DAY} from '@/components/now';

/**
 *
 * @param {{
 *   name: import('vue').GetterOrRef<string>,
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
  // Mapping the device or zone names to an array
  const mappedDeviceNames = computed(() => {
    if (typeof props.name === 'string') {
      return [props.name];
    } else if (Array.isArray(props.name)) {
      return props.name;
    }
    return devicesData.value.map(device => device.name);
  });
  // const mappedDeviceNames = computed(() => [props.name]);
  const deviceOptions = computed(() => mappedDeviceNames.value.map(device => {
    return {
      label: device,
      value: device
    };
  }
  ));
  const airDevice = ref('');
  const previousAirDevice = ref(''); // Store the previous device name

  // Time period and now
  const {now} = useNow(() => props.span);
  const {periodStart: downloadPeriodStart, periodEnd: downloadPeriodEnd} = useTimePeriod(
      now, () => props.timeFrame * 30 * DAY, () => props.span
  );
  const {periodStart, periodEnd} = useTimePeriod(now, () => props.timeFrame, () => props.span);
  const queryStart = computed(() => new Date(periodStart.value.getTime() - props.span * 2));
  const queryEnd = computed(() => periodEnd.value);

  // Flag to indicate if the component is still mounted
  const isMounted = ref(false);
  const isFetching = ref(false);

  // Data and polling state
  const airQualitySensorHistoryValues = reactive({});

  // Initialize data and polling state for each device
  const initializeAirQualityData = (deviceName) => {
    if (!airQualitySensorHistoryValues[deviceName]) {
      Vue.set(airQualitySensorHistoryValues, deviceName, {
        data: [],
        fetching: false,
        lastSuccessfulFetchTime: 0,
        pollHandler: 0,
        records: [],
        request: {}
      });
    }
    setUpRequest(deviceName);
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

  // Watch airDevice.value to trigger re-initialization and data fetch
  watch(airDevice, (newDevice) => {
    if (newDevice && !airQualitySensorHistoryValues[newDevice]) {
      initializeAirQualityData(newDevice);
    }
  });

  // Function to fetch data
  // This function will fetch all pages of data and store it in the records array
  // It will also set the fetching flag to true while fetching
  const fetchData = async (device) => {
    const currentTime = Date.now();
    const deviceData = airQualitySensorHistoryValues[device];

    if (deviceData.fetching || isFetching.value) return;

    // Check if enough time has passed since the last fetch
    // Separate variable to track the time since the last successful fetch
    const timeSinceLastFetch = currentTime - (deviceData.lastSuccessfulFetchTime || 0);

    // Check if enough time has passed since the last successful fetch
    if (timeSinceLastFetch < props.pollDelay) {
      console.warn('Not enough time has passed since last fetch for device:', device);
      return;
    }

    deviceData.fetching = true;
    isFetching.value = true;

    try {
      const listAction = listAllPages(deviceData.request, deviceData.records);
      deviceData.records = await listAction;
      deviceData.lastSuccessfulFetchTime = Date.now(); // Update last successful fetch time
    } catch (e) {
      console.error('Error fetching air quality data:', e);
    } finally {
      deviceData.fetching = false;
      isFetching.value = false;
    }

    // Initialize polling for each device
    Vue.set(deviceData, 'pollHandler', setInterval(
        () => fetchData(device), props.pollDelay)
    );
  };


  // If the airQualitySensorHistoryValues initialized and the request object is set up,
  // loop through each device and fetch data
  watchEffect(async () => {
    // Check if there are any device names available
    if (mappedDeviceNames.value.length > 0) {
      // If airDevice is not set, set it to the first device
      if (!airDevice.value) {
        airDevice.value = mappedDeviceNames.value[0];
      }

      const currentDevice = airDevice.value;

      // Check if the current device is different from the previous one
      if (previousAirDevice.value !== currentDevice) {
        // If there was a previous device, delete its data
        if (previousAirDevice.value && airQualitySensorHistoryValues[previousAirDevice.value]) {
          Vue.delete(airQualitySensorHistoryValues, previousAirDevice.value);
        }

        // Initialize data for the current device
        if (!airQualitySensorHistoryValues[currentDevice]) {
          await initializeAirQualityData(currentDevice);
        }

        // Fetch data for the current device if it hasn't been fetched yet
        if (airQualitySensorHistoryValues[currentDevice]?.request &&
            !airQualitySensorHistoryValues[currentDevice].lastFetchTime) {
          await fetchData(currentDevice);
        }

        // Update the previousDevice value
        previousAirDevice.value = currentDevice;
      }
    }
  });


  // ---- Lifecycle hooks ---- //
  onMounted(async () => {
    isMounted.value = true;
  });

  onUnmounted(() => {
    // Clear pollHandler
    mappedDeviceNames.value.forEach(device => {
      const deviceData = airQualitySensorHistoryValues[device];
      if (deviceData.pollHandler) {
        clearInterval(deviceData.pollHandler);
      }
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
   * @param {string} type
   * @return {Promise<AirQualityRecord.AsObject[]>}
   */
  async function listAllPages(request, existingRecords = [], type = '') {
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

    if (type) return uniqueRecords;

    Vue.set(airQualitySensorHistoryValues[request.name], 'records', uniqueRecords);
    Vue.set(airQualitySensorHistoryValues[request.name], 'data', processRecordsForGraph(uniqueRecords));

    return uniqueRecords;
  }

  // Function to read the comfort value
  const readComfortValue = (level) => {
    if (AirQuality.Comfort.COMFORTABLE) {
      return 'Comfortable';
    } else if (AirQuality.Comfort.UNCOMFORTABLE) {
      return 'Uncomfortable';
    } else {
      return 'Unknown';
    }
  };

  // ---- Acronyms ---- //
  const acronyms = {
    carbonDioxideLevel: {label: 'CO₂', unit: 'ppm'},
    volatileOrganicCompounds: {label: 'VOC', unit: 'ppm'},
    airPressure: {label: 'Air Pressure', unit: 'hPa'},
    comfort: {label: 'Comfort', unit: ''},
    infectionRisk: {label: 'Infection Risk', unit: '%'},
    score: {label: 'Air Quality Score', unit: '%'}
  };

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
  const capitalize = (str) => str.charAt(0).toUpperCase() + str.slice(1);

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

    // Create a mapping for headers
    const headerMap = {};
    Object.keys(flattenedRecords[0]).forEach(key => {
      const label = acronyms[key]?.label || camelToSentence(key); // Fallback to camelToSentence if label is not defined
      const unit = acronyms[key]?.unit ? `(${acronyms[key].unit})` : '';
      headerMap[key] = `${capitalize(label)} ${unit}`;
    });

    // Generate the header row
    const header = Object.values(headerMap);

    // Generate the data rows
    const rows = flattenedRecords.map(record =>
      Object.keys(headerMap).map(key => record[key]).join(',')
    );

    return [header.join(',')].concat(rows).join('\n');
  };


  const downloadCSV = async () => {
    // Check if there's a selected device
    if (!airDevice.value) {
      console.warn('No device selected for downloading CSV');
      return;
    }

    const device = airDevice.value;
    const records = await fetchRecordsForCSV(device);

    // Check if there are records to process
    if (!records || records.length === 0) {
      console.warn('No records found for device:', device);
      return;
    }

    // Process records to CSV format
    const csvString = processRecordsForCSV(records, device);

    // Create a Blob from the CSV String
    const blob = new Blob([csvString], {type: 'text/csv;charset=utf-8;'});
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);

    // Trigger the download
    const date = new Date();
    const dateString = `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;
    link.href = url;
    link.setAttribute('download', `AirQualityData - ${device} - ${dateString}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  const fetchRecordsForCSV = async (deviceName) => {
    const request = {
      name: deviceName,
      pageSize: 1000,
      pageToken: '',
      timePeriod: {
        startTime: downloadPeriodStart.value,
        endTime: downloadPeriodEnd.value
      }
    };

    return await listAllPages(request, [], 'csv');
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
    acronyms,
    airQualitySensorHistoryValues,
    deviceOptions,
    airDevice,
    downloadCSV,
    isMounted,
    isFetching,
    mappedDeviceNames,
    periodEnd,
    periodStart,
    queryEnd,
    queryStart,
    setUpRequest,
    readComfortValue
  };
}
