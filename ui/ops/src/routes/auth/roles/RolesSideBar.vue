<template>
  <side-bar @close="onCloseClick">
    <template #actions>
      <v-btn v-if="editable && !editMode" @click="onEditClick" icon="mdi-pencil" variant="plain" size="small"/>
    </template>
    <template v-if="!editMode">
      <p class="px-4 my-4" v-if="role?.description">{{ role.description }}</p>
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
          <v-btn type="submit" variant="flat" color="primary">Save</v-btn>
          <v-btn @click="onCancelClick" variant="flat">Cancel</v-btn>
        </div>
      </v-form>
      <v-expand-transition>
        <div v-if="saveError">
          <v-alert type="error" tile :text="saveErrorStr"/>
        </div>
      </v-expand-transition>
    </template>
    <template v-if="roleIsProtected">
      <v-alert type="info" tile variant="text" text="Built-in roles cannot be modified."/>
    </template>
    <v-list v-else density="compact">
      <v-list-subheader>{{ permissionsListTitle }}</v-list-subheader>
      <v-expand-transition group>
        <v-list-item
            v-for="perm in permissionsToggleList" :key="perm.id"
            density="compact"
            @click="onPermissionClick(perm)">
          <template #prepend>
            <v-checkbox-btn :model-value="perm.assigned" color="primary" density="compact" class="mr-2"/>
          </template>
          <v-list-item-title>{{ perm.displayName }}</v-list-item-title>
          <v-menu activator="parent" location="left" open-on-hover>
            <v-card width="20em">
              <v-card-title class="pb-0 text-wrap">{{ perm.displayName }}</v-card-title>
              <v-card-subtitle>{{ perm.id }}</v-card-subtitle>
              <v-card-text>{{ perm.description }}</v-card-text>
              <template v-if="perm.implies.length > 0">
                <v-card-subtitle class="pt-2">Implies</v-card-subtitle>
                <v-card-text class="pt-0">
                  <div v-for="p of perm.implies" :key="p.id">{{ p.displayName }}</div>
                </v-card-text>
              </template>
              <template v-if="perm.dependsOn.length > 0">
                <v-card-subtitle class="pt-2">Depends On</v-card-subtitle>
                <v-card-text class="pt-0">
                  <div v-for="p of perm.dependsOn" :key="p.id">{{ p.displayName }}</div>
                </v-card-text>
              </template>
            </v-card>
          </v-menu>
        </v-list-item>
      </v-expand-transition>
    </v-list>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {useAssignedPermissions} from '@/routes/auth/roles/roles.js';
import {useSidebarStore} from '@/stores/sidebar.js';
import {computed, ref} from 'vue';
import {useRouter} from 'vue-router';

const sidebar = useSidebarStore();
const role = computed(() => sidebar.data?.role);

const roleIsProtected = computed(() => role.value?.pb_protected ?? false);

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

const {
  toggleList: permissionsToggleList,
  addPermission, removePermission
} = useAssignedPermissions(permissionsList);
const onPermissionClick = async (perm) => {
  const newRole = {
    ...role.value
  }
  if (perm.assigned) {
    newRole.permissionIdsList = removePermission(perm);
  } else {
    newRole.permissionIdsList = addPermission(perm);
  }
  await sidebar.data.updateRole({role: newRole, updateMask: ['permission_ids']});
}

const editable = computed(() => !roleIsProtected.value)
const editMode = ref(false);
const formValid = ref(false);
const saving = ref(false);
const saveError = ref(null);
const saveErrorStr = computed(() => {
  const err = saveError.value;
  if (!err) return null;
  const msg = `${err.error?.message ?? err.message ?? err}`;
  return 'Error saving role: ' + msg;
});

const router = useRouter();
const onCloseClick = () => {
  router.push({name: 'roles'});
}
const onEditClick = () => {
  editDisplayNameModel.value = role.value?.displayName ?? '';
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
      displayName: editDisplayNameModel.value,
      description: editDescriptionModel.value
    }
    await sidebar.data.updateRole({role: newRole, updateMask: ['display_name', 'description']});
    saveError.value = null;
  } catch (e) {
    saveError.value = e;
  } finally {
    saving.value = false;
  }
  editMode.value = false;
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