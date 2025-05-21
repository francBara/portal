# Portal 🌀

A no-code tool to edit existing codebases, the developers way.

**Portal** parses your code looking for *@portal* comments, then serves a web-based dashboard to edit, preview and push changes.

1. Annotate the code you want to be editable.
    ```
    @portal
    let minItems = 2;
    ```

2. Spin up the Portal container locally or in your cloud environment, configuring repo access and users.

    `docker run portal -p 8080:8080`

3. Log into the web dashboard and edit your project.

## Features
- Interactive, user friendly no-code dashboard.
- Total customization and limitation of code interaction.
- Seamless remote repo integration.
- For frontends, real time UI preview.


## Syntax

### Variable portal
Makes the annotated variable editable in the web dashboard.
```js
@portal
let minItems = 2;
```

### All variables portal
Makes all variables under the annotation in the current file editable.
```js
@portal all

let minItems = 2;

const maxItems = 10;
```

### UI variables portal
If there is a React/Tailwind component tree, it makes all CSS parameters in the component tree editable.
```html
@portal ui

return (
    <div className="h-8 w-8">
        <span className="flex mt-8">
        ...
);
```

## Server configuration

The server can be configured via env variables.
The scope of the configuration are remote repo access and dashboard authentication.

| Variable Name         | Required | Default    | Description                                  |
|-----------------------|----------|------------|----------------------------------------------|
| `REPO_OWNER`          | No       | —          | Owner of the repository. Defaults to `USER_NAME` if unset. |
| `GITHUB_USERNAME`           | Yes      | —          | Git username used for cloning/pushing.       |
| `REPO_NAME`           | Yes      | —          | Name of the repository to work with.         |
| `REPO_BRANCH`         | No       | `main`     | The branch to use for operations.            |
| `PAC`                 | Yes      | —          | Personal access token or authentication secret. |
| `OPEN_PULL_REQUEST`   | No       | `true`     | Whether to automatically open pull requests, in this case a new branch will be created. |
| `SERVE_PREVIEW`       | No       | `true`     | Whether to serve a realtime preview, in case of React projects.            |
| `ADMIN_USERNAME`      | No       | `admin`    | The username to log in the dashboard.        |
| `ADMIN_PASSWORD`      | No       | `admin`    | The password to log in the dashboard.        |

Minimum working `.env`
```
REPO_NAME      = portal_demo
USER_NAME      = myGithubUser
PAC            = my_personal_access_token

ADMIN_USERNAME = pippo
ADMIN_PASSWORD = mypassword
```


## Components

- **Parser**: Scans your codebase for custom *@portal* annotations and extracts structured variables.
- **Patcher**: Edits your codebase, replacing annotated variables with new values.
- **Webserver**: Spins up a secure web-based dashboard where users can update variables and preview/push changes.

## Currently supported technologies

- Javascript
- Typescript
- React/Tailwind