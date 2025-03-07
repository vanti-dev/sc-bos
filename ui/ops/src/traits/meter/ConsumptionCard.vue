<template>
  <v-card>
    <v-tooltip activator="parent" location="bottom">
      <!-- eslint-disable-next-line vue/no-v-html -->
      <span v-html="tooltipStr"/>
    </v-tooltip>
    <v-card-title v-if="titleStr">{{ titleStr }}</v-card-title>
    <template v-if="isMeterReset">
      <v-card-text class="text-h6 font-weight-light">The meter was reset</v-card-text>
    </template>
    <template v-else>
      <v-card-text>
        <span class="value text-h1 font-weight-light">{{ usageStr }}</span>
        <span v-if="unit" class="unit text-subtitle-1 font-weight-light">&nbsp;{{ unit }}</span>
      </v-card-text>
    </template>
    <v-card-text class="text-body-1 pt-0">{{ timePeriodStr }}</v-card-text>
  </v-card>
</template>

<script setup>
import {usePeriod} from '@/composables/time.js';
import {usePullMetadata} from '@/traits/metadata/metadata.js';
import {
  useDescribeMeterReading,
  useMeterReading,
  useMeterReadingAt,
  usePullMeterReading
} from '@/traits/meter/meter.js';
import {isNullOrUndef} from '@/util/types.js';
import {computed, effectScope, reactive, toRef, watch} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  paused: {
    type: Boolean,
    default: false
  },
  period: {
    type: [String],
    default: 'day' // 'minute', 'hour', 'day', 'month', 'year'
  },
  offset: {
    type: [Number, String],
    default: 0 // Used via Math.abs, {period: 'day', offset: 1} means yesterday, and so on
  },
  title: {
    type: String,
    default: null // default to metadata.appearance.title
  }
});
const _offset = computed(() => -Math.abs(parseInt(props.offset)));
const {start, end} = usePeriod(toRef(props, 'period'), toRef(props, 'period'), _offset)

const {response: meterReadingInfo} = useDescribeMeterReading(() => props.name);

// calculate the reading at the start date, which we assume is in the past.
const readingAtStart = useMeterReadingAt(() => props.name, start);

// calculate the reading at the end.
// If offset is 0 then the reading is the live reading and we can use pullMeterReading.
// If not then we have to fetch the reading from the history.
// We assume that the offset is old enough that history includes the values we need without rechecking.
const endIsLive = computed(() => _offset.value === 0);
let endCalcScope = null;
const readingAtEnd = reactive({value: null});
watch(endIsLive, (endIsLive) => {
  if (endCalcScope) {
    endCalcScope();
  }
  const scope = effectScope();
  endCalcScope = () => scope.stop();
  scope.run(() => {
    if (endIsLive) {
      const {value: meterReading} = usePullMeterReading(() => props.name);
      readingAtEnd.value = meterReading;
    } else {
      readingAtEnd.value = useMeterReadingAt(() => props.name, end);
    }
  });
}, {immediate: true});

const readingDiff = computed(() => {
  const start = readingAtStart.value;
  const end = readingAtEnd.value;
  if (isNullOrUndef(start) || isNullOrUndef(end)) {
    return null;
  }
  return {usage: end.usage - start.usage};
});
const isMeterReset = computed(() => readingDiff.value?.usage < 0);

// Properties for the tooltip, e.g. "On 24th: 100 kWh / On 25th: 200 kWh"
const {usageStr: startUsageStr} = useMeterReading(readingAtStart, meterReadingInfo);
const {usageStr: endUsageStr} = useMeterReading(() => readingAtEnd.value, meterReadingInfo);
const startStr = computed(() => {
  if (!start.value) return 'loading';
  switch (props.period) {
    case 'minute':
    case 'hour':
      return 'At ' + start.value.toLocaleTimeString();
    default:
      return 'On ' + start.value.toLocaleDateString();
  }
});
const endStr = computed(() => {
  if (_offset.value === 0) return 'Last reading';
  switch (props.period) {
    case 'minute':
    case 'hour':
      return 'At ' + end.value.toLocaleTimeString();
    default:
      return 'On ' + end.value.toLocaleDateString();
  }
})
const tooltipStr = computed(() => {
  const unitStr = meterReadingInfo.value?.unit ?? '';
  return `${startStr.value}: ${startUsageStr.value} ${unitStr}<br/>${endStr.value}: ${endUsageStr.value} ${unitStr}`;
});

// The card title is either the title prop or grabbed from the metadata
const {value: md} = usePullMetadata(() => props.name, () => Boolean(props.title));
const titleStr = computed(() => {
  if (props.title) {
    return props.title;
  }
  return md.value?.appearance?.title ?? '';
})
const {usageStr, unit} = useMeterReading(readingDiff, meterReadingInfo);
// outputs text like "yesterday" or "2 days ago" or "last month"
const relativeTimeFormat = new Intl.RelativeTimeFormat(undefined, {
  numeric: 'auto',
  style: 'long'
});
const timePeriodStr = computed(() => {
  if (_offset.value === 0) {
    switch (props.period) {
      case 'minute':
        return 'So far this minute';
      case 'hour':
        return 'So far this hour';
      case 'day':
        return 'So far today';
      case 'week':
        return 'So far this week';
      case 'month':
        return 'So far this month';
      case 'year':
        return 'So far this year';
    }
    return 'So far';
  } else {
    return relativeTimeFormat.format(_offset.value, props.period);
  }
});
</script>

<style scoped>

</style>