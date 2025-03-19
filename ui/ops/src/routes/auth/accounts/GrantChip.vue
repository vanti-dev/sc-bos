<template>
  <v-chip v-tooltip:bottom="tooltipStr">
    <template #prepend>{{ prependStr }}</template>
    {{ valueStr }}
  </v-chip>
</template>

<script setup>
import {ResourceTypeById} from '@/api/ui/account.js';
import {RoleAssignment} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed} from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: undefined
  },
  value: {
    type: String,
    default: undefined
  },
  type: {
    type: Number,
    default: RoleAssignment.ResourceType.RESOURCE_TYPE_UNSPECIFIED
  }
});

const prependStr = computed(() => {
  switch (props.type) {
    case RoleAssignment.ResourceType.NAMED_RESOURCE:
      return 'ID';
    case RoleAssignment.ResourceType.NAMED_RESOURCE_PATH_PREFIX:
      return 'ID+';
    case RoleAssignment.ResourceType.NODE:
      return 'On Node';
    case RoleAssignment.ResourceType.SUBSYSTEM:
      return 'System';
    case RoleAssignment.ResourceType.ZONE:
      return 'Zone';
    default:
      return ResourceTypeById[props.type] ?? 'Unknown';
  }
});
const valueStr = computed(() => {
  return props.title;
});
const tooltipStr = computed(() => {
  return props.value ?? props.title;
})
</script>

<style scoped>
:deep(.v-chip__prepend) {
  height: 100%;
  margin-left: -10px;
  padding-left: 10px;
  margin-right: .5em;
  padding-right: .5em;
  background: rgba(var(--v-theme-primary), 0.4);
}
</style>