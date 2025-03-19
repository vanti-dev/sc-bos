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
        <grant-role-btn v-if="showGrantBtn" variant="outlined" :accounts="allSelectedAccountIds"/>
        <new-account-btn variant="flat" color="primary" @save="onNewAccountSave"/>
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
          v-model="selectedAccounts"
          :row-props="tableRowProps"
          @click:row="onRowClick">
        <template #top>
          <v-expand-transition>
            <div v-if="tableErrorStr">
              <v-alert type="error" :text="tableErrorStr"/>
            </div>
          </v-expand-transition>
        </template>
        <template #item.type="{item, internalItem, isSelected, toggleSelect}">
          <div class="select--container" :class="{selected: isSelected(internalItem)}">
            <v-checkbox-btn :model-value="isSelected(internalItem)"
                            @click="toggleSelect(internalItem)"
                            color="primary"/>
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
import {updateAccount} from '@/api/ui/account.js';
import {useDataTableCollection} from '@/composables/table.js';
import {toAddChange, toRemoveChange, toUpdateChange} from '@/routes/auth/accounts.js';
import {
  accountTypeIcon,
  accountTypeStr,
  useAccountsCollection,
  useGetAccount
} from '@/routes/auth/accounts/accounts.js';
import CopySecretAlert from '@/routes/auth/accounts/CopySecretAlert.vue';
import DeleteAccountsBtn from '@/routes/auth/accounts/DeleteAccountsBtn.vue';
import GrantRoleBtn from '@/routes/auth/accounts/GrantRoleBtn.vue';
import NewAccountBtn from '@/routes/auth/accounts/NewAccountBtn.vue';
import {useSidebarStore} from '@/stores/sidebar.js';
import {Account} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed, ref, watch} from 'vue';
import {useRouter} from 'vue-router';

const props = defineProps({
  accountId: {
    type: String,
    default: null,
  }
});

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

const tableRowProps = ({item}) => {
  return {
    class: {
      'row-selected': selectedAccounts.value.includes(item),
      'row-active': item.id === props.accountId,
    }
  };
};
const tableErrorStr = computed(() => {
  const errors = accountsCollection.errors.value;
  if (errors.length === 0) return null;
  return 'Error fetching accounts: ' + errors.map((e) => (e.error ?? e).message ?? e).join(', ');
});

const router = useRouter();
const onRowClick = (_, {item}) => {
  if (item.id === props.accountId) return; // don't click on the same item
  router.push({name: 'accounts', params: {accountId: item.id}});
}
const sidebar = useSidebarStore();
const {response: sidebarItem, refresh: refreshSidebarItem} = useGetAccount(() => {
  if (!props.accountId) return null;
  return {id: props.accountId};
});
watch(sidebarItem, (item) => {
  if (!item) {
    sidebar.closeSidebar();
    return;
  }
  sidebar.title = item.displayName || `Account ${props.roleId}`;
  sidebar.data = {account: item, updateAccount: onAccountUpdate};
  sidebar.visible = true;
}, {immediate: true});

const latestAccount = ref(null);
const latestServiceCredential = ref(null);

const onNewAccountSave = ({account, serviceCredential}) => {
  if (pullAccountsResource) {
    pullAccountsResource.lastResponse = toAddChange(account);
  }
  latestAccount.value = account;
  latestServiceCredential.value = serviceCredential;
};
const onAccountUpdate = async ({account}) => {
  const oldAccount = (() => {
    if (sidebarItem.value?.id === account.id) return sidebarItem.value;
    return accountsCollection.items.value.find((r) => r.id === account.id);
  })();

  const newAccount = await updateAccount({account})

  if (!oldAccount) {
    // we can't do a dynamic of the account, so just refresh the whole list
    accountsCollection.refresh();
  } else {
    pullAccountsResource.lastResponse = toUpdateChange(oldAccount, newAccount);
  }
  if (latestAccount.value?.id === account.id) {
    latestAccount.value = newAccount;
  }
  selectedAccounts.value = selectedAccounts.value.map((r) => {
    if (r.id === account.id) return newAccount;
    return r;
  });
  refreshSidebarItem();

  return newAccount;
}

const onSecretClose = () => {
  latestServiceCredential.value = null;
};

const selectedAccounts = ref([]);
const allSelectedAccountIds = computed(() => {
  const ids = {};
  for (const account of selectedAccounts.value) {
    ids[account.id] = true;
  }
  if (props.accountId) {
    ids[props.accountId] = true;
  }
  return Object.keys(ids).sort();
})

const showGrantBtn = computed(() => allSelectedAccountIds.value.length > 0);

const showDeleteAccountsBtn = computed(() => selectedAccounts.value.length > 0);
const onDelete = () => {
  if (pullAccountsResource) {
    for (const account of selectedAccounts.value) {
      pullAccountsResource.lastResponse = toRemoveChange(account);
    }
  }
  const _latest = latestAccount.value;
  for (const account of selectedAccounts.value) {
    if (_latest && account.id === _latest.id) {
      latestAccount.value = null;
      latestServiceCredential.value = null;
    }
    if (account.id === props.accountId) {
      router.push({name: 'accounts'});
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

.v-data-table {
  :deep(.row-selected) {
    background-color: rgba(var(--v-theme-primary), 0.1);
  }

  :deep(.row-active) {
    background-color: rgba(var(--v-theme-primary), 0.4);
  }
}
</style>