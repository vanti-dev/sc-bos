import {ChangeType} from '@smart-core-os/sc-api-grpc-web/types/change_pb';

/**
 * @typedef {UseCollectionOptions<T>} ListOnlyCollectionOptions
 * @property {function(R, ResourceCollection<T, any>): void} pullFn
 * @template R
 * @template T
 */

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
      return changes.changesList.map(change => ({toObject() { return change; }}));
    }
  }
}

/**
 * Returns an object that looks like a Pull Change that deletes the given val.
 * Useful when an API doesn't support Pull but you want to reuse utilities that do.
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
      return changes.changesList.map(change => ({toObject() { return change; }}));
    }
  }
}

/**
 * Returns an object that looks like a Pull Change that updates the val.
 * Useful when an API doesn't support Pull but you want to reuse utilities that do.
 *
 * @param {T} oldVal
 * @param {T} newVal
 * @return {Object}
 * @template T
 */
export function toUpdateChange(oldVal, newVal) {
  const changes = {
    changesList: [{
      type: ChangeType.UPDATE,
      oldValue: oldVal,
      newValue: newVal,
    }]
  };
  return {
    toObject() {
      return changes;
    },
    getChangesList() {
      return changes.changesList.map(change => ({toObject() { return change; }}));
    }
  }
}
