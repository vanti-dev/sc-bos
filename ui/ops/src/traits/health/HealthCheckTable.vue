<template>
  <v-data-table-server
      v-bind="tableAttrs"
      :headers="headers"
      no-data-text="No health checks">
    <template #item.device="{ item }">
      <md-text :value="item.metadata"/>
    </template>
    <template #item.totalCount="{ item }">
      <check-count-cell v-bind="countChecks(item.healthChecksList)"/>
    </template>
  </v-data-table-server>
</template>

<script setup>
import MdText from '@/components/MdText.vue';
import {useDevices} from '@/composables/devices.js';
import {useDataTableCollection} from '@/composables/table.js';
import CheckCountCell from '@/traits/health/CheckCountCell.vue';
import {countChecks} from '@/traits/health/health';
import {computed, ref} from 'vue';

const props = defineProps({
  conditions: {
    type: Array, // of Device.Query.Condition.AsObject
    default: () => ([
      {'field': 'health_checks.normality', 'stringIn': {'stringsList': ['ABNORMAL', 'HIGH', 'LOW']}}
    ])
  }
});

const wantCount = ref(20);
const _useDevicesOpts = computed(() => {
  return {
    conditions: props.conditions,
    wantCount: wantCount.value
  }
});
const devices = useDevices(_useDevicesOpts);
const tableAttrs = useDataTableCollection(wantCount, devices);
const headers = computed(() => {
  return [
    {title: 'Device', key: 'device'},
    {title: 'Issue Count', key: 'totalCount', align: 'end'},
  ]
})
</script>

<style scoped>
.v-data-table {
  background: transparent;
}
</style>