<template>
  <thead>
    <tr>
      <th
          v-for="header in props.tableHeaders"
          :key="header.text"
          :align="header?.align"
          :width="header?.width">
        {{ header.text }}
      </th>
    </tr>
  </thead>
</template>

<script setup>

const props = defineProps({
  tableHeaders: {
    type: Array,
    default: () => []
  }
});
</script>
<style lang="scss" scoped>
:deep(.v-data-table-header) {
  margin-left: 8px;
}
</style>
