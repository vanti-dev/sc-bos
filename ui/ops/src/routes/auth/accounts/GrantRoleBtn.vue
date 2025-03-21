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
              hide-details/>
        </v-card-text>
        <v-card-text>
          <scope-autocomplete v-model="selectedScopes" :disabled="scopeDisabled"/>
        </v-card-text>
        <v-card-actions>
          <v-btn @click="onGrantClick" color="primary" :disabled="!selectedRole || selectedScopes.length === 0">Grant</v-btn>
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
  }
});
const emit = defineEmits(['save', 'cancel']);

const selectedRole = ref(null);
const selectedScopes = ref([]);

const menuModel = ref(false);

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
      for (const scope of selectedScopes.value) {
        const res = await createRoleAssignment({
          name: props.name,
          roleAssignment: {
            accountId: accountId,
            roleId: selectedRole.value.id,
            scope: {
              resource: scope.value,
              resourceType: scope.type,
            }
          }
        });
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

const scopeDisabled = computed(() => !selectedRole.value);
</script>

<style scoped>
</style>