<template>
  <v-card elevation="0" class="rounded-lg">
    <v-toolbar title="Accounts" color="transparent" class="px-4 pt-2">
      <template #append>
        <delete-accounts-btn
            v-if="showDeleteAccountsBtn"
            color="error"
            variant="outlined"
            :accounts="selectedAccounts"
            @delete="onDelete"/>
        <v-btn variant="flat" color="primary">
          New Account...
          <v-menu v-model="newAccountMenu" activator="parent" :close-on-content-click="false">
            <new-account-card @save="onNewAccountSave" @cancel="onNewAccountCancel"/>
          </v-menu>
        </v-btn>
      </template>
    </v-toolbar>
    <v-expand-transition>
      <div v-if="latestServiceCredential">
        <copy-secret-alert :credential="latestServiceCredential.secret" @close="onSecretClose"/>
      </div>
    </v-expand-transition>
    <v-card-text>
      <v-data-table-server
          v-bind="tableAttrs"
          :headers="tableHeaders"
          disable-sort
          return-object
          v-model="selectedAccounts">
        <template #item.type="{item, internalItem, isSelected, toggleSelect}">
          <div class="select--container" :class="{selected: isSelected(internalItem)}">
            <v-checkbox-btn :model-value="isSelected(internalItem)" @click="toggleSelect(internalItem)" color="primary"/>
            <v-icon :icon="accountTypeIcon(item.type)" v-tooltip:bottom="accountTypeStr(item.type)"/>
          </div>
        </template>
        <template #item.displayName="{item}">
          <span>{{ item.displayName }}</span>
          <span class="opacity-50 ml-2" v-if="item.description">{{ item.description }}</span>
        </template>
        <template #item.username="{item}">
          {{ item.type === Account.Type.USER_ACCOUNT ? item.username : item.id }}
        </template>
        <template #item.createTime="{item}">
          {{ timestampToDate(item.createTime).toLocaleDateString() }}
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {useDataTableCollection} from '@/composables/table.js';
import {
  accountToAddChange,
  accountToRemoveChange,
  accountTypeIcon,
  accountTypeStr,
  useAccountsCollection
} from '@/routes/auth/accounts/accounts.js';
import CopySecretAlert from '@/routes/auth/accounts/CopySecretAlert.vue';
import DeleteAccountsBtn from '@/routes/auth/accounts/DeleteAccountsBtn.vue';
import NewAccountCard from '@/routes/auth/accounts/NewAccountCard.vue';
import {Account} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed, ref, watch} from 'vue';

// used to fake PullAccounts when creating a new account.
// the var isn't reactive, but the value will be whenever it's set
let pullAccountsResource = /** @type {ResourceCollection<Account.AsObject, *> | null} */ null;

const wantCount = ref(20);
const accountsCollectionOpts = computed(() => {
  return {
    wantCount: wantCount.value,
    pullFn: (_, resource) => {
      pullAccountsResource = resource;
    }
  };
})
const accountsCollection = useAccountsCollection({}, accountsCollectionOpts);
const tableAttrs = useDataTableCollection(wantCount, accountsCollection);
const tableHeaders = computed(() => {
  return [
    {key: 'type', width: '1.5rem', cellProps: {class: 'pr-0'}},
    {title: 'Name', key: 'displayName', maxWidth: '10em', cellProps: {class: 'text-overflow-ellipsis'}},
    {title: 'Username / Client ID', key: 'username', maxWidth: '10em', cellProps: {class: 'text-overflow-ellipsis'}},
    {title: 'Created', key: 'createTime', label: 'Created'},
  ]
});

const latestAccount = ref(null);
const latestServiceCredential = ref(null);

const newAccountMenu = ref(false);

const onNewAccountSave = ({account, serviceCredential}) => {
  if (pullAccountsResource) {
    pullAccountsResource.lastResponse = accountToAddChange(account);
  }
  latestAccount.value = account;
  latestServiceCredential.value = serviceCredential;
  newAccountMenu.value = false;
};
const onNewAccountCancel = () => {
  newAccountMenu.value = false;
}
// reset the new account form once the menu is hidden
const newAccountCardRef = ref(null);
watch(newAccountMenu, (value) => {
  if (value) {
    setTimeout(() => {
      const form = newAccountCardRef.value;
      if (!form) return;
      form.reset();
    }, 250);
  }
});

const onSecretClose = () => {
  latestServiceCredential.value = null;
};

const selectedAccounts = ref([]);

const showDeleteAccountsBtn = computed(() => selectedAccounts.value.length > 0);
const onDelete = () => {
  if (pullAccountsResource) {
    for (const account of selectedAccounts.value) {
      pullAccountsResource.lastResponse = accountToRemoveChange(account);
    }
  }
  const _latest = latestAccount.value;
  for (const account of selectedAccounts.value) {
    if (account.id === _latest.id) {
      latestAccount.value = null;
      latestServiceCredential.value = null;
    }
  }
  selectedAccounts.value = [];
}
</script>

<style scoped lang="scss">
:deep(.v-toolbar__append) {
  gap: 1rem;
}

.select--container {
  display: grid;
  justify-items: center;
  align-items: center;

  .v-checkbox-btn {
    grid-column: 1 / -1;
    grid-row: 1 / -1;
  }

  .v-icon {
    grid-column: 1 / -1;
    grid-row: 1 / -1;
  }

  &:not(:hover, .selected) {
    .v-checkbox-btn {
      visibility: hidden;
    }
  }

  &:hover, &.selected {
    .v-icon {
      visibility: hidden;
    }
  }
}
</style>