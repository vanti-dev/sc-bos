<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" two-line>
      <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Access Attempt</v-list-subheader>

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
import {useAccessAttempt} from '@/traits/access/access.js';
import {camelToSentence} from '@/util/string';

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

const {accessAttemptInfo} = useAccessAttempt(() => props.value);
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
