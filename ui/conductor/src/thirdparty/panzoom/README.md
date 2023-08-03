# About

This file is built using the https://github.com/vanti-public/panzoom repository. We have forked the upstream repo to add
a feature to configure a minimum distance before pan/zoom is begun.

### Updating

To update the library

1. Clone the above repository
2. Checkout the `min-pan-distance` branch
3. Make your change
4. Run `npm run build` if you can (requires linux like systems, at least `rm`)
   1. Alternatively on windows, remove the `dist` folder then run `rollup --config`
5. Copy the `dist/panzoom.js` file into this project to overwrite the ./panzoom.js file
6. Test your changes


### Upstreaming

Ideally all changes will be integrated upstream into the base panzoom repository and released. The owner is fairly 
active if you raise issues or pull requests.

To replace our library fork with the official library

1. Update the panzoom dependency in the package.json and yarn install
2. Replace `import Panzoom from '@/thirdparty/panzoom/panzoom'` with `import Panzoom from '@panzoom/panzoom'` in 
   `src/components/TwPanZoom.vue`
3. Delete this folder

