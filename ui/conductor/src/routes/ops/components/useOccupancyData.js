import {timestampToDate} from '@/api/convpb';
import {listOccupancySensorHistory} from '@/api/sc/traits/occupancy';
import {Occupancy} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';

/**
 *
 * @param {*} props
 * @return {{
 *  chartSeries: ComputedRef<Array<{name: string, data: Array<{x: Date, y: number, occupancyState: null|string}>}>>,
 *  chartSegments: ComputedRef<number[]>
 * }}
 */
export default function(props) {
  const pollDelay = computed(() => props.span / 10);
  const now = ref(Date.now());
  const nowHandle = ref(0);

  //
  //
  // Structure to hold the series data
  const seriesMap = reactive({
    occupancy: {
      baseRequest: computed(() => baseRequest(props.name)),
      data: [],
      handle: 0,
      records: /** @type {OccupancyRecord.AsObject[]} */ []
    }
  });

  /**
   * Return an array of object with request details
   *
   * @param {string} name
   * @return {ListOccupancyHistoryRequest.AsObject|undefined}
   */
  const baseRequest = (name) => {
    if (!name) return undefined;

    const period = {
      startTime: new Date(now.value - 24 * 60 * 60 * 1000)
    };

    return {
      name,
      period,
      pageSize: 1000,
      pageToken: ''
    };
  };

  //
  // Fill the series with empty data
  const chartSeries = computed(() => {
    return Object.entries(seriesMap)
        .map(([seriesName, seriesData]) =>
            seriesData.data.length > 0 ? {name: 'Occupancy', data: seriesData.data} : null
        )
        .filter((obj) => obj !== null);
  });

  const occupancyStates = ['state_unspecified', 'occupied', 'unoccupied', 'idle'];

  /**
   * @param {OccupancyRecord.AsObject[]} records
   * @return {Array<{x: Date, y: number, occupancyState: string}>}
   */
  const readRecords = (records) => {
    const graphIntervals = [];
    const currentDate = new Date();

    // Adjusting current date to nearest half-hour mark
    currentDate.setMinutes(currentDate.getMinutes() - (currentDate.getMinutes() % 30), 0, 0);

    // Populate the bar chart with data
    // 24 hour divided into 30 min intervals
    for (let i = 0; i < 48; i++) {
      const dataPoint = {
        x: new Date(currentDate.getTime() - i * 30 * 60 * 1000),
        y: 0,
        occupancyState: null
      };

      graphIntervals.unshift(dataPoint); // Update the array of objects depending on the currentDate
    }

    // Split each hour into 30min intervals and group the records while finding the highest value
    for (const record of records) {
      const recordTime = new Date(timestampToDate(record.recordTime));
      const minute = recordTime.getMinutes();
      const intervalStart = new Date(recordTime); // Separating the hours into 30 min intervals
      intervalStart.setMinutes(minute < 30 ? 0 : 30, 0, 0); // Start of the interval

      // Looking for the existing interval record
      const existingInterval = graphIntervals.find((interval) => interval.x.getTime() === intervalStart.getTime());
      const lastInterval = graphIntervals[graphIntervals.length - 1];

      // Initialize lastUpdatedIndex to the index of the last interval
      let lastUpdatedIndex = -1;

      // Updating the interval record if a higher peopleCount comes in
      // and updating the occupancy state
      if (existingInterval) {
        // Update everything before the last record
        if (lastUpdatedIndex !== graphIntervals.length - 1) {
          // Updating the peopleCount value if a higher value comes in
          if (record.occupancy.peopleCount > existingInterval.y) {
            existingInterval.y = record.occupancy.peopleCount;
          }

          // Update the occupancy state value
          // if an interval has occupied state, store and jump to the next interval
          if (record.occupancy.state === 1) {
            existingInterval.occupancyState = occupancyStates[record.occupancy.state];
            continue;
          }

          // Updating the last updated index to avoid updating the same record again
          lastUpdatedIndex = graphIntervals.indexOf(existingInterval);
        } else if (lastUpdatedIndex === graphIntervals.length - 1) {
          // Updating the peopleCount value if a higher value comes in
          if (record.occupancy.peopleCount > lastInterval.y) {
            lastInterval.y = record.occupancy.peopleCount;
          }

          // Update the occupancy state value
          // if the occupancy state changes to occupied, store it
          if (record.occupancy.state === 1) {
            lastInterval.occupancyState = occupancyStates[record.occupancy.state];
          }
        }
      }
    }

    return graphIntervals;
  };

  /**
   *
   * @param {OccupancySensorHistoryRequest.AsObject} req
   * @param {string} type
   */
  async function pollReadings(req, type) {
    const all = [];
    try {
      while (true) {
        const page = await listOccupancySensorHistory(req, {});
        all.push(...page.occupancyRecordsList);
        req.pageToken = page.nextPageToken;
        if (!req.pageToken) break;
      }
    } catch (e) {
      console.error('error getting occupancy readings', e);
    }
    seriesMap[type].records = all;
    seriesMap[type].handle = setTimeout(() => pollReadings(req, type), pollDelay.value);
  }

  /**
   * Find the segments in the data
   * so we can color the chart background
   * depending on the occupancy state
   *
   * @param {number[]} series
   * @return {number[][]}
   */
  function findSegments(series) {
    const segments = [];
    let startIdx = 0;
    let currentNumber = series[0];
    for (let i = 1; i < series.length; i++) {
      if (series[i] !== currentNumber) {
        segments.push([startIdx, i - 1, series[startIdx]]);
        startIdx = i;
        currentNumber = series[i];
      }
    }
    segments.push([startIdx, series.length - 1, series[startIdx]]);

    return segments;
  }

  //
  // Find the segments in the data
  const chartSegments = computed(() => findSegments(chartSeries.value[0].data.map((data) => data.occupancyState)));

  /**
   * This function takes two arrays of records and combines them by matching the x values and updating the y values.
   *
   * @param {*} newRecords
   * @param {*} existingRecords
   * @return {Array} combinedRecords
   *
   * Update the existing data with the new data - highest occupancy count wins / 30 minute intervals
   * Map over the newRecords array and find the corresponding record in the existingRecords array.
   * If a match is found, update the y value with the larger of the two values.
   * If no match is found, return the new record.
   */
  const handleRecords = (newRecords, existingRecords) => {
    return newRecords.map((newRecord) => {
      const existingRecord = existingRecords.find((record) => record.x.getTime() === newRecord.x.getTime());
      if (existingRecord) {
        return {
          x: newRecord.x, // keep the existing x
          y: newRecord.y > existingRecord.y ? newRecord.y : existingRecord.y, // highest occupancy count wins
          occupancyState: newRecord.occupancyState // add the occupancy state
        };
      }

      return newRecord;
    });
  };

  //
  // Generate request for occupancy sensor history
  // and watch for changes then update the series data
  Object.entries(seriesMap).forEach(([name, series]) => {
    watch(
        () => series.baseRequest,
        (request) => {
          clearTimeout(series.handle);
          series.records = [];
          if (request) pollReadings(request, name);
        },
        {immediate: true, deep: true, flush: 'sync'}
    );

    watch(
        () => series.records,
        (records) => {
          const newRecords = readRecords(records);
          series.data = series.data.length ? handleRecords(newRecords, series.data) : newRecords;
        },
        {immediate: true, deep: true, flush: 'sync'}
    );
  });

  onMounted(() => {
    nowHandle.value = setInterval(() => (now.value = Date.now()), pollDelay.value);
  });
  onUnmounted(() => clearInterval(nowHandle.value));

  return {
    chartSeries,
    chartSegments
  };
}
