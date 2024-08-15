<template>
  <v-card tile flat class="row py-4 justify-center">
    <div id="building" class="d-none d-md-block">
      <svg xmlns="http://www.w3.org/2000/svg" width="100%" height="350" fill="none">
        <mask id="a" width="350" height="349" x="57.5671" y=".5" fill="#000" maskUnits="userSpaceOnUse">
          <path fill="#fff" d="M57.5671.5h350v349h-350z"/>
          <path d="M62.5671 344.5V5.5H401.878v339H62.5671Z"/>
        </mask>

        <rect x="111.463" y="8" width="287.451" height="23" class="floor-status" :class="floorStatus['L10']"/>
        <rect x="111.463" y="38" width="287.451" height="22" class="floor-status" :class="floorStatus['L09']"/>
        <rect x="111.463" y="68" width="287.451" height="22" class="floor-status" :class="floorStatus['L08']"/>
        <rect x="65.0366" y="98" width="333.878" height="22" class="floor-status" :class="floorStatus['L07']"/>
        <rect x="65.0366" y="128" width="333.878" height="22" class="floor-status" :class="floorStatus['L06']"/>
        <rect x="65.0366" y="158" width="333.878" height="22" class="floor-status" :class="floorStatus['L05']"/>
        <rect x="65.0366" y="188" width="333.878" height="22" class="floor-status" :class="floorStatus['L04']"/>
        <rect x="65.0366" y="218" width="333.878" height="22" class="floor-status" :class="floorStatus['L03']"/>
        <rect x="65.0366" y="248" width="333.878" height="22" class="floor-status" :class="floorStatus['L02']"/>
        <rect x="106.524" y="278" width="292.39" height="22" class="floor-status" :class="floorStatus['L01']"/>
        <rect x="106.524" y="308" width="292.39" height="37" class="floor-status" :class="floorStatus['L00']"/>

        <!-- eslint-disable max-len -->
        <path fill="#00BED6" d="M62.5671 5.5v-5c-2.7614 0-5 2.23858-5 5h5Zm0 339h-5v5h5v-5Zm339.3109 0v5h5v-5h-5Zm0-339h5c0-2.76142-2.239-5-5-5v5Zm-344.3109 0v339h10V5.5h-10Zm5 344H401.878v-10H62.5671v10Zm344.3109-5V5.5h-10v339h10Zm-5-344H62.5671v10H401.878V.5Z" mask="url(#a)"/>
        <path stroke="#00BED6" stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M1.81708 348.5H442.378m-382.2804-75h43.4634m0 0H218.64v74H103.561v-74ZM60.0976 94h47.9084V2.5"/>
        <path fill="#793F17" d="M18.0976 347v1.5h3V347h-3Zm3-17.5c0-.828-.6716-1.5-1.5-1.5-.8285 0-1.5.672-1.5 1.5h3Zm0 17.5v-17.5h-3V347h3Zm13.7927 0v1.5h3V347h-3Zm3-15c0-.828-.6716-1.5-1.5-1.5-.8285 0-1.5.672-1.5 1.5h3Zm0 15v-15h-3v15h3Z"/>
        <path stroke="#00C27C" stroke-width="3" d="M36.3902 332.5c3.1981 0 5.828-2.178 7.5768-5.129 1.7682-2.983 2.8135-7.007 2.8135-11.371 0-4.364-1.0453-8.388-2.8135-11.371-1.7488-2.951-4.3787-5.129-7.5768-5.129-3.198 0-5.8279 2.178-7.5767 5.129C27.0453 307.612 26 311.636 26 316c0 4.364 1.0453 8.388 2.8135 11.371 1.7488 2.951 4.3787 5.129 7.5767 5.129Z"/>
        <path stroke="#C5CC3C" stroke-width="3" d="M19.5975 329.5c3.7092 0 6.8356-2.417 8.9577-5.834 2.136-3.441 3.4082-8.099 3.4082-13.166s-1.2722-9.725-3.4082-13.166c-2.1221-3.417-5.2485-5.834-8.9577-5.834-3.7091 0-6.8355 2.417-8.9576 5.834-2.13604 3.441-3.40821 8.099-3.40821 13.166s1.27217 9.725 3.40821 13.166c2.1221 3.417 5.2485 5.834 8.9576 5.834Z"/>
      <!-- eslint-enable max-len -->
      </svg>
    </div>
    <div id="status-grid">
      <span/>
      <span v-for="sys in systems" :key="sys.shortName" class="text-title system-name">{{ sys.shortName }}</span>
      <template v-for="floor in floors.slice().reverse()" :key="floor.name">
        <span class="floor-number text-title">{{ floor.name }}</span>
        <span
            v-for="sys in systems"
            :key="floor.name+sys.shortName"
            class="status-block"
            :class="systemStatus[sys.shortName.toLowerCase()]?.[floor.zone]"/>
      </template>
    </div>
  </v-card>
</template>

<script setup>

import {computed} from 'vue';

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
const systems = [
  {
    name: 'Fire',
    shortName: 'Fire'
  },
  {
    name: 'Lighting',
    shortName: 'Ltg'
  },
  {
    name: 'Power',
    shortName: 'Pwr'
  },
  {
    name: 'BMS',
    shortName: 'BMS'
  },
  {
    name: 'Toilets',
    shortName: 'WC'
  }
];

const systemStatus = computed(() => {
  const status = {};
  for (const s of systems) {
    status[s.shortName.toLowerCase()] = {};
    for (const f of floors) {
      status[s.shortName.toLowerCase()][f.zone] =
          ['online', 'online', 'online', 'online', 'online', 'online', 'online',
            'warning', 'error', 'offline'][Math.floor(Math.random()*10)];
    }
  }
  return status;
});

const systemStatusByFloor = computed(() => {
  const status = {};
  for (const f of floors) {
    status[f.zone] = {};
    for (const sys of systems) {
      status[f.zone][sys.shortName.toLowerCase()] = systemStatus.value[sys.shortName.toLowerCase()]?.[f.zone];
    }
  }
  return status;
});

const floorStatus = computed(() => {
  const status = {};
  for (const f of floors) {
    const s = Object.values(systemStatusByFloor.value[f.zone]);
    if (s.indexOf('error') >= 0 || s.indexOf('offline') >= 0) {
      status[f.zone] = 'error';
    } else if (s.indexOf('warning') >= 0) {
      status[f.zone] = 'warning';
    } else {
      status[f.zone] = 'online';
    }
  }
  return status;
});

</script>

<style lang="scss" scoped>
#building {
  max-width: 444px;
  margin-top: 25px;
  flex-grow: inherit;

  .floor-status {
    fill: var(--v-secondaryTeal-darken1);
    opacity: 0.5;
  }
  .online {
    fill: var(--v-success-lighten2);
  }
  .warning {
    fill: var(--v-warning-base);
  }
  .error {
    fill: var(--v-error-lighten1);
  }
}

#status-grid {
  display: grid;
  grid-template-columns: 0fr repeat(5, 1fr);
  grid-template-rows: repeat(12, 0fr);
  gap: 8px 10px;
  max-width: 235px;

  .floor-number, .system-name {
    display: block;
    text-align: center;
    align-self: center;
    text-transform: uppercase;
    min-width: 32px;
  }

  .status-block {
    display: block;
    height: 22px;
    background-color: var(--v-neutral-darken1);
    border: 1px solid var(--v-neutral-lighten4);
  }
  .online {
    background-color: var(--v-success-lighten2);
    border: none;
  }
  .warning {
    background-color: var(--v-warning-base);
  }
  .error {
    background-color: var(--v-error-lighten1);
  }
}

</style>
