import {computed, ref, toValue} from 'vue';
import * as colors from 'vuetify/util/colors';

const toChartDataset = (title, consumption) => {
  return {
    label: title,
    data: toValue(consumption).map((c) => c.y),
  }
}

/**
 * Creates a Chart.js dataset for each subValue and a remaining dataset for the total.
 *
 * @param {string} key - appended to titles to distinguish between different datasets, e.g. 'Consumption' or 'Production'
 * @param {import('vue').MaybeRefOrGetter<{x:Date, y:number|null}[]>} totals
 * @param {import('vue').MaybeRefOrGetter<(string|{name:string})[]>} subNames - specifies the order of the subValues
 * @param {import('vue').Reactive<Record<string, SubConsumption>>} subValues
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
    const dataset = toChartDataset(toValue(subValue.title), subValue.consumption);
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
/**
 * Helper to give type assistance to chart.js plugins.
 *
 * @template {import('chart.js').Plugin} T
 * @param {T} plugin
 * @return {T}
 */
export function definePlugin(plugin) {
  return plugin;
}

/**
 * Captures legend items using a Chart.js plugin.
 *
 * @return {{
 *   legendItems: Ref<{
 *     text: string,
 *     hidden: boolean,
 *     bgColor: string,
 *     onClick: (e: MouseEvent) => void
 *   }[]>,
 *   vueLegendPlugin: import('chart.js').Plugin
 * }}
 */
export function useVueLegendPlugin() {
  const legendItems = ref([]);
  return {
    legendItems,
    vueLegendPlugin: definePlugin({
      id: 'vueLegend',
      afterUpdate(chart) {
        const items = chart.options.plugins.legend.labels.generateLabels(chart);
        legendItems.value = items.map((item) => {
          return {
            text: item.text,
            hidden: item.hidden,
            bgColor: item.strokeStyle,
            onClick: (e) => {
              const {type} = chart.config;
              if (type === 'pie' || type === 'doughnut') {
                // Pie and doughnut charts only have a single dataset and visibility is per item
                chart.setDatasetVisibility(item.index, e);
              } else {
                chart.setDatasetVisibility(item.datasetIndex, e);
              }
              chart.update();
            }
          };
        });
      }
    })
  }
}

/**
 * Colours dataset based on the current theme.
 *
 * @return {{themeColorPlugin: import('chart.js').Plugin}}
 */
export function useThemeColorPlugin() {
  const datasetColors = computed(() => {
    return [
      colors.blue.base,
      colors.green.base,
      colors.orange.base,
      colors.yellow.base,
      colors.red.base,
    ].filter(Boolean);
  });
  return {
    themeColorPlugin: definePlugin({
      id: 'themeColor',
      beforeLayout(chart) {
        const colors = datasetColors.value;
        let i = 0;
        chart.data.datasets.forEach((dataset) => {
          if (dataset.backgroundColor && !dataset._pluginColor) return;
          const color = colors[i % colors.length];
          i++;
          dataset.backgroundColor = color + '80'; // 80 is 50% opacity
          dataset.borderColor = color;
          dataset._pluginColor = true;
        });
      }
    })
  }
}


/**
 * @typedef {Object} TooltipData
 * @property {number} x
 * @property {number} y
 * @property {number} opacity
 * @property {import('chart.js').TooltipItem[]} dataPoints
 * @property {Record<string,string>} displayFormats
 */

/**
 * @typedef {function} ExternalFunc
 * @this {import('chart.js').TooltipModel}
 * @property {{
 *   chart: import('chart.js').Chart,
 *   tooltip: import('chart.js').TooltipModel
 * }} args
 * @return {void}
 */

/**
 * Returns a function that can be used by chart.js tooltip.external and a ref containing the data.
 *
 * @return {{
 *   data: import('vue').Ref<TooltipData | null>,
 *   external: ExternalFunc
 * }}
 */
export function useExternalTooltip() {
  const data = ref(null);
  return {
    data,
    external: (ctx) => {
      if (!ctx.tooltip) {
        data.value = null;
        return;
      }
      const canvasBounds = ctx.chart.canvas.getBoundingClientRect();
      data.value = {
        x: ctx.tooltip.caretX + canvasBounds.left,
        y: ctx.tooltip.caretY + canvasBounds.top,
        opacity: ctx.tooltip.opacity,
        dataPoints: ctx.tooltip.dataPoints,
        displayFormats: ctx.chart.options.scales.x.time.displayFormats,
      };
    }
  }
}
