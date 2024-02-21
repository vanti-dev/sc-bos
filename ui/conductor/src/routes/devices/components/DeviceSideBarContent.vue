<template>
  <span>
    <device-info-card/>
    <WithStatus v-if="traits['smartcore.bos.Status']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <status-log-card v-bind="resource"/>
    </WithStatus>
    <WithAirTemperature v-if="traits['smartcore.traits.AirTemperature']" :name="deviceId" v-slot="{resource, update}">
      <v-divider class="mt-4 mb-1"/>
      <air-temperature-card v-bind="resource" @updateAirTemperature="update"/>
    </WithAirTemperature>
    <WithAirQuality v-if="traits['smartcore.traits.AirQualitySensor']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <air-quality-card v-bind="resource"/>
    </WithAirQuality>
    <v-divider v-if="traits['smartcore.traits.Light']" class="mt-4 mb-1"/>

    <LightCard v-if="traits['smartcore.traits.Light']" :name="deviceId"/>

    <v-divider v-if="traits['smartcore.traits.OccupancySensor']" class="mt-4 mb-1"/>
    <WithOccupancy v-if="traits['smartcore.traits.OccupancySensor']" :name="deviceId" v-slot="{resource}">
      <occupancy-card v-bind="resource"/>
    </WithOccupancy>
    <WithElectricDemand v-if="traits['smartcore.traits.Electric']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <electric-demand-card v-bind="resource"/>
    </WithElectricDemand>
    <WithMeter v-if="traits['smartcore.bos.Meter']" :name="deviceId" v-slot="{resource, info}">
      <v-divider class="mt-4 mb-1"/>
      <meter-card v-bind="resource" :unit="info?.response?.unit"/>
    </WithMeter>
    <v-divider v-if="traits['smartcore.bsp.EmergencyLight']" class="mt-4 mb-1"/>
    <emergency-light :name="deviceId" v-if="traits['smartcore.bsp.EmergencyLight']"/>
    <v-divider v-if="traits['smartcore.traits.Mode']" class="mt-4 mb-1"/>
    <mode-card :name="deviceId" v-if="traits['smartcore.traits.Mode']"/>
    <v-divider v-if="traits['smartcore.bos.UDMI']" class="mt-4 mb-1"/>
    <udmi-card :name="deviceId" v-if="traits['smartcore.bos.UDMI']"/>
  </span>
</template>

<script setup>
import WithAirQuality from '@/routes/devices/components/renderless/WithAirQuality.vue';
import WithAirTemperature from '@/routes/devices/components/renderless/WithAirTemperature.vue';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import WithMeter from '@/routes/devices/components/renderless/WithMeter.vue';
import WithOccupancy from '@/routes/devices/components/renderless/WithOccupancy.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import AirQualityCard from '@/traits/airQuality/AirQualityCard.vue';
import AirTemperatureCard from '@/traits/airTemperature/AirTemperatureCard.vue';
import DeviceInfoCard from '@/traits/deviceInfo/DeviceInfoCard.vue';
import ElectricDemandCard from '@/traits/electricDemand/ElectricDemandCard.vue';
import EmergencyLight from '@/traits/emergency/EmergencyLight.vue';
import LightCard from '@/traits/lighting/LightCard.vue';
import MeterCard from '@/traits/meter/MeterCard.vue';
import ModeCard from '@/traits/mode/ModeCard.vue';
import OccupancyCard from '@/traits/occupancy/OccupancyCard.vue';
import StatusLogCard from '@/traits/status/StatusLogCard.vue';
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
