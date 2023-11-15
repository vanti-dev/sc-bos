<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" two-line>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3"> Access Attempt </v-subheader>

      <v-col v-for="(val, key) in accessAttemptInfo[0]" :key="key" class="pa-0" cols="align-self">
        <v-list-item class="py-1">
          <v-list-item-content class="py-0 pb-3">
            <v-list-item-title class="text-body-small text-capitalize">
              {{ camelToSentence(key) }}
            </v-list-item-title>
            <v-list-item-subtitle class="text-subtitle-1 py-1 font-weight-medium text-wrap ml-2">
              {{ val }}
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-col>

      <v-col v-for="(subValue, subKey) in accessAttemptInfo[1]" :key="subKey" class="pa-0" cols="align-self">
        <v-list-item class="py-1">
          <v-list-item-content class="py-0 pb-3">
            <v-list-item-title class="text-body-small text-capitalize">
              {{ camelToSentence(subKey) }}
            </v-list-item-title>
            <v-row v-if="!Array.isArray(subValue)" class="py-0">
              <v-col v-for="(val, key) in subValue" :key="key" cols="align-self">
                <v-list-item-subtitle class="py-1 mx-2">
                  <v-col cols="align-self">
                    <v-row class="text-capitalize text-caption">
                      {{ camelToSentence(key) }}
                    </v-row>
                    <v-row v-if="!Array.isArray(val)" class="text-subtitle-1">
                      {{ val }}
                    </v-row>
                    <v-row v-else class="pa-0 ma-0 ml-n3">
                      <v-col v-for="(innerVal, innerKey) in val" :key="innerKey" cols="align-self">
                        <v-col cols="align-self">
                          <v-row class="text-capitalize text-caption">
                            {{ Object.keys(innerVal)[0] }}
                          </v-row>
                          <v-row class="text-subtitle-1">
                            {{ Object.values(innerVal)[0] }}
                          </v-row>
                        </v-col>
                      </v-col>
                    </v-row>
                  </v-col>
                </v-list-item-subtitle>
              </v-col>
            </v-row>
          </v-list-item-content>
        </v-list-item>
      </v-col>
    </v-list>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';
import {camelToSentence} from '@/util/string';
import {AccessAttempt} from '@sc-bos/ui-gen/proto/access_pb';

const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  },
  loading: {
    type: Boolean,
    default: false
  },
  showChangeDuration: {
    type: Number,
    default: 30 * 1000
  }
});

const grantNamesByID = Object.entries(AccessAttempt.Grant).reduce((all, [name, id]) => {
  all[id] = name.toLowerCase();
  return all;
}, {});

const accessAttemptInfo = computed(() => {
  // Initialize variables for info and subInfo
  const info = {};
  const subInfo = {};

  // Check if sidebarData has metadata property
  if (props?.value) {
    // Get all properties of metadata as an array of [key, value] pairs
    const data = Object.entries(props.value);

    // // Flatten out the data
    data.forEach(([key, value]) => {
      // If value is not empty
      if (value) {
        // If value is not an object
        if (typeof value !== 'object') {
          if (key === 'grant') {
            const state = grantNamesByID[value].split('_').join(' ');
            info[key] = state.charAt(0).toUpperCase() + state.slice(1);
          } else info[key] = value;
        } else {
          // Loop through the object for inner values
          for (const subValue in value) {
            if (subValue && value[subValue]) {
              // If subInfo[key] does not exist, create it
              if (!subInfo[key]) {
                subInfo[key] = {};
              }

              // If subValue is an array, map it to an object
              if (value[subValue].length) {
                if (Array.isArray(value[subValue])) {
                  // Map array to an object
                  const mappedArray = value[subValue].map(([key, value]) => {
                    return {[key]: value};
                  });

                  // Add subValue and its mapped value to subInfo[key]
                  subInfo[key][subValue] = mappedArray;
                } else {
                  // Add subValue and its value to subInfo[key]
                  subInfo[key][subValue] = value[subValue];
                }
              }
            }
          }
        }
      }
    });
  }

  // Return info and subInfo objects within an array
  return [info, subInfo];
});
</script>

<style scoped>
.granted {
  color: green;
}
.denied,
.forced,
.failed {
  color: red;
}
.pending,
.aborted,
.tailgate {
  color: orange;
}
.grant_unknown {
  color: grey;
}
</style>
