import type { StorybookConfig } from "@storybook/angular";
const config: StorybookConfig = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|ts|tsx)"],
  staticDirs: [{ from: '../src/assets', to: '/assets' }],
  addons: [
    // SB9: addon-essentials (controls, actions, viewport, …) and addon-interactions
    // are built into core. addon-docs provides docs/autodocs and the MDX blocks.
    "@storybook/addon-links",
    "@storybook/addon-docs",
  ],
  framework: {
    name: "@storybook/angular",
    options: {},
  },
};
export default config;
