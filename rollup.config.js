import typescript from '@rollup/plugin-typescript';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import terser from '@rollup/plugin-terser';
import commonjs from '@rollup/plugin-commonjs';

export default {
  input: 'client/main.ts',
  output: {
    name: 'webtimer',
    dir: 'dist',
    format: 'iife',
    sourcemap: true,
    globals: {
      '@textfit': 'textfit',
    },
    sourcemapPathTransform: (relativeSourcePath, sourcemapPath) => {
      return `${relativeSourcePath.replace(new RegExp('^\.\./'), '')}`;
    },
  },
  plugins: [
    typescript({
      "compilerOptions": {
        "outDir": "dist/build-tsc"
      },
    }),
    commonjs(),
    nodeResolve(),
  ],
  external: ['@textfit']
};
