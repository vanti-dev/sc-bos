import {getAccount, listAccounts} from '@/api/ui/account.js';
import {useAction} from '@/composables/action.js';
import useCollection from '@/composables/collection.js';
import {Account} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed, toValue} from 'vue';

/**
 * @param {import('vue').MaybeRefOrGetter<Partial<ListAccountsRequest.AsObject>>} request
 * @param {import('vue').MaybeRefOrGetter<Partial<ListOnlyCollectionOptions<ListAccountsRequest.AsObject, Account.AsObject>>>?} options
 * @return {UseCollectionResponse<Account.AsObject>}
 */
export function useAccountsCollection(request, options) {
  const normOpts = computed(() => {
    return {
      cmp: (a, b) => a.id.localeCompare(b.id, undefined, {numeric: true}),
      ...toValue(options)
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listAccounts(req, tracker);
      return {
        items: res.accountsList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      const opts = toValue(normOpts);
      if (opts.pullFn) {
        opts.pullFn(req, resource);
      }
    }
  };

  return useCollection(request, client, normOpts);
}

/**
 * @param {import('vue').MaybeRefOrGetter<Partial<GetAccountRequest.AsObject>>} request
 * @return {ToRefs<UnwrapNestedRefs<UseActionResponse<Account.AsObject>>>}
 */
export function useGetAccount(request) {
  return useAction(request, getAccount);
}

/**
 * @param {Account.Type} accountType
 * @return {string}
 */
export function accountTypeIcon(accountType) {
  switch (accountType) {
    case Account.Type.USER_ACCOUNT:
      return 'mdi-account';
    case Account.Type.SERVICE_ACCOUNT:
      return 'mdi-server';
    default:
      return 'mdi-account-question';
  }
}

/**
 * @param {Account.Type} accountType
 * @return {string}
 */
export function accountTypeStr(accountType) {
  switch (accountType) {
    case Account.Type.USER_ACCOUNT:
      return 'User Account';
    case Account.Type.SERVICE_ACCOUNT:
      return 'Service Account';
    default:
      return `Unknown Account Type ${accountType}`;
  }
}