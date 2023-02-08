<template>
  <v-container fluid class="pa-0 pt-2">
    <v-row>
      <v-col md="8" sm="12">
        <h3 class="text-h3 pt-2 pb-6">Status: Building Overview</h3>
      </v-col>
      <v-col>
        <content-card class="justify-center">
          <v-card-text class="py-1">
            <span class="text-title">Smart Core OS: </span>
            <span class="text-title-bold text-uppercase">Online</span>
          </v-card-text>
        </content-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col lg="8" md="12" sm="12">
        <content-card class="mb-6">
          <h4 class="text-h4">System Monitor</h4>
        </content-card>
        <content-card>
          <h4 class="text-h4">Energy</h4>
        </content-card>
      </v-col>
      <v-col>
        <content-card class="mb-6">
          <h4 class="text-h4">Occupancy</h4>
        </content-card>
        <content-card>
          <h4 class="text-h4">Environmental</h4>
          <circular-gauge width="300" :value="gaugeVal" min="-10" max="30"/>
          <v-slider v-model="gaugeVal" min="-10" max="30" step="0.5"/>
        </content-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script setup>
import ContentCard from '@/components/ContentCard.vue';
import CircularGauge from '@/components/CircularGauge.vue';
import {ref} from 'vue';

const gaugeVal = ref(0);

const floors = [
  {
    name: 'G',
    title: 'Ground Floor',
    zone: 'L00'
  },
  {
    name: '1',
    title: '1st Floor',
    zone: 'L01'
  },
  {
    name: '2',
    title: '2nd Floor',
    zone: 'L02'
  },
  {
    name: '3',
    title: '3rd Floor',
    zone: 'L03'
  },
  {
    name: '4',
    title: '4th Floor',
    zone: 'L04'
  },
  {
    name: '5',
    title: '5th Floor',
    zone: 'L05'
  },
  {
    name: '6',
    title: '6th Floor',
    zone: 'L06'
  },
  {
    name: '7',
    title: '7th Floor',
    zone: 'L07'
  },
  {
    name: '8',
    title: '8th Floor',
    zone: 'L08'
  },
  {
    name: '9',
    title: '9th Floor',
    zone: 'L09'
  },
  {
    name: '10',
    title: '10th Floor',
    zone: 'L10'
  }
];

</script>
