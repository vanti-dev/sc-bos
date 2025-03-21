<template>
  <router-link :to="toAttr" v-tooltip:bottom="tooltipStr">{{ textStr }}</router-link>
</template>

<script setup>
import {ResourceTypeById} from '@/api/ui/account.js';
import {RoleAssignment} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed} from 'vue';

const props = defineProps({
  roleAssignment: {
    type: Object, // type of RoleAssignment.AsObject & {role: Role.AsObject} & {device: Device.AsObject}
    required: true,
  }
});

const toAttr = computed(() => {
  const roleAssignment = props.roleAssignment;
  if (!roleAssignment) return null;
  return {
    name: 'roles',
    params: {
      roleId: roleAssignment.roleId,
    }
  };
});
const textStr = computed(() => {
  let str = '';
  const scope = props.roleAssignment.scope;
  if (!scope) {
    str += 'Global';
  } else {
    const typeStr = (() => {
      const base = (() => {
        const resourceType = scope.resourceType;
        switch (resourceType) {
          case RoleAssignment.ResourceType.NAMED_RESOURCE:
            return 'Device'
          case RoleAssignment.ResourceType.NAMED_RESOURCE_PATH_PREFIX:
            return 'Device+'
          default:
            return ResourceTypeById[scope.resourceType];
        }
      })();
      return base[0] + base.slice(1).toLowerCase();
    })();
    if (props.roleAssignment.device?.metadata?.appearance?.title) {
      str += `${typeStr} ${props.roleAssignment.device.metadata.appearance.title}`;
    } else {
      str += `${typeStr} ${scope.resource}`;
    }
  }
  str += ` ${props.roleAssignment.role?.displayName ?? '...'}`;
  return str;
});
const tooltipStr = computed(() => {
  if (props.roleAssignment.device?.metadata?.appearance?.title) {
    return props.roleAssignment.scope?.resource;
  }
  return undefined;
})
</script>

<style scoped>

</style>