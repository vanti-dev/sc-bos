<template>
  <span>
    <metadata-card/>
    <with-status v-if="traits['smartcore.bos.Status']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <status-log-card v-bind="resource"/>
    </with-status>
    <with-air-temperature v-if="traits['smartcore.traits.AirTemperature']" :name="deviceId" v-slot="{resource, update}">
      <v-divider class="mt-4 mb-1"/>
      <air-temperature-card v-bind="resource" @updateAirTemperature="update"/>
    </with-air-temperature>
    <with-air-quality v-if="traits['smartcore.traits.AirQualitySensor']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <air-quality-card v-bind="resource"/>
      <air-quality-history-card v-bind="resource" :name="deviceId"/>
    </with-air-quality>
    <v-divider v-if="traits['smartcore.traits.Light']" class="mt-4 mb-1"/>

    <light-card v-if="traits['smartcore.traits.Light']" :name="deviceId"/>

    <v-divider v-if="traits['smartcore.traits.OccupancySensor']" class="mt-4 mb-1"/>
    <with-occupancy v-if="traits['smartcore.traits.OccupancySensor']" :name="deviceId" v-slot="{resource}">
      <occupancy-card v-bind="resource"/>
      <occupancy-history-card v-bind="resource" :name="deviceId"/>
    </with-occupancy>
    <with-electric-demand v-if="traits['smartcore.traits.Electric']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <electric-demand-card v-bind="resource"/>
    </with-electric-demand>
    <with-meter v-if="traits['smartcore.bos.Meter']" :name="deviceId" v-slot="{resource, info}">
      <v-divider class="mt-4 mb-1"/>
      <meter-card v-bind="resource" :info="info?.response" :name="deviceId"/>
      <!-- todo we really only want to display the history card if this device has any history -->
      <meter-history-card v-bind="resource" :name="deviceId"/>
    </with-meter>
    <v-divider v-if="traits['smartcore.bsp.EmergencyLight']" class="mt-4 mb-1"/>
    <emergency-light :name="deviceId" v-if="traits['smartcore.bsp.EmergencyLight']"/>
    <v-divider v-if="traits['smartcore.traits.Mode']" class="mt-4 mb-1"/>
    <mode-card :name="deviceId" v-if="traits['smartcore.traits.Mode']"/>
    <v-divider v-if="traits['smartcore.bos.UDMI']" class="mt-4 mb-1"/>
    <udmi-card :name="deviceId" v-if="traits['smartcore.bos.UDMI']"/>
  </span>
</template>

<script setup>
import AirQualityCard from '@/traits/airQuality/AirQualityCard.vue';
import AirQualityHistoryCard from '@/traits/airQuality/AirQualityHistoryCard.vue';
import WithAirQuality from '@/traits/airQuality/WithAirQuality.vue';
import AirTemperatureCard from '@/traits/airTemperature/AirTemperatureCard.vue';
import WithAirTemperature from '@/traits/airTemperature/WithAirTemperature.vue';
import ElectricDemandCard from '@/traits/electricDemand/ElectricDemandCard.vue';
import WithElectricDemand from '@/traits/electricDemand/WithElectricDemand.vue';
import EmergencyLight from '@/traits/emergency/EmergencyLight.vue';
import LightCard from '@/traits/light/LightCard.vue';
import MetadataCard from '@/traits/metadata/MetadataCard.vue';
import MeterCard from '@/traits/meter/MeterCard.vue';
import MeterHistoryCard from '@/traits/meter/MeterHistoryCard.vue';
import WithMeter from '@/traits/meter/WithMeter.vue';
import ModeCard from '@/traits/mode/ModeCard.vue';
import OccupancyCard from '@/traits/occupancy/OccupancyCard.vue';
import OccupancyHistoryCard from '@/traits/occupancy/OccupancyHistoryCard.vue';
import WithOccupancy from '@/traits/occupancy/WithOccupancy.vue';
import StatusLogCard from '@/traits/status/StatusLogCard.vue';
import WithStatus from '@/traits/status/WithStatus.vue';
import UdmiCard from '@/traits/udmi/UdmiCard.vue';

defineProps({
  deviceId: {
    type: String,
    default: ''
  },
  traits: {
    type: Object,
    default: () => {
    }
  }
});
</script>
