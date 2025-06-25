<template>
  <v-btn>
    Delete
    <v-dialog activator="parent" max-width="440" v-model="dialogVisible">
      <v-card :title="dialogTitle">
        <v-card-text>Are you sure you want to delete these accounts? This cannot be undone!</v-card-text>
        <v-card-actions>
          <v-btn text="Cancel" @click="onCancel" :disabled="deleteLoading"/>
          <v-btn text="Delete" @click="onDelete" type="submit" color="error" variant="flat" :loading="deleteLoading"/>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-btn>
</template>

<script setup>
import {deleteAccount} from '@/api/ui/account.js';
import {computed, ref} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: undefined,
  },
  accounts: {
    type: Array,
    default: () => [], // of Account.AsObject
  }
});
const emit = defineEmits(['cancel', 'delete']);

const dialogVisible = ref(false);
const dialogTitle = computed(() => {
  if (props.accounts.length === 1) {
    return `Delete Account ${props.accounts[0].displayName}`;
  }
  return `Delete ${props.accounts.length} Accounts`;
});
const deleteLoading = ref(false);
const onCancel = () => {
  emit('cancel');
  dialogVisible.value = false;
}
const onDelete = async () => {
  deleteLoading.value = true;
  try {
    const selected = props.accounts;
    for (const account of selected) {
      await deleteAccount({id: account.id, allowMissing: true});
    }
    emit('delete');
    dialogVisible.value = false;
  } finally {
    deleteLoading.value = false;
  }
}
</script>

<style scoped>

</style>