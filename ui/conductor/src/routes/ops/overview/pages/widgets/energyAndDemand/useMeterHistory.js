import {timestampToDate} from '@/api/convpb';
import {listMeterReadingHistory} from '@/api/sc/traits/meter';
import {HOUR, MINUTE, useNow} from '@/components/now';
import {csvDownload} from '@/util/downloadCSV';
import {toValue} from '@/util/vue';
import debounce from 'debounce';
import {computed, ref, watch, watchEffect} from 'vue';

/**
 * @param {MaybeRefOrGetter<string>} name
 * @param {MaybeRefOrGetter<Date>} periodStart
 * @param {MaybeRefOrGetter<Date>} periodEnd
 * @param {MaybeRefOrGetter<number>} spanSize
 * @return {{
 *  records: import('vue').Ref<MeterReadingRecord.AsObject[]>,
 *  missingPeriods: import('vue').ComputedRef<Array<{start: Date, end: Date, type: 'set' | 'unshift' | 'push'}>>,
 *  allSeriesData: import('vue').ComputedRef<{x: Date, y: number, incomplete: boolean}[]>,
 *  seriesData: import('vue').ComputedRef<{x: Date, y: number}[]>,
 *  fetching: import('vue').Ref<boolean>,
 *  lastFetchTime: import('vue').Ref<Date|null>,
 *  fetchPeriod: import('vue').ComputedRef<number>,
 *  now: import('vue').Ref<Date>,
 *  shouldFetch: import('vue').ComputedRef<boolean>,
 *  firstRecordTime: import('vue').ComputedRef<Date|null>,
 *  lastRecordTime: import('vue').ComputedRef<Date|null>,
 *  exportData: (fileName) => Promise<void>
 * }}
 */
