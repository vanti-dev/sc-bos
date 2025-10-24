import {nextTick} from 'vue';

/**
 * Iterates over large arrays asynchronously
 * @param array
 * @param cb
 * @param batched
 */
export function iterateLargeArray(array, cb, batched = false) {
  if (batched) {
    array.forEach(async batch => {
      batch.forEach(cb);
      await nextTick();
    });

    return;
  }

  array.forEach(async iter => {
    await nextTick(() => cb(iter));
  });
}

/**
 * Maps large arrays asynchronously
 * @param array
 * @param cb
 * @param batched
 * @return {Array}
 */
export function mapLargeArray(array, cb, batched = false) {
  if (batched) {
    const ret = [];
    array.map(async batch => {
      ret.push(...batch.map(cb));
      await nextTick();
    });
    return ret;
  }

  return array.map(async (item) => {
    return await nextTick(() => cb(item));
  });
}

export function batchLargeArray(array, batchSize = 50) {
  let batches = [];
  for (let i = 0; i < array.length; i += batchSize) {
    batches.push(array.slice(i, i + batchSize));
  }
  return batches;
}