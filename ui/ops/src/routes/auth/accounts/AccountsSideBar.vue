<template>
  <side-bar @close="onCloseClick">
    <template #actions>
      <v-btn v-if="!editMode" @click="onEditClick" icon="mdi-pencil" variant="plain" size="small"/>
    </template>
    <template v-if="!editMode">
      <div v-if="account" class="details pt-4">
        <p v-if="account.description" class="px-4">{{ account.description }}</p>
        <template v-if="account.username">
          <v-list-subheader>Username</v-list-subheader>
          <p class="px-4">{{ account.username }}</p>
        </template>
        <template v-if="account.type === Account.Type.SERVICE_ACCOUNT">
          <v-list-subheader>Client ID</v-list-subheader>
          <p class="px-4">
            <copy-div :text="account.id" icon-size="18" :btn-props="{size: 'small', class: 'ma-n2'}"/>
          </p>
        </template>
        <v-list-subheader>Account Created</v-list-subheader>
        <p class="px-4">{{ timestampToDate(account.createTime).toLocaleString() }}</p>
      </div>
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
    <v-list>
      <v-list-subheader>
        {{ roleAssignments.length }} Role{{ roleAssignments.length === 1 ? '' : 's' }}
      </v-list-subheader>
      <v-list-item v-for="roleAssignment in roleAssignments" :key="roleAssignment.id" min-height="0px">
        <v-list-item-title>
          <role-assignment-link :role-assignment="roleAssignment"/>
        </v-list-item-title>
        <template #append>
          <v-list-item-action>
            <v-btn
                icon="mdi-delete"
                size="small"
                variant="plain"
                v-tooltip:bottom="'Remove role'"
                class="my-n2"
                :loading="roleDeleteLoading[roleAssignment.id]"
                @click="onRoleRemoveClick(roleAssignment)"/>
          </v-list-item-action>
        </template>
      </v-list-item>
    </v-list>
  </side-bar>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import CopyDiv from '@/components/CopyDiv.vue';
import SideBar from '@/components/SideBar.vue';
import RoleAssignmentLink from '@/routes/auth/accounts/RoleAssignmentLink.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {Account} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed, reactive, ref} from 'vue';
import {useRouter} from 'vue-router';

const sidebar = useSidebarStore();
const account = computed(() => sidebar.data?.account);
const roleAssignmentsCollection = computed(() => sidebar.data?.roleAssignments);

const roleAssignments = computed(() => roleAssignmentsCollection.value?.items ?? []);

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

const roleDeleteLoading = reactive({});
const roleDeleteError = reactive({});
const onRoleRemoveClick = async (roleAssignment) => {
  roleDeleteLoading[roleAssignment.id] = true;
  try {
    delete roleDeleteError[roleAssignment.id];
    await sidebar.data.removeRole(roleAssignment);
  } catch (e) {
    roleDeleteError[roleAssignment.id] = e;
  } finally {
    delete roleDeleteLoading[roleAssignment.id];
  }
}
</script>

<style scoped>
.details .v-list-subheader {
  min-height: initial;
}
</style>