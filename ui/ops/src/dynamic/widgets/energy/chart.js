import {toValue} from 'vue';

const toChartDataset = (title, data) => {
  return {
    label: title,
    data: toValue(data).map((c) => c.y),
  }
}

export const datasetSourceName = Symbol('datasetSourceName');

/**
 * Creates a Chart.js dataset for each subValue and a remaining dataset for the total.
 *
 * @param {string} key - appended to titles to distinguish between different datasets, e.g. 'Consumption' or 'Production'
 * @param {import('vue').MaybeRefOrGetter<{x:Date, y:number|null}[]>} totals
 * @param {import('vue').MaybeRefOrGetter<(string|{name:string})[]>} subNames - specifies the order of the subValues
 * @param {import('vue').Reactive<Record<string, Series>>} subValues
 * @param {boolean} invert - if true, the dataset will be placed below the x-axis
 * @return {import('chart.js').ChartDataSets[]}
 */
export function computeDatasets(key, totals, subNames, subValues, invert = false) {
  const _subNames = toValue(subNames) ?? [];

  const datasets = [];
  const remaining = toChartDataset(_subNames.length === 0 ? `Total ${key}` : `Other ${key}`, totals);
  if (_subNames.length > 0) {
    // muted colours for the remaining dataset
    remaining.backgroundColor = '#cccccc80';
    remaining.borderColor = '#cccccc';
  }
  for (const name of _subNames) {
    const subValue = subValues[name];
    if (!subValue) continue;
    const dataset = toChartDataset(toValue(subValue.title), subValue.data);
    dataset[datasetSourceName] = name;
    let hasAny = false;
    for (let i = 0; i < dataset.data.length; i++) {
      if (dataset.data[i] !== null) {
        hasAny = true;
        if (remaining.data[i] !== null) {
          remaining.data[i] -= dataset.data[i];
        }
      }
    }
    if (!hasAny) continue;
    datasets.push(dataset);
  }

  // get rid of negative remaining values
  remaining.data = remaining.data.map((v) => v <= 0 ? null : v);
  if (remaining.data.some((v) => v !== null)) {
    datasets.push(remaining);
  }

  if (invert) {
    // make all values negative so the bars appear below the x-axis
    for (const dataset of datasets) {
      dataset.data = dataset.data.map((v) => v === null ? null : -v);
      dataset._inverted = true;
    }
  }

  return datasets;
}
