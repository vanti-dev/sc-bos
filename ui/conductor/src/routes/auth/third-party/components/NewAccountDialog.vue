<template>
  <v-dialog v-model="dialog" max-width="512">
    <template #activator="actBindings">
      <slot name="activator" v-bind="actBindings"/>
    </template>
    <v-form @submit.prevent="addTenant">
      <v-card class="px-2 pb-4 pt-1">
        <v-card-title>Add Tenant</v-card-title>
        <v-list>
          <v-list-item>
            <v-text-field label="Name" v-model="name" variant="filled" hide-details/>
          </v-list-item>
        </v-list>
        <v-card-actions>
          <v-spacer/>
          <v-btn color="error" @click="cancel">Cancel</v-btn>
          <v-btn color="primary" :disabled="!name" type="submit">Add</v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
<script setup>

import {newActionTracker} from '@/api/resource';
import {createTenant} from '@/api/ui/tenant';
import {useErrorStore} from '@/components/ui-error/error';
import {onMounted, onUnmounted, ref} from 'vue';

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
  if (name.value === '') {
    return;
  }

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
