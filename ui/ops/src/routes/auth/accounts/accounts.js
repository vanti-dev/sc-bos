import {listAccounts} from '@/api/ui/account.js';
import useCollection from '@/composables/collection.js';
import {ChangeType} from '@smart-core-os/sc-api-grpc-web/types/change_pb';
import {Account} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed, toValue} from 'vue';

/**
 * @typedef {UseCollectionOptions<T>} ListOnlyCollectionOptions
 * @property {function(R, ResourceCollection<T, any>): void} pullFn
 * @template R
 * @template T
 */

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
 * Returns an object that looks like a Pull Change that adds the given val.
 * Useful when an API doesn't support Pull but you want to reuse utilities that do.
 *
 * @param {T} val
 * @return {Object}
 * @template T
 */
export function toAddChange(val) {
  const changes = {
    changesList: [{
      type: ChangeType.ADD,
      newValue: val,
    }]
  };
  return {
    toObject() {
      return changes;
    },
    getChangesList() {
      return changes.changesList.map(change => ({toObject() { return change; } }));
    }
  }
}

/**
 * Returns an object that looks like a Pull Change that deletes the given val.
 * Useful when an API doesn't support Pull but you want to reuse
 *
 * @param {T} val
 * @return {Object}
 * @template T
 */
export function toRemoveChange(val) {
  const changes = {
    changesList: [{
      type: ChangeType.REMOVE,
      oldValue: val,
    }]
  };
  return {
    toObject() {
      return changes;
    },
    getChangesList() {
      return changes.changesList.map(change => ({toObject() { return change; } }));
    }
  }
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