// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const { themes } = require("prism-react-renderer");
const lightTheme = themes.github;
const darkTheme = themes.dracula;

const tabBlocksRemarkPlugin = require("docusaurus-remark-plugin-tab-blocks");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Jalapeno",
  tagline: "CLI for creating, managing and sharing spiced up project templates",
  url: "https://futurice.github.io",
  baseUrl: "/jalapeno/",
  trailingSlash: false,
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/logo.svg",

  // GitHub pages deployment config
  organizationName: "futurice",
  projectName: "jalapeno",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  markdown: {
    mermaid: true,
  },
  themes: ["@docusaurus/theme-mermaid"],

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/futurice/jalapeno/tree/main/docs/site",
          remarkPlugins: [tabBlocksRemarkPlugin],
        },
        blog: false,
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: "Jalapeno",
        logo: {
          alt: "Jalapeno Logo",
          src: "img/logo.svg",
        },
        items: [
          {
            type: "doc",
            docId: "installation",
            position: "left",
            label: "Installation",
          },
          {
            type: "doc",
            docId: "usage",
            position: "left",
            label: "Usage",
          },
          {
            type: "doc",
            docId: "api",
            position: "left",
            label: "API",
          },
          {
            href: "https://github.com/futurice/jalapeno",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: "Pages",
            items: [
              {
                label: "Installation",
                to: "/installation",
              },
              {
                label: "Usage",
                to: "/usage",
              },
              {
                label: "API",
                to: "/api",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "GitHub",
                href: "https://github.com/futurice/jalapeno",
              },
            ],
          },
        ],
      },
      prism: {
        theme: lightTheme,
        darkTheme: darkTheme,
        additionalLanguages: ["bash", "powershell"],
      },
    }),
};

module.exports = config;
