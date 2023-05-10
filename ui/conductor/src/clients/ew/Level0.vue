<template>
  <v-card title flat class="pa-4 justify-center">
    <svg
        width="800px"
        height="100%"
        viewBox="0 0 45989 33818"
        xmlns="http://www.w3.org/2000/svg"
        style="fill-rule:evenodd;clip-rule:evenodd;stroke-linecap:round;stroke-miterlimit:1.5;">
      <g transform="matrix(1,0,0,1,-3182,-27027)">
        <!--eslint-disable max-len -->
        <path
            d="M3543.34,27388L36940.8,27388L36940.8,33301.2L48810,33301.2L48810,60484L9938,60484L9938,51447L3543.34,51447L3543.34,27388Z"
            style="fill:none;stroke:rgb(0,129,213);stroke-width:400px;"/>
        <g id="Doors" style="fill:#03BD5B" v-if="showDoors">
          <rect x="9615" y="54479" width="672" height="2921"/>
          <g transform="matrix(6.12323e-17,1,-0.637111,3.90118e-17,49608.1,17437)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(6.12323e-17,1,-0.637111,3.90118e-17,52202.1,17437)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(6.12323e-17,1,-0.637111,3.90118e-17,58681.1,17437)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(6.12323e-17,1,-0.637111,3.90118e-17,67544.1,17437)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(6.12323e-17,1,-0.415953,2.54698e-17,33816.7,17437)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(6.12323e-17,1,-0.415953,2.54698e-17,61069.7,23352)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(6.12323e-17,1,-0.415953,2.54698e-17,52088.7,50533)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(-1,-1.22465e-16,6.00794e-17,-0.490585,19902,86874.6)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
          <g transform="matrix(-1,-1.22465e-16,5.09396e-17,-0.415953,58761,63862.7)">
            <rect x="9615" y="54479" width="672" height="2921"/>
          </g>
        </g>
      </g>
    </svg>
  </v-card>
</template>

<script setup>

defineProps({
  showDoors: {
    type: Boolean,
    default: false
  }
});

</script>

<style scoped>

</style>
