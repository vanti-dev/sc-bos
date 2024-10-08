<template>
  <side-bar>
    <v-list>
      <account-zone-list-card :zone-list="sidebar.data.zoneNamesList ?? []" @update:zone-list="saveZones"/>
      <v-divider/>
      <account-secrets-card :account="sidebar.data"/>
      <v-divider/>
      <v-list-item class="pt-3">
        <delete-confirmation-dialog
            title="Delete Account"
            :progress-bar="deleteTracker.loading"
            @confirm="deleteAccount">
          Are you sure you want to delete the account "{{ sidebar.title }}"?
          <template #alert-content>
            Deleting this account will stop all integrations that connect using this account.
            <br><br>
            This action cannot be undone.
          </template>
          <template #confirmBtn>I understand, delete account</template>
          <template #activator="{ props }">
            <v-btn
                variant="outlined"
                color="error"
                :disabled="blockActions"
                width="100%"
                v-bind="props">
              Delete Account
            </v-btn>
          </template>
        </delete-confirmation-dialog>
      </v-list-item>
    </v-list>
  </side-bar>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import {deleteTenant, updateTenant} from '@/api/ui/tenant';
import SideBar from '@/components/SideBar.vue';
import useAuthSetup from '@/composables/useAuthSetup';
import AccountSecretsCard from '@/routes/auth/third-party/components/AccountSecretsCard.vue';
import AccountZoneListCard from '@/routes/auth/third-party/components/AccountZoneListCard.vue';
import DeleteConfirmationDialog from '@/routes/auth/third-party/components/DeleteConfirmationDialog.vue';
import {useTenantStore} from '@/routes/auth/third-party/tenantStore';
import {useSidebarStore} from '@/stores/sidebar';
import {reactive} from 'vue';

defineProps({
  // This is passed as part of the routing, but we aren't currently using it.
  accountId: {
    type: String,
    default: ''
  }
});

const sidebar = useSidebarStore();
const tenantStore = useTenantStore();

const deleteTracker = reactive(
    /** @type {ActionTracker<DeleteTenantResponse.AsObject>} */ newActionTracker()
);
const updateZonesTracker = reactive(
    /** @type {ActionTracker<Tenant.AsObject>} */
    newActionTracker()
);

/**
 *
 */
async function deleteAccount() {
  await deleteTenant(
      {
        id: sidebar.data.id
      },
      deleteTracker
  );
  tenantStore.refreshTenants();
  sidebar.visible = false;
}

/**
 * @param {string[]} zones
 * @return {Promise<void>}
 */
async function saveZones(zones) {
  await updateTenant({
    tenant: {
      id: sidebar.data.id,
      zoneNamesList: zones
    },
    updateMask: {
      pathsList: ['zone_names']
    }
  }, updateZonesTracker);
  tenantStore.refreshTenants();
}

// ------------------------------ //
// ----- Authentication settings ----- //

const {blockActions} = useAuthSetup();
</script>

<style scoped></style>
