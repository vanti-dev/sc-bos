<template>
  <v-navigation-drawer
      v-model="drawer"
      absolute
      temporary
      right
      width="300px"
      text
      elevation="0">
    <v-list dense tile>
      <div class="d-flex justify-space-between align-center py-1 pl-2">
        <p class="mb-0">{{ selectedItem.device_id }}</p>
        <div>
          <v-btn text plain small> <v-icon>mdi-cog</v-icon></v-btn>
          <v-btn @click.stop="drawer = !drawer" text plain small>
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </div>
      </div>
      <v-divider/>
      <v-subheader>INFORMATION</v-subheader>
      <v-list-item v-for="item in lightData" :key="item.title" link>
        <v-list-item-content>
          <v-list-item-title class="text-caption">
            {{
              item.title
            }}
          </v-list-item-title>
        </v-list-item-content>
        <v-list-item-content>
          <v-list-item-title class="text-caption">
            {{
              item.content
            }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-divider/>

      <v-subheader>STATUS</v-subheader>

      <v-row column class="pa-2">
        <v-col class="px-3">
          <v-row>
            <v-col>
              <p class="mb-0 text-caption">Brightness</p>
            </v-col>
            <v-col>
              <p class="text-right mb-0 text-caption">
                {{ selectedItem.brightness }} %
              </p>
            </v-col>
          </v-row>
          <v-progress-linear
              :value="`${selectedItem.brightness}`"
              color="amber"
              height="15"/>
        </v-col>
        <v-col class="d-flex">
          <v-col class="d-flex">
            <v-btn @click="store.turnOn()" class="mr-2" color="green" small>
              ON
            </v-btn>
            <v-btn @click="store.turnOff()" color="red" small> OFF </v-btn>
          </v-col>
          <v-spacer/>
          <v-col class="d-flex">
            <v-btn
                @click="store.increaseBrightness()"
                class="mr-2"
                color="orange"
                small>
              UP
            </v-btn>
            <v-btn @click="store.decreaseBrightness()" color="orange" small>
              DOWN
            </v-btn>
          </v-col>
        </v-col>
      </v-row>
      <v-divider/>

      <v-subheader>EMERGENCY LIGHTING</v-subheader>

      <v-row column class="pa-2">
        <v-col class="px-3">
          <v-row>
            <v-col>
              <p class="mb-0 text-caption">Battery Level</p>
            </v-col>
            <v-col>
              <p class="text-right mb-0 text-caption">
                {{ selectedItem.battery_status }}%
              </p>
            </v-col>
          </v-row>

          <v-progress-linear
              :value="`${selectedItem.battery_status}`"
              color="amber"
              height="15"/>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-subheader id="test-header">Testing History</v-subheader>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0">
                28.09.22
              </v-list-item-title>
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title class="green--text text-right text-caption mb-0">
                PASS
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0">
                28.09.22
              </v-list-item-title>
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title class="green--text text-right text-caption mb-0">
                PASS
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0">
                28.09.22
              </v-list-item-title>
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title class="red--text text-right text-caption mb-0">
                FAIL
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0">
                28.09.22
              </v-list-item-title>
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title class="green--text text-right text-caption mb-0">
                PASS
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-action>
              <v-btn color="green" small>Test Now</v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-col>
      </v-row>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import {useLightingStore} from '@/stores/devices/lighting.js';
import {storeToRefs} from 'pinia';

const store = useLightingStore();

const {search, drawer, selectedItem, lightData} = storeToRefs(store);
</script>

<style scoped>
.v-list--dense .v-list-item,
.v-list-item--dense {
  min-height: 20px;
}
</style>
