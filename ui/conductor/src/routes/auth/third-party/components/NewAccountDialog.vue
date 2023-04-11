<template>
  <v-dialog v-model="dialog" max-width="512">
    <template #activator="actBindings">
      <slot name="activator" v-bind="actBindings"/>
    </template>
    <v-card class="px-2 pb-4 pt-1">
      <v-card-title>Add Tenant</v-card-title>
      <v-list>
        <v-list-item>
          <v-text-field label="Name" v-model="name" filled hide-details/>
        </v-list-item>
      </v-list>
      <v-card-actions>
        <v-spacer/>
        <v-btn color="error" @click="cancel">Cancel</v-btn>
        <v-btn color="primary" @click="addTenant">Add</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup>

import {onMounted, onUnmounted, ref} from 'vue';
import {createTenant} from '@/api/ui/tenant';
import {newActionTracker} from '@/api/resource';
import {useErrorStore} from '@/components/ui-error/error';

const dialog = ref(false);
const addTenantTracker = ref(newActionTracker());
const name = ref('');

const emit = defineEmits(['finished']);

// UI error handling
const errorStore = useErrorStore();
let unwatchErrors;
onMounted(() => {
  unwatchErrors = errorStore.registerTracker(addTenantTracker);
});

onUnmounted(() => {
  unwatchErrors();
  clearForm();
});

/**
 *
 */
function clearForm() {
  name.value = '';
}

/**
 *
 */
function cancel() {
  clearForm();
  dialog.value = false;
}

/**
 *
 */
async function addTenant() {
  const req = {
    tenant: {
      title: name.value
    }
  };
  await createTenant(req, addTenantTracker);
  clearForm();
  dialog.value = false;
  emit('finished');
}

</script>

<style scoped>

</style>