export default function(name, periodStart, periodEnd, spanSize) {
  // Contains all the raw (well AsObject) records we've fetched from the server
  const records = ref(/** @type {MeterReadingRecord.AsObject[]} */ []);
  // A boolean indicating whether the async fetch is in progress
  const fetching = ref(false);
  // A Date recording the last time we fetched data from the server.
  // Used to limit fetches to a reasonable rate.
  const lastFetchTime = ref(/** @type {Date|null} */null);
  watch(() => [toValue(periodStart), toValue(periodEnd)], () => {
    lastFetchTime.value = null; // reset if the period changes
  });
  // How often do we retry a fetch that didn't return all the data we're after.
  const fetchPeriod = computed(() => Math.max(toValue(spanSize) / 4, 10 * MINUTE));
  // Track the current time, so we know how long it's been since the last fetch.
  const {now} = useNow(() => fetchPeriod.value);
  // A boolean indicating whether an attempted fetch should proceed.
  const shouldFetch = computed(() => {
    if (fetching.value) return false;
    if (lastFetchTime.value === null) return true;
    return now.value.getTime() > lastFetchTime.value.getTime() + fetchPeriod.value;
  });

  // Watch for changes to the span size and name, and reset the records and the last fetch if any changes.
  watch([() => toValue(spanSize), () => toValue(name)], () => {
    console.log(toValue(name));
    records.value = []; // reset if the span size changes
    lastFetchTime.value = null; // reset if the span size or name changes
  }, {immediate: true});

  const firstRecordTime = computed(() => {
    if (records.value.length === 0) return null;
    return timestampToDate(records.value[0].recordTime);
  });
  const lastRecordTime = computed(() => {
    if (records.value.length === 0) return null;
    return timestampToDate(records.value[records.value.length - 1].recordTime);
  });

  const queryExtraLeadingSpans = 2; // how many extra spans at the start should we attempt to fetch
  const queryExtraLeadingTime = computed(() => Math.max(toValue(spanSize) * queryExtraLeadingSpans, 20 * MINUTE));
  const queryStart = computed(
      () => new Date(toValue(periodStart).getTime() - queryExtraLeadingTime.value));
  const queryEnd = computed(() => toValue(periodEnd));

  // Inspects both [periodStart,periodEnd] and records to calculate which slices of time we don't have.
  const missingPeriods = computed(() => {
    if (records.value.length === 0) {
      return [{start: toValue(queryStart), end: toValue(queryEnd), type: 'set'}];
    }
    const periodStartDate = toValue(queryStart);
    const periodEndDate = toValue(queryEnd);
    const firstRecordDate = firstRecordTime.value;
    const lastRecordDate = lastRecordTime.value;
    if (firstRecordDate > periodEndDate || lastRecordDate < periodStartDate) {
      // the new period does not overlap in any way with the existing records
      return [{start: periodStartDate, end: periodEndDate, type: 'set'}];
    }

    const spans = [];
    if (firstRecordDate > periodStartDate) {
      spans.push({start: periodStartDate, end: firstRecordDate, type: 'unshift'});
    }
    if (lastRecordDate < periodEndDate) {
      spans.push({start: lastRecordDate, end: periodEndDate, type: 'push'});
    }
    return spans;
  });

  // Debounce the data fetch to prevent multiple calls to the server within a certain time frame.
  const debouncedDataFetch = debounce(() => {
    // We use these later to trim the records array and remove items that aren't needed any more.
    // We grab them now so that they are tracked for reactivity before we await anything.
    const periodStartDate = toValue(queryStart);
    const periodEndDate = toValue(queryEnd);
    const periods = missingPeriods.value;

    fetching.value = true;
    listAllPeriods(toValue(name), toValue(periods), records)
        .then(() => {
          if (periods.length > 0) {
            // We've fetched new r
            records.value = deleteGarbageRecords(records.value, periodStartDate, periodEndDate);
            // If we need to sort the records, here's where we'd do it.
          }
        })
        .catch((err) => {
          console.error('Error fetching meter history', err);
        })
        .finally(() => fetching.value = false);
  }, 500, false);

  // This watch tracks the missingPeriods that we have to query the server for.
  // It will fetch the missing data from the server and store the results, along with any existing records,
  // in the records ref.
  watchEffect(() => {
    if (!shouldFetch.value || !toValue(name)) {
      return;
    }
    lastFetchTime.value = toValue(now);

    debouncedDataFetch();
  });

  const shouldSampleData = computed(() => toValue(spanSize) > 0);
  // Converts the records into series data for the chart.
  const allSeriesData = computed(() => {
    if (shouldSampleData.value) {
      return sampledSeriesData.value;
    } else {
      return unsampledSeriesData.value;
    }
  });
  // We work out the reading value before the start of a span, and the reading before the end of a span,
  // and diff them to calculate the y value for that span.
  // The x value will be the end date of the span.
  //
  // Incomplete spans are included but marked with "incomplete: true".
  const sampledSeriesData = computed(() => {
    const series =
        /** @type {Array<{x:Date,y:number,incomplete:boolean,len:number}>} */
        [];

    const size = toValue(spanSize); // how big are each span
    const lastSpanEnd = toValue(periodEnd).getTime(); // when do we stop

    let spanStart = toValue(periodStart).getTime(); // tracks the start of the current span
    let spanEnd = spanStart + size; // tracks the end of the current span

    let recordBeforeStart = /** @type {MeterReadingRecord.AsObject|null} */ null;
    let recordIndex = 0; // the record we're currently looking at

    // account for records that are before the start of this span
    const availableRecords = records.value;
    for (let i = 0; i < availableRecords.length; i++) {
      const record = availableRecords[i];
      const recordTime = timestampToDate(record.recordTime).getTime();
      if (recordTime >= spanStart) {
        recordIndex = i;
        break;
      } else {
        recordBeforeStart = record;
      }
    }

    for (; spanEnd <= lastSpanEnd; spanStart += size, spanEnd += size) {
      let recordBeforeEnd = recordBeforeStart;
      for (let i = recordIndex; i < availableRecords.length; i++) {
        const record = availableRecords[i];
        const recordTime = timestampToDate(record.recordTime).getTime();
        if (recordTime >= spanEnd) {
          recordIndex = i;
          break;
        } else {
          recordBeforeEnd = record;
        }
      }

      if (recordBeforeStart === null) {
        // we have no records before the start of this span, so we can't calculate a value
        series.push({x: new Date(spanEnd), y: 0, incomplete: true, len: 0});
      } else if (recordBeforeEnd === null) {
        // we have no records before the end of this span, so we can't calculate a value
        series.push({x: new Date(spanEnd), y: 0, incomplete: true, len: 0});
      } else {
        // we have a record before the start and before the end, so we can calculate a value
        const startValue = recordBeforeStart.meterReading.usage;
        const endValue = recordBeforeEnd.meterReading.usage;
        const spanValue = endValue - startValue;
        const startRecordTime = timestampToDate(recordBeforeStart.recordTime).getTime();
        const endRecordTime = timestampToDate(recordBeforeEnd.recordTime).getTime();
        const len = endRecordTime - startRecordTime;
        series.push({
          x: new Date(spanEnd),
          y: spanValue,
          incomplete: spanValue < 0 || len <= 0,
          len: len
        });
      }

      recordBeforeStart = recordBeforeEnd;
    }

    return series;
  });
  // Every pair of readings (within the period) that change the value of the reading are exported as a data point.
  const unsampledSeriesData = computed(() => {
    const series =
        /** @type {Array<{x:Date,y:number,incomplete:boolean,len:number}>} */
        [];
    const start = toValue(periodStart).getTime(); // tracks the start of the current span
    let lastRecord = null;
    const availableRecords = records.value;
    for (const record of availableRecords) {
      const recordDate = timestampToDate(record.recordTime);
      const recordTimestamp = recordDate.getTime();
      if (recordTimestamp < start) {
        continue;
      }
      if (lastRecord === null) {
        lastRecord = record;
        continue;
      }
      const diff = record.meterReading.usage - lastRecord.meterReading.usage;
      if (diff === 0) {
        continue;
      }
      const len = recordTimestamp - timestampToDate(lastRecord.recordTime).getTime();
      series.push({
        x: recordDate,
        y: diff,
        incomplete: false,
        len
      });
      lastRecord = record;
    }
    return series;
  });

  // the series data, but with incomplete spans set to 0
  const seriesData = computed(() => {
    const data = allSeriesData.value.map(val => {
      if (val.incomplete) {
        return {x: val.x.getTime(), y: null};
      } else {
        return {x: val.x.getTime(), y: val.y / val.len * HOUR};
      }
    });

    // add predicted values to fill in gaps in the series
    let lastItemWithValue = null;
    let lastIndexWithValue = -1;
    for (let i = 0; i < data.length; i++) {
      const item = data[i];
      if (item.y !== null) {
        if (lastItemWithValue !== null) {
          // fill in items from lastIndexWithValue to i using an average progression between lastIndexWithValue and i
          const diff = item.y - lastItemWithValue.y;
          const inc = diff / (i - lastIndexWithValue);
          for (let j = lastIndexWithValue + 1; j < i; j++) {
            data[j].y = lastItemWithValue.y + inc * (j - lastIndexWithValue);
            data[j].predicted = true;
          }
        }
        lastItemWithValue = item;
        lastIndexWithValue = i;
      }
    }

    return data;
  });

  // --------- Export data as CSV --------- //

  // The acronyms object contains the labels and units for the CSV file.
  const acronyms = {
    usage: {label: 'Meter Reading', unit: 'kWh'}
  };

  /**
   * Processes the records into spans of time and returns only the last record of each span.
   *
   * @template T
   * @param {Array<T>} records
   * @param {number} span
   * @return {Array<T>}
   */
  const processRecordsForCSV = (records, span) => {
    const spanRecords = []; // The records to be included in the output
    let currentSpanIndex = 0; // The index of the current span
    let spanStartTime = timestampToDate(records[0].recordTime).getTime(); // The start time of the current span

    // Iterate through the records and group them into spans
    records.forEach((record, index) => {
      const recordTime = timestampToDate(record.recordTime).getTime(); // The time of the current record
      // Check if the current record's time exceeds the current span's start time by the span duration
      if (recordTime - spanStartTime >= span) {
        // If the current record's time exceeds the current span's start time by the span duration,
        // push the last record of the previous span to spanRecords
        spanRecords.push(records[index - 1]);
        // Update the span start time to the current record's time for the next span
        spanStartTime = recordTime;
        currentSpanIndex = index; // Update the current span index
      }
    });

    // Add the last record if the last span has at least one record
    if (currentSpanIndex < records.length) {
      spanRecords.push(records[records.length - 1]);
    }

    return spanRecords;
  };

  /**
   * Flattens the records into a format that can be exported as a CSV file.
   * The output format is an array of objects, each with the following properties:
   * - deviceName: The name of the device
   * - usage: The usage value of the last record in the span
   * - recordTime: The time of the last record in the span
   *
   * @param {MeterReadingRecord.AsObject[]} records
   * @param {string} recordValue
   * @param {number} span
   * @param {string} deviceName
   * @return {Array<{deviceName: string, usage: number, recordTime: string}>}
   */
  const flattenedRecords = (records, recordValue, span, deviceName) => {
    // Process the records into spans
    const spanRecords = processRecordsForCSV(records, span);
    const findDeepestValue = (record, path) => {
      let obj = record;
      for (const part of path.split('.')) {
        obj = obj[part];
      }
      return obj;
    };

    // Convert the spanRecords to the desired output format
    return spanRecords.map(record => {
      const recordDate = timestampToDate(record.recordTime);
      return {
        deviceName,
        usage: findDeepestValue(record, recordValue),
        recordTime: `${recordDate.toLocaleDateString()} ${recordDate.toLocaleTimeString()}`
      };
    });
  };


  /**
   * Exports the data as a CSV file.
   * We process the existing records and convert them into a CSV format.
   *
   * @param {string} fileName
   * @return {Promise<void>}
   */
  const exportData = async (fileName) => {
    csvDownload({
      acronyms,
      docType: fileName,
      flattenRecords: (records) => {
        return flattenedRecords(records, 'meterReading.usage', toValue(spanSize), toValue(name));
      },
      records: () => records.value,
      deviceName: toValue(name)
    });
  };

  return {
    // the important data
    seriesData,

    // debug fields
    records,
    firstRecordTime,
    lastRecordTime,
    missingPeriods,
    allSeriesData,

    fetching,
    lastFetchTime,
    fetchPeriod,
    now,
    shouldFetch,
    exportData
  };
};

