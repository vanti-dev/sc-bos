<template>
  <dl class="val-obj md">
    <template v-for="(value, key) in props.metadata" :key="key">
      <!-- Check if the value is not empty or an empty array/object -->
      <template v-if="isValueAvailable(value)">
        <dt>
          {{ camelToSentence(key) }}:
        </dt>
        <dd class="md-prop-val">
          <ul v-if="isArray(value)" class="val-arr">
            <li v-for="(item, idx) in value" :key="idx">
              {{ item.name || item }}
            </li>
          </ul>
          <!-- Handle nested objects -->
          <dl v-else-if="isObject(value)" class="val-obj">
            <template v-for="(subValue, subKey) in value" :key="subKey">
              <template v-if="isValueAvailable(subValue)">
                <dt class="md-prop-key">
                  {{ camelToSentence(subKey) }}:
                </dt>
                <dd class="md-prop-val">
                  <!-- Handle simple values inside nested objects -->
                  <ul v-if="isArray(subValue)" class="val-arr">
                    <li v-for="(item, idx) in value" :key="idx">
                      {{ item.name || item }}
                    </li>
                  </ul>
                  <!-- Handle arrays and objects inside nested objects -->
                  <dl v-else-if="isObject(subValue)">
                    <template v-for="(deepValue, deepKey) in subValue" :key="deepKey">
                      <div v-if="isValueAvailable(deepValue)" class="md-prop-row">
                        <dt class="md-prop-key">{{ camelToSentence(deepKey) }}</dt>
                        <dd class="md-prop-val">
                          {{ isObject(deepValue) ? deepValue.name || JSON.stringify(deepValue) : deepValue }}
                        </dd>
                      </div>
                    </template>
                  </dl>
                  <template v-else>
                    <span class="val-plain">{{ subValue }}</span>
                  </template>
                </dd>
              </template>
            </template>
          </dl>
          <!-- Handle simple values and arrays with simple values -->
          <template v-else>
            <span class="val-plain">{{ value }}</span>
          </template>
        </dd>
      </template>
    </template>
  </dl>
</template>

<script setup>
import {camelToSentence} from '@/util/string';
import {isArray, isObject, isValueAvailable} from '@/util/types';

const props = defineProps({
  metadata: {
    type: Object,
    default: () => ({})
  }
});
</script>

<style scoped>
.md {
  --base-font-weight: 800;
}

dl {
  display: grid;
  grid-template-columns: minmax(10em, auto) 1fr;
  gap: .5em 1em;
  --child-font-weight: calc(var(--base-font-weight) * .6);
}

dt {
  text-transform: capitalize;
  font-weight: var(--base-font-weight);
}

dd {
  --base-font-weight: calc(var(--child-font-weight));
  font-weight: lighter;
}
</style>
