<template>
  <v-card elevation="0" class="rounded-lg">
    <v-toolbar title="Roles" color="transparent" class="px-4 pt-2">
      <template #append>
        <delete-roles-btn
            v-if="showDeleteRolesBtn"
            color="error"
            variant="outlined"
            :roles="selectedRoles"
            @delete="onDelete"/>
        <new-role-btn variant="flat" color="primary" @save="onNewRoleSave"/>
      </template>
    </v-toolbar>
    <v-card-text>
      <v-data-table-server
          v-bind="tableAttrs"
          :headers="tableHeaders"
          disable-sort
          return-object
          show-select
          v-model="selectedRoles">
        <template #top>
          <v-expand-transition>
            <div v-if="tableErrorStr">
              <v-alert type="error" :text="tableErrorStr"/>
            </div>
          </v-expand-transition>
        </template>
        <!-- no header select item -->
        <template #header.data-table-select/>
        <template #item.displayName="{item}">
          <span>{{ item.displayName }}</span>
          <span class="opacity-50 ml-2" v-if="item.description">{{ item.description }}</span>
        </template>
        <template #item.permissions="{item}">
          {{ (item.permissionsList ?? []).length.toLocaleString() }}
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>

<script setup>
import {useDataTableCollection} from '@/composables/table.js';
import {toAddChange, toRemoveChange, useRolesCollection} from '@/routes/auth/accounts.js';
import DeleteRolesBtn from '@/routes/auth/roles/DeleteRolesBtn.vue';
import NewRoleBtn from '@/routes/auth/roles/NewRoleBtn.vue';
import {computed, ref} from 'vue';

// used to fake PullRoles when creating a new role.
// the var isn't reactive, but the value will be whenever it's set
let pullRolesResource = /** @type {ResourceCollection<Role.AsObject, *> | null} */ null;

const wantCount = ref(20);
const rolesCollectionOpts = computed(() => {
  return {
    wantCount: wantCount.value,
    pullFn: (_, resource) => {
      pullRolesResource = resource;
    }
  };
})
const rolesCollection = useRolesCollection({}, rolesCollectionOpts);
const tableAttrs = useDataTableCollection(wantCount, rolesCollection);
const tableHeaders = computed(() => {
  return [
    {title: 'Name', key: 'displayName', maxWidth: '10em', cellProps: {class: 'text-overflow-ellipsis'}},
    {title: 'Permissions', key: 'permissions', align: 'end', maxWidth: '10em', cellProps: {class: 'text-overflow-ellipsis'}},
  ]
});

const tableErrorStr = computed(() => {
  const errors = rolesCollection.errors.value;
  if (errors.length === 0) return null;
  return 'Error fetching roles: ' + errors.map((e) => (e.error ?? e).message ?? e).join(', ');
});

const latestRole = ref(null);

const onNewRoleSave = ({role}) => {
  if (pullRolesResource) {
    pullRolesResource.lastResponse = toAddChange(role);
  }
  latestRole.value = role;
};

const selectedRoles = ref([]);

const showDeleteRolesBtn = computed(() => selectedRoles.value.length > 0);
const onDelete = () => {
  if (pullRolesResource) {
    for (const role of selectedRoles.value) {
      pullRolesResource.lastResponse = toRemoveChange(role);
    }
  }
  const _latest = latestRole.value;
  for (const role of selectedRoles.value) {
    if (_latest && role.id === _latest.id) {
      latestRole.value = null;
    }
  }
  selectedRoles.value = [];
}
</script>

<style scoped>
:deep(.v-toolbar__append) {
  gap: 1rem;
}
</style>