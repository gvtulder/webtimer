// @ts-check

import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import eslintConfigPrettier from "eslint-config-prettier";

export default tseslint.config({
  files: ["client/**/*.ts", "webpack.config.js"],
  extends: [
    eslint.configs.recommended,
    tseslint.configs.recommended,
    eslintConfigPrettier,
  ],
});
