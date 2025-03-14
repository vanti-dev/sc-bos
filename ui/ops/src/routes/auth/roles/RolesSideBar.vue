<template>
  <side-bar>
    <template #actions>
      <v-btn v-if="!editMode" @click="onEditClick" icon="mdi-pencil" variant="plain" size="small"/>
    </template>
    <template v-if="!editMode">
      <p class="px-4 my-4" v-if="role?.description">{{ role.description }}</p>
    </template>
    <template v-else>
      <v-textarea
          class="ma-4"
          v-model="editDescriptionModel"
          label="Description"
          :counter="250"
          autocomplete="off"/>
      <div class="d-flex ga-2 ma-4">
        <v-btn @click="onSaveClick" variant="flat" color="primary">Save</v-btn>
        <v-btn @click="onCancelClick" variant="flat">Cancel</v-btn>
      </div>
    </template>
    <v-list density="compact">
      <v-list-subheader>{{ permissionsListTitle }}</v-list-subheader>
      <v-list-item v-for="perm in permissionsList" :key="perm">
        <v-list-item-title>{{ perm }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {computed, ref} from 'vue';

const sidebar = useSidebarStore();
const role = computed(() => sidebar.data?.role);
const permissionsList = computed(() => role.value?.permissionIdsList ?? []);
const permissionsListTitle = computed(() => {
  const len = permissionsList.value.length;
  switch (len) {
    case 0:
      return 'No Permissions';
    case 1:
      return '1 Permission';
    default:
      return `${len} Permissions`;
  }
});

const editMode = ref(false);
const saving = ref(false);
const saveError = ref(null);
const onEditClick = () => {
  editDescriptionModel.value = role.value?.description ?? '';
  editMode.value = true;
};
const onCancelClick = () => {
  editMode.value = false;
};
const onSaveClick = async () => {
  try {
    saving.value = true;
    const newRole = {
      ...role.value,
      description: editDescriptionModel.value
    }
    await sidebar.data.updateRole({role: newRole});
    saveError.value = null;
  } catch (e) {
    saveError.value = e;
  } finally {
    saving.value = false;
  }
  editMode.value = false;
};

const editDescriptionModel = ref(null);
</script>

<style scoped>

</style>