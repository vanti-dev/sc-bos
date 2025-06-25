<template>
  <v-btn>
    Delete
    <v-dialog activator="parent" max-width="440" v-model="dialogVisible">
      <v-card :title="dialogTitle">
        <v-card-text>Are you sure you want to delete these roles? This cannot be undone!</v-card-text>
        <v-card-actions>
          <v-btn text="Cancel" @click="onCancel" :disabled="deleteLoading"/>
          <v-btn text="Delete" @click="onDelete" type="submit" color="error" variant="flat" :loading="deleteLoading"/>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-btn>
</template>

<script setup>
import {deleteRole} from '@/api/ui/account.js';
import {computed, ref} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: undefined,
  },
  roles: {
    type: Array,
    default: () => [], // of Role.AsObject
  }
});
const emit = defineEmits(['cancel', 'delete']);

const dialogVisible = ref(false);
const dialogTitle = computed(() => {
  if (props.roles.length === 1) {
    return `Delete Role ${props.roles[0].displayName}`;
  }
  return `Delete ${props.roles.length} Roles`;
});
const deleteLoading = ref(false);
const onCancel = () => {
  emit('cancel');
  dialogVisible.value = false;
}
const onDelete = async () => {
  deleteLoading.value = true;
  try {
    const selected = props.roles;
    for (const account of selected) {
      await deleteRole({id: account.id, allowMissing: true});
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