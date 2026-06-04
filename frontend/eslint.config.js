// @ts-check
const eslint = require("@eslint/js");
const tseslint = require("typescript-eslint");
const angular = require("angular-eslint");

// Lenient baseline for this upstream-inherited codebase (foundation #12, Phase 0).
// tslint/codelyzer were dead (the tslint builder was removed in Angular 13), so the
// code had never been linted. Rather than rewrite ~2k findings at once, we:
//   - turn OFF opinionated / migration-nudge / legacy-style rules (no real bugs),
//   - downgrade valuable-but-voluminous rules to "warn" so they stay visible as debt,
//   - leave every other recommended rule at "error" to catch NEW problems.
// Ratchet rules back up to "error" incrementally as the debt is paid down.
module.exports = tseslint.config(
  {
    files: ["**/*.ts"],
    extends: [
      eslint.configs.recommended,
      ...tseslint.configs.recommended,
      ...tseslint.configs.stylistic,
      ...angular.configs.tsRecommended,
    ],
    processor: angular.processInlineTemplates,
    rules: {
      // -- off: opinionated / Angular-version migration nudges (constructor DI and
      //    NgModules are still valid; selectors are upstream-inconsistent) --
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-empty-function": "off",
      "@typescript-eslint/ban-ts-comment": "off",
      "@typescript-eslint/no-empty-object-type": "off",
      "@typescript-eslint/no-this-alias": "off",
      "@typescript-eslint/no-unsafe-function-type": "off",
      "@typescript-eslint/class-literal-property-style": "off",
      "@typescript-eslint/prefer-for-of": "off",
      "@angular-eslint/prefer-inject": "off",
      "@angular-eslint/prefer-standalone": "off",
      "@angular-eslint/component-selector": "off",
      "@angular-eslint/directive-selector": "off",
      "@angular-eslint/no-empty-lifecycle-method": "off",
      "@angular-eslint/no-output-on-prefix": "off",
      // -- warn: keep visible as debt, don't block --
      "@typescript-eslint/no-unused-vars": "warn",
      "@typescript-eslint/no-unused-expressions": "warn",
      "no-case-declarations": "warn",
      "no-unsafe-optional-chaining": "warn",
      "no-useless-escape": "warn",
      "no-empty": "warn",
      "prefer-const": "warn",
    },
  },
  {
    files: ["**/*.html"],
    extends: [
      ...angular.configs.templateRecommended,
      ...angular.configs.templateAccessibility,
    ],
    rules: {
      // -- warn: accessibility + template debt to pay down incrementally --
      "@angular-eslint/template/label-has-associated-control": "warn",
      "@angular-eslint/template/interactive-supports-focus": "warn",
      "@angular-eslint/template/click-events-have-key-events": "warn",
      "@angular-eslint/template/alt-text": "warn",
      "@angular-eslint/template/elements-content": "warn",
      "@angular-eslint/template/eqeqeq": "warn",
    },
  }
);
