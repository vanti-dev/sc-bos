<template>
  <div class="impact-totals">
    <health-impact-table
        v-for="table in impactTables"
        :key="table.title"
        v-bind="table"/>
  </div>
</template>

<script setup>
import {useComplianceImpactTable, useEquipmentImpactTable, useOccupantImpactTable} from '@/traits/health/health.js';
import HealthImpactTable from '@/traits/health/HealthImpactTable.vue';
import {computed} from 'vue';

const {table: occupantTable} = useOccupantImpactTable();
const {table: equipmentTable} = useEquipmentImpactTable();
const {table: complianceTable} = useComplianceImpactTable();

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