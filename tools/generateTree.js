const fs = require('fs');
const parser = require('@babel/parser');
const traverse = require('@babel/traverse').default;
const t = require('@babel/types');

const sourceCode = fs.readFileSync(process.argv[2], 'utf8');

const ast = parser.parse(sourceCode, {
    sourceType: 'module',
    plugins: ['typescript', 'jsx'],
});

const components = {};
const comments = {};
const props = {};

const rootIds = {};

traverse(ast, {
    FunctionDeclaration(path) {
        const rootName = path.node.id.name;

        if (!path.node.leadingComments) {
            return;
        }

        let hasComment = false;

        for (let comment of path.node.leadingComments) {
            const trimmedComment = comment.value.trim();

            if (trimmedComment.includes("@portal") && trimmedComment.includes("ui")) {
                hasComment = true;
                comments[rootName] = [];
                props[rootName] = [];

                const params = path.node.params;

                if (params[0] && params[0].type === 'ObjectPattern') {
                    for (let p of params[0].properties) {
                        if (p.key && p.key.name) {
                            props[rootName].push(p.key.name);
                        }
                    }
                }
                else if (params[0] && params[0].type === 'Identifier') {
                    props[rootName].push(params[0].name);
                }
            }
            if (hasComment) {
                comments[rootName].push(trimmedComment);
            }
        }

        if (!hasComment) {
            return;
        }

        for (let el of path.node.body.body) {
            if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                rootIds[rootName] = 0;
                components[rootName] = collectJSX(el.argument, rootName);
            }
        }
    },
    VariableDeclarator(path) {
        const rootName = path.node.id.name;

        if (!path.node.leadingComments) {
            return;
        }

        let hasComment = false;

        for (let comment of path.node.leadingComments) {
            const trimmedComment = comment.value.trim();

            if (trimmedComment.includes("@portal") && trimmedComment.includes("ui")) {
                hasComment = true;
                comments[rootName] = [];
            }
            if (hasComment) {
                comments[rootName].push(trimmedComment);
            }
        }

        if (!hasComment) {
            return;
        }

        if (path.node.init && path.node.init.type === "ArrowFunctionExpression" && path.node.init.body && path.node.init.body.body) {
            for (let el of path.node.init.body.body) {
                if (t.isReturnStatement(el) && t.isJSXElement(el.argument)) {
                    rootIds[rootName] = 0;
                    components[rootName] = collectJSX(el.argument, rootName);
                }
            }
        }
    }
});

function collectJSX(node, rootName) {
    if (!t.isJSXElement(node)) return;

    const element = {
        type: node.openingElement.name.name || 'Unknown',
        row: node.openingElement.loc.start.line,
        id: rootIds[rootName],
        properties: [],
        children: [],
    };

    rootIds[rootName]++;

    node.openingElement.attributes.forEach(attr => {
        if (t.isJSXAttribute(attr) && t.isJSXIdentifier(attr.name) && attr.name.name === "className") {
            if (t.isStringLiteral(attr.value)) {
                element.row = attr.loc.start.line;
                element.properties = parseTailwindString(attr.value.value);
            }
            else if (t.isJSXExpressionContainer(attr.value) && t.isTemplateLiteral(attr.value.expression)) {
                if (attr.value.expression.expressions.length > 0) {
                    //TODO: Implement expression container parsing
                }
            }
        }
    });

    node.children.forEach(child => {
        if (t.isJSXElement(child)) {
            element.children.push(collectJSX(child, rootName));
        }
        else if (t.isJSXText(child)) {
            const trimmed = child.value.trim();
            if (trimmed) {
                element.children.push({ type: 'text', properties: [{prefix: "content", value: trimmed}], id: rootIds[rootName], children: [] });
                rootIds[rootName]++;
            }
        }
    });

    return element;
}

function parseTailwindString(tailwind) {
    const splitted = tailwind.split(" ");

    const result = [];

    for (let word of splitted) {
        const lastHyphen = word.lastIndexOf("-");

        let prefix = word.substring(0, lastHyphen);
        let value = word.substring(lastHyphen + 1);

        if (prefix.length === 0) {
            prefix = value;
            value = "";
        }

        result.push({ prefix, value });
    }

    return result;
}

console.log(JSON.stringify({components, props, comments}, null, 2));
