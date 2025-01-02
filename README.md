# NextgoStarter

A fullstack starter using Golang and Next.js with [Nx](https://nx.dev) integrated Monorepo

## Run tasks

To see all available targets to run for a project, run:

```sh
# shows available tasks for `api`
npx nx show project api

# shows available tasks for `web`
npx nx show project web
```

To visually explore what was created for workspace, run:

```sh
npx nx graph
```

These targets are either [inferred automatically](https://nx.dev/concepts/inferred-tasks?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects) or defined in the `project.json` or `package.json` files.

[More about running tasks in the docs &raquo;](https://nx.dev/features/run-tasks?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)

## Applications' Description

| Name      | Path             | Description                                                                                                                                       |
| --------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `api`     | `./apps/api`     | A REST API template in Golang. See details and available targets in [README](https://github.com/pohsiangchen/nextgo-starter/tree/main/apps/api)   |
| `web`     | `./apps/web`     | A web application using Next.js. See details and available targets in [README](https://github.com/pohsiangchen/nextgo-starter/tree/main/apps/web) |
| `web-e2e` | `./apps/web-e2e` | An end-to-end testing for `web` application using [Playwright](https://playwright.dev/).                                                          |

## Add new projects

While you could add new projects to your workspace manually, you might want to leverage [Nx plugins](https://nx.dev/concepts/nx-plugins?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects) and their [code generation](https://nx.dev/features/generate-code?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects) feature.

Use the plugin's generator to create new projects.

To generate a new Next.js application, use:

```sh
npx nx g @nx/next:app demo
```

To generate a new React library, use:

```sh
npx nx g @nx/react:lib mylib
```

To generate a new [Golang application](https://github.com/nx-go/nx-go/blob/main/docs/generators/application.md), use:

```sh
nx g @nx-go/nx-go:application <YOUR_APPLICATION_NAME> --directory apps/<YOUR_APPLICATION_NAME>
```

To generate a new [Golang library](https://github.com/nx-go/nx-go/blob/main/docs/generators/library.md), use:

```sh
nx g @nx-go/nx-go:library <YOUR_LIBRARY_NAME> --directory libs/<YOUR_LIBRARY_NAME>
```

You can use `npx nx list` to get a list of installed plugins. Then, run `npx nx list <plugin-name>` to learn about more specific capabilities of a particular plugin. Alternatively, [install Nx Console](https://nx.dev/getting-started/editor-setup?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects) to browse plugins and generators in your IDE.

[Learn more about Nx plugins &raquo;](https://nx.dev/concepts/nx-plugins?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects) | [Browse the plugin registry &raquo;](https://nx.dev/plugin-registry?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)


[Learn more about Nx on CI](https://nx.dev/ci/intro/ci-with-nx#ready-get-started-with-your-provider?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)

## Install Nx Console

Nx Console is an editor extension that enriches your developer experience. It lets you run tasks, generate code, and improves code autocompletion in your IDE. It is available for VSCode and IntelliJ.

[Install Nx Console &raquo;](https://nx.dev/getting-started/editor-setup?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)

## Useful links

[Learn more about Nx workspace setup and its capabilities](https://nx.dev/nx-api/next?utm_source=nx_project&amp;utm_medium=readme&amp;utm_campaign=nx_projects).

Learn more:

- [Learn more about this workspace setup](https://nx.dev/nx-api/next?utm_source=nx_project&amp;utm_medium=readme&amp;utm_campaign=nx_projects)
- [Learn about Nx on CI](https://nx.dev/ci/intro/ci-with-nx?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)
- [Releasing Packages with Nx release](https://nx.dev/features/manage-releases?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)
- [What are Nx plugins?](https://nx.dev/concepts/nx-plugins?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)

And join the Nx community:
- [Discord](https://go.nx.dev/community)
- [Follow us on X](https://twitter.com/nxdevtools) or [LinkedIn](https://www.linkedin.com/company/nrwl)
- [Our Youtube channel](https://www.youtube.com/@nxdevtools)
- [Our blog](https://nx.dev/blog?utm_source=nx_project&utm_medium=readme&utm_campaign=nx_projects)
