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
```
@portal
let minItems = 2;
```

### All variables portal
Makes all variables under the annotation in the current file editable.
```
@portal all

let minItems = 2;

const maxItems = 10;
```

### UI variables portal
If there is a React/Tailwind component tree, it makes all CSS parameters in the component tree editable.
```
@portal ui

return (
    <div className="h-8 w-8">
        <span className="flex mt-8">
        ...
);
```

## Server configuration

The server can be configured via env variables, config file and CLI arguments. 

The scope of configuration is to repo access and Portal admin user.

`config.json`
```
{
    "repoOwner": "francBara",
    "repoName": "portal-demo",
    "repoBranch": "main",
    "pac": "YOUR_PERSONAL_ACCESS_TOKEN",
    "adminUsername": "admin",
    "adminPassword": "XXX"
}
```

## Components

- **Parser**: Scans your codebase for custom *@portal* annotations and extracts structured variables.
- **Patcher**: Edits your codebase, replacing annotated variables with new values.
- **Webserver**: Spins up a secure web-based dashboard where users can update variables and preview/push changes.

## Currently supported technologies

- Javascript
- Typescript
- React/Tailwind