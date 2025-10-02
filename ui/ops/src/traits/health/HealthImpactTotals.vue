<template>
  <div class="impact-totals">
    <health-impact-table
        v-for="table in impactTables"
        :key="table.title"
        v-bind="table"/>
  </div>
</template>

<script setup>
import {newResourceValue} from '@/api/resource.js';
import {pullDevicesMetadata} from '@/api/ui/devices';
import {
  useComplianceImpactTable,
  useEquipmentImpactTable,
  useOccupantImpactTable,
  useRollingHistory
} from '@/traits/health/health.js';
import HealthImpactTable from '@/traits/health/HealthImpactTable.vue';
import {computed, reactive} from 'vue';

// tracks how many abnormal checks are in each category
const abnormalCounts = reactive(newResourceValue())
pullDevicesMetadata({
  query: {
    conditionsList: [
      {field: 'health_checks.check.state', stringIn: {stringsList: ['ABNORMAL', 'HIGH', 'LOW']}},
    ]
  },
  includes: {
    fieldsList: [
      'health_checks.id',
      'health_checks.occupant_impact',
      'health_checks.equipment_impact',
      'health_checks.compliance_impacts.contribution',
    ]
  }
}, abnormalCounts);

const {oldValue: initialCounts} = useRollingHistory(() => abnormalCounts.value);

const {table: occupantTable} = useOccupantImpactTable(abnormalCounts, initialCounts);
const {table: equipmentTable} = useEquipmentImpactTable(abnormalCounts, initialCounts);
const {table: complianceTable} = useComplianceImpactTable(abnormalCounts, initialCounts);

const impactTables = computed(() => {
  return [
    occupantTable.value,
    equipmentTable.value,
    complianceTable.value,
  ].filter(t => t.totalCount > 0);
})
</script>

<style scoped lang="scss">
.impact-totals {
  display: flex;
  justify-content: stretch;
  flex-wrap: wrap;
  gap: 40px;

  > * {
    flex: 1;
    min-width: 30em;
    max-width: 40em;
  }
}
</style>