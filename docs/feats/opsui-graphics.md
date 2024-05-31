# Ops UI - Interactive Graphics

One of the widgets developed for the Ops UI is an interactive graphic.
This widget shows an SVG image (created elsewhere) and updates the styles and behaviour of that SVG based on data from
devices.

For example your SVG might show a floor in a building with light fixtures positioned on it.
The widget can be configured to adjust the colour of those light elements based on the brightness of different light
devices in Smart Core.

Clicking on interactive elements in the SVG will show the elements first source device in the sidebar.

## Getting Started

The widget is named `builtin:graphic/LayeredGraphic` and can be used on any dashboard, configured via the ui config
json file.
Here's an example snippet taken from [the UGS sample](../../example/config/vanti-ugs/ui-config.json):

```json
{
  "component": "builtin:graphic/LayeredGraphic",
  "props": {
    "height": "max(70vh, 700px)",
    "background": {"svgPath": "./assets/floor0-bg.svg"},
    "layers": [
      {
        "title": "BMS",
        "svgPath": "./assets/floor0-bms.svg",
        "configPath": "./assets/floor0-bms.json"
      },
      {
        "title": "ACS",
        "svgPath": "./assets/floor0-acs.svg",
        "configPath": "./assets/floor0-acs.json"
      }
    ]
  }
}
```

The SVG and the JSON are stored in separate files and loaded by the widget when that layer is displayed on the app.
SVG and JSON paths can be relative to the ui-config.json file or absolute.

The SVG is a standard SVG file, exported from any tool that can create SVG images.
The config file describes which elements in the SVG should be made interactive.

## Element Config File

Each element in the SVG that should be interactive is described in the JSON file and identified via a CSS selector.
Sources of data are configured per element using a trait and request payload.
Finally, the data is linked to the element via one or more effects that can be parameterised with the data.

Here's an example where an element representing a lighting fixture changes colour based on the Light trait brightness.
This example is slightly modified from the [UGS sample](../../example/config/vanti-ugs/assets/floor0-lights.json):

```json5
{
  "elements": [
    {
      // css selector matching the element in the SVG for this light
      "selector": "#LTF-L00-01",
      // sources of data for this element
      "sources": {
        "light": {
          "trait": "smartcore.traits.Light",
          "request": {"name": "van/uk/brum/ugs/devices/LTF-L00-01"}
        }
      },
      // effects to apply to the element based on the data
      "fill": {
        "source": {
          // "light" refers to the source defined above
          "ref": "light",
          // the property path to the data to use, supports foo.bar[1].baz syntax
          "property": "levelPercent"
        },
        // how to convert the data into a colour
        "interpolate": {
          "steps": [
            {"value": 0, "color": "#40464d"}, // grey when off
            {"value": 100, "color": "#ffaf25"} // yellow when on
          ]
        }
      }
    }
  ]
}
```

### Element Groups

Sometimes the device you are trying to represent graphically doesn't exist as a single element in the SVG.
An example of this are the lights in a room, there might be 4 fixtures in the ceiling and 4 elements in the SVG,
but they are only controllable as a single unit.

The CSS selector in the element config can match against multiple elements at once, in this case all elements will
become interactive and behave the same way. If configured to turn yellow when at `100`, then they will all turn yellow,
clicking on any of the elements will show the same device.

The lighting config in the UGS example uses this to configure the lights. The SVG has a group for each circuit with an
id attribute, then within those groups are the elements that represent the light fixtures.
In the config the `"selector"` is set to `"#LTF-L00-01 rect`, which matches all `<rect>` elements that are a child of
elements with `id="LTF-L00-01"`.

### Element Templates

As the number of interactive elements grows the config file quickly becomes unwieldy.
Generally the elements in the SVG will have a common structure, for example all lights might be `<rect>` elements with a
`<g id="...">` element.
To simplify the config, and reduce typos and errors, the config file allows you to define templates that can be reused
by elements to describe themselves.

A template is a JSON object that matches the structure of the elements in the `elements` array.
String keys and values can optionally have placeholders, using mustache syntax `{{foo}}`, in them that will later be
replaced when the template is used.

As an example we can convert the lighting example above to use templates and define more interactive elements in the
SVG.

```json5
{
  "templates": {
    "lightGroup": {
      // The {{id}} placeholder is filled by the "id" property in the "template" object of each element
      "selector": "#{{id}} rect",
      "sources": {
        "light": {
          "trait": "smartcore.traits.Light",
          // {{id}} here is the same property as above
          "request": {"name": "van/uk/brum/ugs/devices/{{id}}"}
        }
      },
      "fill": {
        "source": {"ref": "light", "property": "levelPercent"},
        "interpolate": {
          "steps": [
            {"value": 0, "color": "#40464d"},
            {"value": 100, "color": "#ffaf25"}
          ]
        }
      }
    }
  },
  "elements": [
    // If an element has a "template" property, that template will be loaded and
    // hydrated using the properties of that property.
    {"template": {"ref": "lightGroup", "id": "LTF-L00-01"}},
    // The template is merged with the element config, with the element config winning if both are present.
    {"selector": "#specialLight", "template": {"ref": "lightGroup", "id": "LTF-L00-02"}},
    {"template": {"ref": "lightGroup", "id": "LTF-L00-03"}},
    {"template": {"ref": "lightGroup", "id": "LTF-L00-04"}},
    {"template": {"ref": "lightGroup", "id": "LTF-L00-05"}}
  ]
}
```

Any number of templates can be defined to match the different structures and behaviours of the elements in the SVG.

### Effects

Effects describe how an SVG element should change appearance or behaviour based on the data from a source.
The effects are configured via properties on the element object, one property per distinct effect.

We've seen the `fill` effect in the previous examples, this effect adjusts the `style.fill` property of the SVG element
based on an interpolation between colour steps. Another similar effect is the `stroke` effect, which changes the stroke
colour.

All available effects are listed in [effects array](../../ui/conductor/src/widgets/graphic/svg.js) in the js source.

### Sidebar Interaction

By default, clicking on an element will show the app sidebar using the first source with a named request as the subject.
If this is not what you require from the element you can configure a specific subject name for the sidebar on element
click via the `sidebar` property:

```json5
{
  "elements": [
    {
      "selector": "#rooms #teal",
      // When clicked, show "van/uk/brum/ugs/zones/rooms/teal" in the sidebar
      // instead of inferring the subject from the sources
      "sidebar": {"name": "van/uk/brum/ugs/zones/rooms/teal"}
    }
  ]
}
```

The `sidebar` property is useful when the element only exists as a click target and doesn't have any sources or effects.
The usual example is when drawing zones or rooms on a map, clicking on the zone should open that zone in the sidebar.
You can see an example of this in the [UGS sample](../../example/config/vanti-ugs/assets/floor0-zones.json).