import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource.js';
import {nextTick, reactive, set} from 'vue';

export class Collection {
  constructor(listFn, pullFn) {
    this._queryVersion = null;
    this.queryRaw = null;
    this.listFn = listFn;
    this.pullFn = pullFn;

    // _resources holds all the items we've fetched or updates via listFn or pullFn
    this._resources = reactive(newResourceCollection());
    // _fetchingPage tracks the fetching of the next page of results
    this._fetchingPage = reactive(newActionTracker());

    this._needsMorePages = false;
    this._nextPageToken = null;
  }

  query(q = undefined) {
    this.reset();

    this.queryRaw = q;
    this._queryVersion = Math.random();
    this.fetchPages()
        .catch(err => console.error(err));
  }

  pullIfNeeded() {
    if (this.pullFn && !this._resources.loading && !this._resources.stream) {
      this.pullFn(this.queryRaw, this._resources);
    }
  }

  async fetchPages() {
    if (!this._needsMorePages) return;
    this.pullIfNeeded();
    const queryVersion = this._queryVersion;
    this._nextPageToken = await this.listFn(this.queryRaw, this._fetchingPage, this._nextPageToken, (item, id) => {
      if (queryVersion === this._queryVersion) {
        set(this._resources.value, id, item);
      }
    });

    if (!this._nextPageToken) return; // the server has no more pages

    // give the ui a chance to update, then check if the ui wants us to fetch more pages.
    await nextTick();
    await this.fetchPages();
  }

  set needsMorePages(b) {
    if (b === this._needsMorePages) return;

    this._needsMorePages = b;
    this.fetchPages()
        .catch(err => console.error(err));
  }

  get needsMorePages() {
    return this._needsMorePages;
  }

  reset() {
    closeResource(this._resources);
    this._resources.value = {};
    this._nextPageToken = null;
  }

  get nextPageResource() {
    return this._fetchingPage;
  }

  get resources() {
    return this._resources;
  }
}
