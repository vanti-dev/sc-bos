import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useLightingStore = defineStore("lightingStore", () => {
  const lights = ref([
    {
      device_id: "LIT-L02_12-001",
      location: "L02_12",
      status: "On",
      battery_status: "100",
      brightness: "100",
      model: "Philips LED 1245813",
    },
    {
      device_id: "LIT-L02_12-002",
      location: "L02_12",
      status: "Off",
      battery_status: "10",
      brightness: "100",
      model: "Philips LED 1245812",
    },
    {
      device_id: "LIT-L02_12-003",
      location: "L02_12",
      status: "On",
      battery_status: "0",
      brightness: "60",
      model: "Philips LED 1245814",
    },
    {
      device_id: "LIT-L02_12-004",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "100",
      model: "Philips LED 1245812",
    },
    {
      device_id: "LIT-L02_12-005",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "50",
      model: "Philips LED 1245812",
    },
    {
      device_id: "LIT-L02_12-006",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "100",
      model: "Philips LED 1245818",
    },
    {
      device_id: "LIT-L02_12-007",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "100",
      model: "Philips LED 1245819",
    },
    {
      device_id: "LIT-L02_12-008",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "100",
      model: "Philips LED 1245810",
    },
    {
      device_id: "LIT-L02_12-009",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "100",
      model: "Philips LED 1245812",
    },
    {
      device_id: "LIT-L02_12-010",
      location: "L02_12",
      status: "On",
      battery_status: "80",
      brightness: "100",
      model: "Philips LED 1245817",
    },
  ]);

  const headers = ref([
    {
      text: "Device ID",
      align: "start",
      sortable: false,
      value: "device_id",
    },
    { text: "Location", value: "location" },
    { text: "Status", value: "status" },
    { text: "Battery Status", value: "battery_status" },
    { text: "Model", value: "model" },
  ]);

  //an array for storing the selected lights for a particular operation
  const selected = ref([]);

  //a object for storing the selected light for a particular operation
  const selectedItem = ref({});

  const status = ref("All");

  const smartCoreStatus = ref("Online");

  const model = ref("All");

  const search = ref("");

  const drawer = ref(false);

  const items = ref([
    {
      title: "Building",
      content: "Upper Gough Street",
    },
    {
      title: "Floor",
      content: "LO1",
    },
    {
      title: "Zone",
      content: "L02_12",
    },
    {
      title: "Manufacturer",
      content: "Philips",
    },
    {
      title: "Model",
      content: "LED 1245812",
    },
    {
      title: "Installed on",
      content: "12.09.22",
    },
    {
      title: "Serial Number",
      content: "12348a7a595",
    },
    {
      title: "DALI Address",
      content: "1234",
    },
    {
      title: "DALI Controller",
      content: "1234",
    },
  ]);

  const models = ref([]);

  models.value = [...new Set(lights.value.map((light) => light.model))];
  models.value.unshift("All");

  const statuses = ref(["All", "On", "Off"]);

  const filteredLights = computed(() =>
    lights.value.filter((light) => {
      if (status.value === "All" && model.value === "All") {
        return true;
      } else if (status.value === "All") {
        return light.model === model.value;
      } else if (model.value === "All") {
        return light.status === status.value;
      } else {
        return light.status === status.value && light.model === model.value;
      }
    })
  );

  const bulkAction = (action) => {
    if (action === "On") {
      selected.value.forEach((light) => {
        light.status = "On";
      });
      selected.value = [];
    } else if (action === "Off") {
      selected.value.forEach((light) => {
        light.status = "Off";
      });
      selected.value = [];
    }
  };

  const toggleDrawer = () => {
    drawer.value = !drawer.value;
  };

  const setSelectedItem = (selected) => {
    selectedItem.value = selected;
  };

  const increaseBrightness = () => {
    if (selectedItem.value.brightness < 100) {
      selectedItem.value.brightness++;
    }
  };
  const decreaseBrightness = () => {
    if (selectedItem.value.brightness > 0) {
      selectedItem.value.brightness--;
    }
  };

  const turnOff = () => {
    selectedItem.value.status = "Off";
  };

  const turnOn = () => {
    selectedItem.value.status = "On";
  };

  const checkStatus = (status) => {};

  return {
    lights,
    headers,
    selected,
    items,
    status,
    smartCoreStatus,
    model,
    search,
    models,
    statuses,
    filteredLights,
    bulkAction,
    drawer,
    toggleDrawer,
    selectedItem,
    setSelectedItem,
    increaseBrightness,
    decreaseBrightness,
    turnOff,
    turnOn,
  };
});