/**
 * @param {string} name
 * @param {Array<{start:Date, end:Date, type:string}>} periods
 * @param {import('vue').Ref<MeterReadingRecord.AsObject[]>} records
 * @return {Promise<void>}
 */
async function listAllPeriods(name, periods, records) {
  for (const {start, end, type} of periods) {
    const baseRequest = /** @type {ListMeterReadingHistoryRequest.AsObject} */{
      name: name,
      period: {
        startTime: start,
        endTime: new Date(end.getTime() - 1)
      }
    };

    const newRecords = await listAllPages(baseRequest);
    switch (type) {
      case 'set':
        records.value = newRecords;
        break;
      case 'unshift':
        // remove duplicates
        let i = newRecords.length - 1;
        const firstRecordTime = records.value[0].recordTime;
        for (; i >= 0; i--) {
          const nrTS = newRecords[i].recordTime;
          if (firstRecordTime.seconds !== nrTS.seconds || firstRecordTime.nanos !== nrTS.nanos) {
            break;
          }
        }
        if (i === -1) {
          // all records are duplicates
          break;
        } else {
          newRecords.splice(0, newRecords.length - i - 1);
        }
        if (newRecords.length === 0) {
          break;
        }
        records.value.unshift(...newRecords);
        break;
      case 'push':
        if (newRecords.length === 0) {
          break;
        }
        records.value.push(...newRecords);
        break;
    }
  }
}

