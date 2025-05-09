<template>
  <v-card min-width="440">
    <v-form @submit.prevent="onSave"
            v-model="formValid"
            ref="formRef"
            :disabled="formDisabled">
      <v-card-title>New Role</v-card-title>
      <v-card-text class="ga-2 d-flex flex-column">
        <v-text-field
            label="Name"
            v-model.trim="displayName"
            :rules="displayNameRules"
            required
            :counter="100"
            autocomplete="off"/>
        <v-textarea
            label="Description"
            v-model.trim="description"
            :rules="descriptionRules"
            :counter="250"
            autocomplete="off"/>
      </v-card-text>
      <v-card-actions>
        <v-btn text="Create"
               color="primary"
               type="submit"
               variant="flat"
               :disabled="!formValid"
               :loading="saveLoading"/>
        <v-btn text="Cancel" @click="onCancel"/>
      </v-card-actions>
      <v-expand-transition>
        <div v-if="errorStr">
          <v-alert type="error" :text="errorStr" tile class="mt-4"/>
        </div>
      </v-expand-transition>
    </v-form>
  </v-card>
</template>

<script setup>
import {createRole} from '@/api/ui/account.js';
import {computed, ref} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: undefined,
  }
});

const emit = defineEmits(['save', 'cancel', 'error']);

const formRef = ref(null);
const formValid = ref(false);
const displayName = ref('');
const displayNameRules = computed(() => {
  return [
    (v) => !!v || 'Name is required',
    (v) => v.length >= 3 || 'Name must be at least 3 characters',
    (v) => v.length <= 100 || 'Name must be less than 100 characters',
    // todo: add validation to check for uniqueness of the name.
    //   Blocked on an API to check this
  ];
})
const description = ref('');
const descriptionRules = computed(() => {
  return [];
});

const saveError = ref(null);
const errorStr = computed(() => {
  if (!saveError.value) return null;
  return saveError.value.message;
})

const formDisabled = ref(false);
const saveLoading = ref(false);
const reset = () => {
  const form = formRef.value;
  if (!form) return;
  form.reset();
  saveError.value = null;
}
const onSave = async () => {
  saveLoading.value = true;
  formDisabled.value = true;
  try {
    const createRoleReq = {
      name: props.name,
      role: {
        displayName: displayName.value,
        description: description.value,
      }
    }
    const role = await createRole(createRoleReq);
    emit('save', {role});
    saveError.value = null;
  } catch (e) {
    saveError.value = e;
  } finally {
    saveLoading.value = false;
    formDisabled.value = false;
  }
};
const onCancel = () => {
  emit('cancel');
}

defineExpose({
  reset,
})
</script>

<style scoped>

</style>