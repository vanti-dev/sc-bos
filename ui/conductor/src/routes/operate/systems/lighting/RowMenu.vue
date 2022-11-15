<template>
  <v-navigation-drawer
    v-model="drawer"
    absolute
    temporary
    right
    width="300px"
    class="bgColor"
    flat
    elevation="0"
  >
    <v-list dense tile>
      <div class="d-flex justify-space-between align-center pa-1 px-2">
        <p class="mb-0">LIT-L02_12-001</p>
        <div>
          <v-btn flat plain small> <v-icon>mdi-cog</v-icon></v-btn>
          <v-btn @click.stop="drawer = !drawer" flat plain small>
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </div>
      </div>
      <v-divider></v-divider>
      <v-subheader>INFORMATION</v-subheader>
      <v-list-item v-for="item in items" :key="item.title" link>
        <v-list-item-content>
          <v-list-item-title class="text-caption">{{
            item.title
          }}</v-list-item-title>
        </v-list-item-content>
        <v-list-item-content>
          <v-list-item-title class="text-caption">{{
            item.content
          }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-subheader>STATE</v-subheader>

      <v-layout row column class="pa-2">
        <v-flex class="px-3">
          <v-row>
            <v-col>
              <p class="mb-0 text-caption">Brightness</p>
            </v-col>
            <v-col>
              <p class="text-right mb-0 text-caption">60%</p>
            </v-col>
          </v-row>
          <v-progress-linear
            value="60"
            color="amber"
            height="15"
          ></v-progress-linear>
        </v-flex>
        <v-flex class="d-flex">
          <v-col class="d-flex">
            <v-btn class="mr-2" color="green" small> ON </v-btn>
            <v-btn color="red" small> OFF </v-btn>
          </v-col>
          <v-spacer></v-spacer>
          <v-col class="d-flex">
            <v-btn class="mr-2" color="orange" small> UP</v-btn>
            <v-btn color="orange" small> DOWN </v-btn>
          </v-col>
        </v-flex>
      </v-layout>
      <v-divider></v-divider>

      <v-subheader>EMERGENCY LIGHTING</v-subheader>

      <v-layout row column class="pa-2">
        <v-flex class="px-3">
          <v-row>
            <v-col>
              <p class="mb-0 text-caption">Battery Level</p>
            </v-col>
            <v-col>
              <p class="text-right mb-0 text-caption">60%</p>
            </v-col>
          </v-row>

          <v-progress-linear
            value="60"
            color="amber"
            height="15"
          ></v-progress-linear>
        </v-flex>
      </v-layout>
      <v-layout>
        <v-flex>
          <v-subheader id="test-header">Testing History</v-subheader>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0"
                >28.09.22</v-list-item-title
              >
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title
                class="green--text text-right text-caption mb-0"
                >PASS</v-list-item-title
              >
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0"
                >28.09.22</v-list-item-title
              >
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title
                class="green--text text-right text-caption mb-0"
                >PASS</v-list-item-title
              >
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0"
                >28.09.22</v-list-item-title
              >
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title class="red--text text-right text-caption mb-0"
                >FAIL</v-list-item-title
              >
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption mb-0"
                >28.09.22</v-list-item-title
              >
            </v-list-item-content>
            <v-list-item-content>
              <v-list-item-title
                class="green--text text-right text-caption mb-0"
                >PASS</v-list-item-title
              >
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-action>
              <v-btn color="green" small>Test Now</v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-flex>
      </v-layout>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();

const { search, drawer, selectedItem, items } = storeToRefs(store);
</script>

<style scoped>
.v-list--dense .v-list-item,
.v-list-item--dense {
  min-height: 20px;
}
</style>