/**
 * Executes the given baseRequest, collecting all subsequent pages, and returning them as a single array.
 *
 * @param {ListMeterReadingHistoryRequest.AsObject} baseRequest
 * @return {Promise<MeterReadingRecord.AsObject[]>}
 */
async function listAllPages(baseRequest) {
  baseRequest.pageSize = 1000;

  const all = /** @type {MeterReadingRecord.AsObject[]} */[];
  let pageToken = '';
  do {
    baseRequest.pageToken = pageToken;
    const page = await listMeterReadingHistory(baseRequest, {});
    all.push(...page.meterReadingRecordsList);
    pageToken = page.nextPageToken;
  } while (pageToken);
  return all;
}

/**
 * Returns an array with any record outside [periodStart,periodEnd] removed.
 *
 * @param {MeterReadingRecord.AsObject[]} records
 * @param {Date} periodStart
 * @param {Date} periodEnd
 * @return {MeterReadingRecord.AsObject[]}
 */
function deleteGarbageRecords(records, periodStart, periodEnd) {
  if (records.length === 0) {
    return records;
  }
  let removeFromStart = 0;
  for (; removeFromStart < records.length; removeFromStart++) {
    const r = records[removeFromStart];
    if (timestampToDate(r.recordTime) >= periodStart) {
      break; // have reached the first record that is within the period (well is after the start of it)
    }
  }

  // quick bail if we're going to be removing all records
  if (removeFromStart === records.length) {
    return [];
  }

  let removeFromEnd = 0;
  for (; removeFromEnd < records.length; removeFromEnd++) {
    const r = records[records.length - removeFromEnd - 1];
    if (timestampToDate(r.recordTime) <= periodEnd) {
      break; // have reached the last record that is within the period (well is before the end of it)
    }
  }

  // avoid an array copy if we can help it
  if (removeFromStart === 0 && removeFromEnd === 0) {
    return records;
  }

  return records.slice(removeFromStart, records.length - removeFromEnd);
}
