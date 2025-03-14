<template>
  <side-bar @close="onCloseClick">
    <template #actions>
      <v-btn v-if="!editMode" @click="onEditClick" icon="mdi-pencil" variant="plain" size="small"/>
    </template>
    <template v-if="!editMode">
      <p class="px-4 my-4" v-if="account?.description">{{ account.description }}</p>
    </template>
    <template v-else>
      <v-form @submit.prevent="onSaveClick" v-model="formValid">
        <v-text-field
            class="ma-4"
            v-model="editDisplayNameModel"
            :rules="editDisplayNameRules"
            :counter="100"
            label="Name"
            autocomplete="off"/>
        <v-textarea
            class="ma-4"
            v-model="editDescriptionModel"
            label="Description"
            :counter="250"
            autocomplete="off"/>
        <div class="d-flex ga-2 ma-4">
          <v-btn @click="onSaveClick" variant="flat" color="primary" type="submit" :disabled="!formValid">Save</v-btn>
          <v-btn @click="onCancelClick" variant="flat">Cancel</v-btn>
        </div>
      </v-form>
      <v-expand-transition>
        <div v-if="saveError">
          <v-alert type="error" tile :text="saveErrorStr"/>
        </div>
      </v-expand-transition>
    </template>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {computed, ref} from 'vue';
import {useRouter} from 'vue-router';

const sidebar = useSidebarStore();
const account = computed(() => sidebar.data?.account);

const editMode = ref(false);
const formValid = ref(false);
const saving = ref(false);
const saveError = ref(null);
const saveErrorStr = computed(() => {
  const err = saveError.value;
  if (!err) return null;
  const msg = `${err.error?.message ?? err.message ?? err}`;
  return 'Error saving account: ' + msg;
})

const router = useRouter();
const onCloseClick = () => {
  router.push({name: 'accounts'});
}
const onEditClick = () => {
  editDisplayNameModel.value = account.value?.displayName ?? '';
  editDescriptionModel.value = account.value?.description ?? '';
  editMode.value = true;
};
const onCancelClick = () => {
  editMode.value = false;
};
const onSaveClick = async () => {
  try {
    saving.value = true;
    const newAccount = {
      ...account.value,
      displayName: editDisplayNameModel.value,
      description: editDescriptionModel.value,
    }
    await sidebar.data.updateAccount({account: newAccount});
    saveError.value = null;
    editMode.value = false;
  } catch (e) {
    saveError.value = e;
  } finally {
    saving.value = false;
  }
};

const editDisplayNameModel = ref(null);
const editDisplayNameRules = computed(() => {
  return [
    (v) => !!v || 'Name is required',
    (v) => v.length >= 3 || 'Name must be at least 3 characters',
    (v) => v.length <= 100 || 'Name must be less than 100 characters',
  ];
})
const editDescriptionModel = ref(null);
</script>

<style scoped>

</style>