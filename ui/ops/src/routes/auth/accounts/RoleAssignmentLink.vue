<template>
  <router-link :to="toAttr">{{ textStr }}</router-link>
</template>

<script setup>
import {ResourceTypeById} from '@/api/ui/account.js';
import {computed} from 'vue';

const props = defineProps({
  roleAssignment: {
    type: Object, // type of RoleAssignment.AsObject & {role: Role.AsObject}
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
  if (!props.roleAssignment.scope) {
    str += 'Global';
  } else {
    const typeStr = (() => {
      const base = ResourceTypeById[props.roleAssignment.scope.resourceType];
      return base[0] + base.slice(1).toLowerCase();
    })();
    str += `${typeStr} ${props.roleAssignment.scope.resource}`;
  }
  str += ` ${props.roleAssignment.role?.displayName ?? '...'}`;
  return str;
})
</script>

<style scoped>

</style>