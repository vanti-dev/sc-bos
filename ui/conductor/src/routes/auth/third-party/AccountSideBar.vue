<template>
  <side-bar>
    <v-list>
      <account-zone-list-card :zone-list="sidebarData.zoneNamesList ?? []"/>
      <v-divider/>
      <account-secrets-card :account-id="sidebarData.id"/>
      <v-divider/>
      <v-list-item class="pt-3">
        <v-dialog v-model="deleteConfirmation" max-width="320">
          <v-card class="pa-2">
            <v-card-title class="text-h4 error--text text--lighten">Delete Account</v-card-title>
            <v-card-text>
              Are you sure you want to delete the account "{{ sidebarTitle }}"?<br><br>
              <span class="font-bold error--text">Note: This action cannot be undone</span>
            </v-card-text>
            <v-card-actions>
              <v-spacer/>
              <v-btn @click="deleteConfirmation = false" color="primary">Cancel</v-btn>
              <v-btn @click="deleteAccount" color="error">Delete</v-btn>
            </v-card-actions>
          </v-card>
          <template #activator="{on}">
            <v-btn outlined color="error" width="100%" v-on="on">Delete Account</v-btn>
          </template>
        </v-dialog>
      </v-list-item>
    </v-list>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {reactive, ref} from 'vue';
import {deleteTenant} from '@/api/ui/tenant';
import {newActionTracker} from '@/api/resource';
import AccountZoneListCard from '@/routes/auth/third-party/components/AccountZoneListCard.vue';
import AccountSecretsCard from '@/routes/auth/third-party/components/AccountSecretsCard.vue';

const pageStore = usePageStore();
const {sidebarTitle, sidebarData} = storeToRefs(pageStore);

const deleteConfirmation = ref(false);
const deleteTracker = reactive(
    /** @type {ActionTracker<DeleteTenantResponse.AsObject>} */ newActionTracker()
);

/**
 *
 */
function deleteAccount() {
  deleteTenant({
    id: sidebarData.value.id
  }, deleteTracker);
  deleteConfirmation.value = false;
  // todo: remove from tenants list
}

</script>

<style scoped>

</style>
