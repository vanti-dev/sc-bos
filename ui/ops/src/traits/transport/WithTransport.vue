<template>
  <div>
    <slot name="transport" :resource="transportValue" :info="transportInfo"/>
  </div>
  <div>
    <slot name="history" :history="transportHistory"/>
  </div>
</template>

<script setup>
import {timestampFromObject} from '@/api/convpb.js';
import {useDescribeTransport, usePullTransport, useTransportHistory} from '@/traits/transport/transport.js';
import {Period} from '@smart-core-os/sc-api-grpc-web/types/time/period_pb';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullTransportRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const transportValue = reactive(usePullTransport(() => props.name || props.request, () => props.paused));
const transportInfo = reactive(useDescribeTransport(() => props.name));

const computePeriod = () => {
  const period = new Period();
  const now = new Date();
  period.setEndTime(timestampFromObject(now));

  const prevMonth = now.getMonth() - 1;
  if (prevMonth < 0) {
    now.setFullYear(now.getFullYear() - 1);
    now.setMonth(11);
  } else {
    now.setMonth(prevMonth);
  }
  period.setStartTime(timestampFromObject(now));

  return period.toObject();
};

const transportHistory = reactive(useTransportHistory(() => props.name, () => computePeriod()));

</script>
