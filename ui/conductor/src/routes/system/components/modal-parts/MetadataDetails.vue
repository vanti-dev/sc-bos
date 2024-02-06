<template>
  <div class="px-6">
    <div class="pt-2 pb-2 px-6">
      <template v-for="(value, key) in props.metadata">
        <!-- Check if the value is not empty or an empty array/object -->
        <v-list-item
            v-if="isValueAvailable(value)"
            class="ma-0 pa-0 mb-n4"
            :key="key">
          <v-list-item-content class="d-flex flex-row flex-nowrap align-start">
            <v-col cols="align-self" class="ma-0 pa-0 mr-4">
              <v-list-item-title class="text-capitalize font-weight-bold ma-0 pa-0">
                {{ camelToSentence(key) }}:
              </v-list-item-title>
            </v-col>
            <v-col cols="10" class="ma-0 pa-0 pl-6 mb-2">
              <!-- Handle simple values and arrays with simple values -->
              <v-list-item-subtitle
                  v-if="!isObject(value) || isArray(value)"
                  class="ma-0 pa-0 text-wrap">
                <div v-if="isArray(value)" class="ml-0">
                  <div v-for="(item, idx) in value" class="pb-1" :key="idx">
                    {{ item.name || item }}
                  </div>
                </div>
                <span v-else class="">
                  {{ value }}
                </span>
              </v-list-item-subtitle>
              <!-- Handle nested objects -->
              <div v-else class="d-flex flex-column ml-n5">
                <template v-for="(subValue, subKey) in value">
                  <v-list-item
                      v-if="isValueAvailable(subValue)"
                      class="ma-0 pa-0 mt-n3 mb-n2"
                      :key="subKey">
                    <v-list-item-content
                        v-if="isValueAvailable(subValue)"
                        class="d-flex flex-row pb-4">
                      <v-col cols="3" class="ma-0 pa-0">
                        <v-list-item-title class="text-capitalize text-body-2 font-weight-medium ma-0 pa-0 ml-5">
                          {{ camelToSentence(subKey) }}:
                        </v-list-item-title>
                      </v-col>
                      <v-col cols="9" class="ma-0 pa-0">
                        <!-- Handle simple values inside nested objects -->
                        <v-list-item-subtitle v-if="!isObject(subValue)" class="ma-0 pa-0 text-wrap">
                          {{ subValue }}
                        </v-list-item-subtitle>
                        <!-- Handle arrays and objects inside nested objects -->
                        <div v-else>
                          <template v-for="(deepValue, deepKey) in subValue">
                            <v-list-item-subtitle
                                v-if="isValueAvailable(deepValue)"
                                class="ma-0 pa-0 text-wrap"
                                :key="deepKey">
                              {{ camelToSentence(deepKey) }}: {{
                                isObject(deepValue) ? deepValue.name || JSON.stringify(deepValue) : deepValue
                              }}
                            </v-list-item-subtitle>
                          </template>
                        </div>
                      </v-col>
                    </v-list-item-content>
                  </v-list-item>
                </template>
              </div>
            </v-col>
          </v-list-item-content>
        </v-list-item>
      </template>
    </div>
  </div>
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

</style>
