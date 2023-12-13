import {toValue} from '@/util/vue';
import {computed, onMounted, onUnmounted, reactive, ref, watch, watchEffect} from 'vue';

import {newActionTracker} from '@/api/resource';
import {startService, stopService} from '@/api/ui/services';

import {useHubStore} from '@/stores/hub';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';

import {useErrorStore} from '@/components/ui-error/error';
import {serviceName} from '@/util/proxy';

/**
 * @param {{
 * name: MaybeRefOrGetter<string>,
 * type: MaybeRefOrGetter<string>
 * }} props
 * @return {{
 * serviceList: ComputedRef<Service.AsObject[]>,
 * serviceCollection: Ref<Collection<Service.AsObject>>,
 * showService: function(Service.AsObject): void,
 * startStopTracker: ActionTracker,
 * search: Ref<string>,
 * node: ComputedRef<Node.AsObject>,
 * nodesListValues: ComputedRef<Node.AsObject[]>,
 * _startService: function(Service.AsObject): Promise<void>,
 * _stopService: function(Service.AsObject): Promise<void>
 * }}
 */
export default function(props) {
  const serviceStore = useServicesStore();
  const pageStore = usePageStore();
  const errors = useErrorStore();
  const hubStore = useHubStore();

  const startStopTracker = reactive(
      /** @type {ActionTracker<Service.AsObject>} */
      newActionTracker()
  );
  const serviceCollection = ref({});
  const search = ref('');

  // The node to query services from
  const node = computed({
    get() {
      return pageStore.sidebarNode;
    },
    set(val) {
      pageStore.sidebarNode = val;
    }
  });

  // Available services to select from
  const serviceList = computed(() => {
    return Object.values(serviceCollection.value?.resources?.value ?? []).filter(service => {
      const type = toValue(props.type);
      return type === '' || type === 'all' || service.type === type;
    });
  });

  // Available nodes to select from
  const nodesListValues = computed(() => Object.values(hubStore.nodesList));


  // Watch for changes in pageStore.sidebarNode.name and update it if needed
  watchEffect(() => {
    if (!pageStore.sidebarNode.name) {
      if (nodesListValues.value.length > 0) {
        pageStore.sidebarNode = nodesListValues.value[0];
      }
    }
  });

  // Watch for changes in the name prop and the name of the node and update the serviceCollection
  watch(
      [() => props.name, node],
      async ([newName, newNode]) => {
        if (serviceCollection.value.reset) serviceCollection.value.reset();

        serviceCollection.value =
          serviceStore.getService(
              newName, await newNode?.commsAddress, await newNode?.commsName
          ).servicesCollection;

        // reinitialize in case this service collection has been previously reset;
        serviceCollection.value.init();
        serviceCollection.value.query(newName);
      },
      {immediate: true}
  );

  watch(serviceCollection, () => {
    // todo: this causes us to load all pages, connect with paging logic instead
    serviceCollection.value.needsMorePages = true;
  });

  //
  //
  // SERVICE ACTIONS
  /**
   *
   * @param {Service.AsObject} service
   */
  function showService(service) {
    pageStore.showSidebar = true;
    pageStore.sidebarTitle = service.id;
    pageStore.sidebarData = {...service, config: JSON.parse(service.configRaw)};
  }

  /**
   *
   * @param {Service.AsObject} service
   */
  async function _startService(service) {
    // Update the sidebarData if the sidebar is open and the service is being started
    if (pageStore.showSidebar && pageStore.sidebarData.id !== service.id) {
      pageStore.sidebarTitle = service.id;
      pageStore.sidebarData = {...service, config: JSON.parse(service.configRaw)};
    }

    await startService({
      name: serviceName(await node.value.commsName, toValue(props.name)),
      id: service.id
    }, startStopTracker);
  }


  /**
   *
   * @param {Service.AsObject} service
   */
  async function _stopService(service) {
    // Update the sidebarData if the sidebar is open and the service is being stopped
    if (pageStore.showSidebar && pageStore.sidebarData.id !== service.id) {
      pageStore.sidebarTitle = service.id;
      pageStore.sidebarData = {...service, config: JSON.parse(service.configRaw)};
    }

    await stopService({
      name: serviceName(await node.value.commsName, toValue(props.name)),
      id: service.id
    }, startStopTracker);
  }

  // Watch for changes in the serviceList and update the sidebarData if needed
  // This is necessary if we want to update the status details in the sidebar
  // simultaneously with the status details in the service list
  // Mainly when the sidebar is open and the service is being started/stopped
  watch(
      serviceList,
      (newServiceList, oldServiceList) => {
        if (pageStore.sidebarData === null || !pageStore.sidebarData.id) return;

        // Find the service in the new list that matches the id in sidebarData
        const updatedService = newServiceList.find(s => s.id === pageStore.sidebarData.id);

        if (updatedService) {
        // Check if the service has been updated by comparing it with the old list
          const oldService = oldServiceList.find(s => s.id === updatedService.id);

          // Perform a deep comparison if necessary, for now, we just check if the old service exists
          if (!oldService || JSON.stringify(updatedService) !== JSON.stringify(oldService)) {
          // Update the sidebarData with the new service data
          // Ensuring to parse the config if it's in a raw JSON string format
            pageStore.sidebarData = {
              ...updatedService,
              config: updatedService.configRaw ? JSON.parse(updatedService.configRaw) : {}
            };
          }
        }
      },
      {immediate: true, deep: true} // Watch for nested changes within the serviceList array
  );

  //
  //
  // UI error handling
  let unwatchErrors; let unwatchStartStopErrors;
  onMounted(() => {
    unwatchErrors = errors.registerCollection(serviceCollection);
    unwatchStartStopErrors = errors.registerTracker(startStopTracker);
  });
  onUnmounted(() => {
    if (unwatchErrors) unwatchErrors();
    if (unwatchStartStopErrors) unwatchStartStopErrors();
    serviceCollection.value.reset();
  });

  return {
    serviceList,
    serviceCollection,
    showService,
    startStopTracker,
    search,
    node,
    nodesListValues,
    _startService,
    _stopService
  };
}
