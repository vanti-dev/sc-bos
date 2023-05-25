import {defineStore} from 'pinia';
import {computed} from 'vue';

import {usePageStore} from '@/stores/page';

export const useTableHeaderStore = defineStore('tableHeader', () => {
  const pageStore = usePageStore();
  const {pageType} = pageStore;

  const activeTableHeader = computed(() => {
    let headers;

    if (pageType.automations || pageType.system) {
      headers = [
        {text: 'ID', value: 'id'},
        {text: 'Status', value: 'active'},
        {text: '', value: 'actions', align: 'end', width: '100'}
      ];
    } else if (pageType.devices || pageType.site) {
      headers = [
        {text: 'Device name', value: 'name'},
        {text: 'Floor', value: 'metadata.location.floor'},
        {text: 'Location', value: 'metadata.location.title'},
        !pageType.site ?? {
          text: '',
          value: 'hotpoints',
          align: 'end',
          width: '100',
          class: 'd-flex justify-end align-center'
        }
      ];
    } else if (pageType.ops) {
      headers = [
        {text: 'Timestamp', value: 'createTime', width: '15em'},
        {text: 'Floor', value: 'floor', width: '10em'},
        {text: 'Zone', value: 'zone', width: '10em'},
        {text: 'Severity', value: 'severity', width: '9em'},
        {text: 'Description', value: 'description', width: '100%'},
        {text: 'Acknowledged', value: 'acknowledged', align: 'center', width: '12em'}
      ];
    }

    return headers;
  });

  const headerCollection = computed(() => {
    const staticDataHeaders = activeTableHeader.value.filter((header) => header.text);
    const liveDataHeaders = activeTableHeader.value.filter((header) => !header.text);

    return {staticDataHeaders, liveDataHeaders};
  });

  return {
    activeTableHeader,
    headerCollection
  };
});
