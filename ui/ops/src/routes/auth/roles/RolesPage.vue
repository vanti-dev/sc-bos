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
          v-model="selectedRoles"
          :row-props="tableRowProps"
          @click:row="onRowClick">
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
import {updateRole} from '@/api/ui/account.js';
import {useDataTableCollection} from '@/composables/table.js';
import {toAddChange, toRemoveChange, toUpdateChange, useGetRole, useRolesCollection} from '@/routes/auth/accounts.js';
import DeleteRolesBtn from '@/routes/auth/roles/DeleteRolesBtn.vue';
import NewRoleBtn from '@/routes/auth/roles/NewRoleBtn.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {computed, ref, watch} from 'vue';
import {useRouter} from 'vue-router';

const props = defineProps({
  roleId: {
    type: String,
    default: null,
  }
});

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
});
const rolesCollection = useRolesCollection({}, rolesCollectionOpts);
const tableAttrs = useDataTableCollection(wantCount, rolesCollection);
const tableHeaders = computed(() => {
  return [
    {title: 'Name', key: 'displayName', maxWidth: '10em', cellProps: {class: 'text-overflow-ellipsis'}},
    {title: 'Permissions', key: 'permissions', align: 'end', maxWidth: '10em'},
  ]
});

const tableErrorStr = computed(() => {
  const errors = rolesCollection.errors.value;
  if (errors.length === 0) return null;
  return 'Error fetching roles: ' + errors.map((e) => (e.error ?? e).message ?? e).join(', ');
});
const tableRowProps = ({item}) => {
  return {
    class: {
      'row-selected': selectedRoles.value.includes(item),
      'row-active': item.id === props.roleId,
    }
  };
};

const router = useRouter();
const onRowClick = (_, {item}) => {
  if (item.id === props.roleId) return; // don't click on the same item
  router.push({name: 'roles', params: {roleId: item.id}});
}
const sidebar = useSidebarStore();
const {response: sidebarItem, refresh: refreshSidebarItem} = useGetRole(() => {
  if (!props.roleId) return null;
  return {id: props.roleId};
});
watch(sidebarItem, (item) => {
  if (!item) {
    sidebar.closeSidebar();
    return;
  }
  sidebar.title = item.displayName || `Role ${props.roleId}`;
  sidebar.data = {role: item, updateRole: onRoleUpdate};
  sidebar.visible = true;
}, {immediate: true});
watch(() => sidebar.visible, (visible) => {
  if (!visible && props.roleId) {
    router.push({name: 'roles'});
  }
});

const latestRole = ref(null);

const onNewRoleSave = ({role}) => {
  if (pullRolesResource) {
    pullRolesResource.lastResponse = toAddChange(role);
  }
  latestRole.value = role;
};
const onRoleUpdate = async ({role}) => {
  const oldRole = (() => {
    if (sidebarItem.value?.id === role.id) return sidebarItem.value;
    return rolesCollection.items.value.find((r) => r.id === role.id);
  })();

  const newRole = await updateRole({role})

  if (!oldRole) {
    // we can't do a dynamic of the role, so just refresh the whole list
    rolesCollection.refresh();
  } else {
    // this will update the role in the collection
    pullRolesResource.lastResponse = toUpdateChange(oldRole, newRole);
  }
  if (latestRole.value?.id === role.id) {
    latestRole.value = newRole;
  }
  selectedRoles.value = selectedRoles.value.map((r) => {
    if (r.id === role.id) return newRole;
    return r;
  });
  refreshSidebarItem();

  return newRole;
}

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
    if (role.id === props.roleId) {
      router.push({name: 'roles'});
    }
  }
  selectedRoles.value = [];
}
</script>

<style scoped lang="scss">
:deep(.v-toolbar__append) {
  gap: 1rem;
}

.v-data-table {
  :deep(.row-selected) {
    background-color: rgba(var(--v-theme-primary), 0.1);
  }

  :deep(.row-active) {
    background-color: rgba(var(--v-theme-primary), 0.4);
  }
}
</style>