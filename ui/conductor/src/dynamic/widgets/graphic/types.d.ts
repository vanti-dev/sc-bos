import {Component, Raw, Reactive} from "vue";

export namespace layer {
  export interface Element extends Record<string, unknown> {
    selector: string;
    sidebar?: {
      name: string;
    };
    sources?: Record<string, Source>;
    effects?: Array<Effect>;
  }

  export interface Source {
    trait: string;
    request: { name: string } & Record<string, unknown>
  }

  export interface SourceRef {
    ref: string;
    property?: string;
  }

  export interface Effect extends Record<string, unknown> {
    type: string;
  }

  export interface WidgetEffect extends Effect {
    type: "widget";
    component: string;
    // Props that are passed to the widget component instance.
    props?: Record<string, SourceRef | unknown>;
    // If the selected element is not the element that should be replaced by the widget, allow specifying a selector.
    // This is useful if the click target, or visual representation of the element is different from the element that should be replaced.
    selector?: string;
    // When true, the element the widget replaces will not be hidden.
    showElement?: boolean;
  }

  export interface WidgetInstance {
    component: Raw<Component>;
    key: string;
    // props is bound to the component instance via v-bind
    props: Reactive<Record<string, unknown>>;
    // bounds represents the style position of the widget instance.
    bounds: { top: string, left: string, width: string, height: string };
  }
}
