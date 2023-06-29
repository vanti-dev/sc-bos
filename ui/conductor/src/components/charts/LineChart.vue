<template>
  <div style="height: 100%; max-height: 275px;">
    <div id="legend-container"/>
    <LineChartGenerator
        :options="props.chartOptions"
        :data="props.chartData"
        :plugins="[htmlLegendPlugin]"
        :dataset-id-key="props.datasetIdKey"
        :css-classes="props.cssClasses"
        :styles="props.styles"/>
  </div>
</template>

<script setup>
import {Line as LineChartGenerator} from 'vue-chartjs';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  Filler,
  LinearScale,
  CategoryScale,
  PointElement
} from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, LineElement, LinearScale, Filler, CategoryScale, PointElement);

const props = defineProps({
  datasetIdKey: {
    type: String,
    default: 'label'
  },
  cssClasses: {
    type: String,
    default: 'position-relative'
  },
  styles: {
    type: Object,
    default: () => {
      return {
        height: ''
      };
    }
  },
  chartData: {
    type: Object,
    default: () => {
      return {};
    }
  },
  chartOptions: {
    type: Object,
    default: () => {
      return {};
    }
  }
});

const getOrCreateLegendList = (id) => {
  const legendContainer = document.getElementById(id);
  let listContainer = legendContainer.querySelector('ul');

  if (!listContainer) {
    listContainer = document.createElement('ul');
    listContainer.style.display = 'flex';
    listContainer.style.flexDirection = 'row';
    listContainer.style.justifyContent = 'end';
    listContainer.style.margin = '-40px 0 20px 0';
    listContainer.style.padding = 0;

    legendContainer.appendChild(listContainer);
  }

  return listContainer;
};

const htmlLegendPlugin = {
  id: 'htmlLegend',
  afterUpdate(chart, args, options) {
    const ul = getOrCreateLegendList('legend-container');

    // Remove old legend items
    while (ul.firstChild) {
      ul.firstChild.remove();
    }

    // Reuse the built-in legendItems generator
    const items = chart.options.plugins.legend.labels.generateLabels(chart);

    items.forEach((item, index) => {
      // HTML Legend Item
      const li = document.createElement('li');
      li.id = 'legend-' + index;
      li.style.alignItems = 'center';
      li.style.cursor = 'pointer';
      li.style.display = 'flex';
      li.style.flexDirection = 'row';
      li.style.marginLeft = '10px';

      // Color box
      const boxSpan = document.createElement('span');
      boxSpan.style.background = item.strokeStyle;
      boxSpan.style.borderColor = item.strokeStyle;
      boxSpan.style.borderWidth = item.lineWidth + 'px';
      boxSpan.style.display = 'inline-block';
      boxSpan.style.height = '5px';
      boxSpan.style.marginRight = '15px';
      boxSpan.style.marginTop = '-15px';
      boxSpan.style.width = '15px';

      // Text
      const textContainer = document.createElement('p');
      textContainer.style.color = 'white';
      textContainer.style.marginRight = '15px';
      textContainer.style.padding = 0;
      textContainer.style.textDecoration = item.hidden ? 'line-through' : '';

      const text = document.createTextNode(item.text);
      textContainer.appendChild(text);

      li.onclick = () => {
        const {type} = chart.config;
        if (type === 'pie' || type === 'doughnut') {
          // Pie and doughnut charts only have a single dataset and visibility is per item
          chart.toggleDataVisibility(item.index);
        } else {
          chart.setDatasetVisibility(item.datasetIndex, !chart.isDatasetVisible(item.datasetIndex));
        }
        chart.update();
      };

      li.appendChild(boxSpan);
      li.appendChild(textContainer);
      ul.appendChild(li);
    });
  }
};
</script>

<style lang="scss">
</style>