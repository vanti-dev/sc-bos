<template>
  <v-btn>
    Grant Role...
    <v-menu v-model="menuModel" activator="parent" :close-on-content-click="false">
      <v-card title="Grant Role" width="440px">
        <v-card-text>
          <v-autocomplete
              v-model="selectedRole"
              :items="rolesCollection"
              item-title="displayName"
              item-value="id"
              return-object
              label="Role"
              hide-details>
            <template #item="{ props: _props, item }">
              <v-list-item v-bind="_props">
                <template #append v-if="item.raw.pb_protected">
                  <v-chip text="Built-in" size="small" class="ml-2"/>
                </template>
              </v-list-item>
            </template>
          </v-autocomplete>
        </v-card-text>
        <v-card-text v-if="!props.globalOnly">
          <scope-autocomplete v-model="selectedScopes" :disabled="scopeDisabled">
            <template #appendSticky>
              <v-btn @click="onGrantClick" color="primary" :disabled="grantBtnDisabled">
                Grant
              </v-btn>
              <v-btn @click="onCancelClick">Cancel</v-btn>
            </template>
          </scope-autocomplete>
          <v-expand-transition>
            <v-alert v-if="selectedRoleIsProtected"
                     type="info"
                     tile
                     variant="text"
                     text="Built-in roles cannot be scoped."/>
          </v-expand-transition>
        </v-card-text>
        <v-card-actions>
          <v-btn @click="onGrantClick" color="primary" :disabled="grantBtnDisabled">
            Grant
          </v-btn>
          <v-btn @click="onCancelClick">Cancel</v-btn>
        </v-card-actions>
        <v-expand-transition>
          <div v-if="grantError">
            <v-alert type="error" :text="grantErrorStr" tile class="mt-2"/>
          </div>
        </v-expand-transition>
      </v-card>
    </v-menu>
  </v-btn>
</template>

<script setup>
import {createRoleAssignment} from '@/api/ui/account.js';
import ScopeAutocomplete from '@/routes/auth/accounts/ScopeAutocomplete.vue';
import {useRolesCollection} from '@/routes/auth/roles/roles.js';
import {computed, onScopeDispose, ref} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: undefined
  },
  accounts: {
    type: [String, Object, Array], // String is the id, Object is an account, Array is of string or object.
    required: true,
  },
  globalOnly: {
    type: Boolean,
    default: false, // If true, only global scope assignments are allowed.
  }
});
const emit = defineEmits(['save', 'cancel']);

const selectedRole = ref(null);
const selectedScopes = ref([]);

const selectedRoleIsProtected = computed(() => selectedRole.value?.pb_protected ?? false);

const menuModel = ref(false);

const grantBtnDisabled = computed(() => {
  if (!selectedRole.value) return true;
  if (selectedRoleIsProtected.value) return false; // Protected roles can be granted without scopes.
  if (props.globalOnly) return false; // Only a role is needed for global scope assignments.
  return selectedScopes.value.length === 0;
});
const grantLoading = ref(false);
const grantError = ref(null);
const grantErrorStr = computed(() => {
  const err = grantError.value;
  if (!err) return null;
  return `Failed to grant roles: ${err.error?.message ?? err.message ?? err}`;
});

const _accountIds = computed(() => {
  if (Array.isArray(props.accounts)) {
    return props.accounts.map((a) => typeof a === 'string' ? a : a.id);
  } else {
    return [typeof props.accounts === 'string' ? props.accounts : props.accounts.id];
  }
})
const onGrantClick = async () => {
  grantLoading.value = true;
  const ras = [];
  try {
    for (const accountId of _accountIds.value) {
      if (selectedScopes.value.length === 0) {
        // global scope assignment
        const res = await createRoleAssignment({
          name: props.name,
          roleAssignment: {
            accountId: accountId,
            roleId: selectedRole.value.id,
          }
        });
        ras.push(res);
      }
      for (const scope of selectedScopes.value) {
        const roleAssignment = {
          accountId: accountId,
          roleId: selectedRole.value.id,
        };
        if (scope.type && scope.value) {
          roleAssignment.scope = {
            resource: scope.value,
            resourceType: scope.type,
          };
        }
        const res = await createRoleAssignment({name: props.name, roleAssignment});
        ras.push(res);
      }
    }
    grantError.value = null;
    emit('save', ras);
    hideAndClear();
  } catch (e) {
    grantError.value = e;
  } finally {
    grantLoading.value = false;
  }
}
const onCancelClick = () => {
  emit('cancel');
  hideAndClear();
}
let hideAndClearHandle = 0;
const hideAndClear = () => {
  menuModel.value = false;
  hideAndClearHandle = setTimeout(() => {
    selectedRole.value = null;
    selectedScopes.value = [];
    grantError.value = null;
  }, 250);
}
onScopeDispose(() => {
  clearTimeout(hideAndClearHandle);
});

// role selection
// todo: support pagination and/or filtering on the server
const rolesWantCount = ref(100);
const {items: rolesCollection} = useRolesCollection(() => ({name: props.name}), () => ({
  wantCount: rolesWantCount.value
}));

const scopeDisabled = computed(() => !selectedRole.value || selectedRoleIsProtected.value);
</script>

<style scoped>
</style>