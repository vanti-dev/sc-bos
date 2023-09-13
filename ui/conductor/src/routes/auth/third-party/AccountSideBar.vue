<template>
  <side-bar>
    <v-list>
      <account-zone-list-card :zone-list="sidebarData.zoneNamesList ?? []"/>
      <v-divider/>
      <account-secrets-card :account="sidebarData"/>
      <v-divider/>
      <v-list-item class="pt-3">
        <delete-confirmation-dialog
            title="Delete Account"
            :progress-bar="deleteTracker.loading"
            @confirm="deleteAccount">
          Are you sure you want to delete the account "{{ sidebarTitle }}"?
          <template #alert-content>
            Deleting this account will stop all integrations that connect using this account.
            <br><br>
            This action cannot be undone.
          </template>
          <template #confirmBtn>I understand, delete account</template>
          <template #activator="{ on, attrs }">
            <v-btn
                outlined
                color="error"
                :disabled="blockActions"
                width="100%"
                v-on="on"
                v-bind="attrs">
              Delete Account
            </v-btn>
          </template>
        </delete-confirmation-dialog>
      </v-list-item>
    </v-list>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {reactive} from 'vue';
import {deleteTenant} from '@/api/ui/tenant';
import {newActionTracker} from '@/api/resource';
import AccountZoneListCard from '@/routes/auth/third-party/components/AccountZoneListCard.vue';
import AccountSecretsCard from '@/routes/auth/third-party/components/AccountSecretsCard.vue';
import DeleteConfirmationDialog from '@/routes/auth/third-party/components/DeleteConfirmationDialog.vue';
import {useTenantStore} from '@/routes/auth/third-party/tenantStore';
import useAuthSetup from '@/composables/useAuthSetup';

const pageStore = usePageStore();
const {sidebarTitle, sidebarData} = storeToRefs(pageStore);
const tenantStore = useTenantStore();

const deleteTracker = reactive(/** @type {ActionTracker<DeleteTenantResponse.AsObject>} */ newActionTracker());

/**
 *
 */
async function deleteAccount() {
  await deleteTenant(
      {
        id: sidebarData.value.id
      },
      deleteTracker
  );
  tenantStore.refreshTenants();
  pageStore.showSidebar = false;
}

// ------------------------------ //
// ----- Authentication settings ----- //

const {blockActions} = useAuthSetup();
</script>

<style scoped></style>
