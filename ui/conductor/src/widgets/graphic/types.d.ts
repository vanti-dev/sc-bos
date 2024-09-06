import {Component, Raw, Reactive} from "vue";

export namespace layer {
  export interface Element extends Record<string, any> {
    selector: string;
    sidebar?: {
      name: string;
    };
    sources?: Record<string, Source>;
    effects?: Array<Effect>;
  }

  export interface Source {
    trait: string;
    request: { name: string } & Record<string, any>
  }

  export interface Effect extends Record<string, any> {
    type: string;
  }

  export interface WidgetEffect extends Effect {
    type: "widget";
    component: string;
    props: Record<string, any>;
  }

  export interface WidgetInstance {
    component: Raw<Component>;
    key: string;
    // props is bound to the component instance via v-bind
    props: Reactive<Record<string, any>>;
    bounds: { top: string, left: string, width: string, height: string };
  }
}
